package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"server/internal/pkg/xerr"
)

const (
	defaultMaxLoginAttempts = 10               // 默认最大尝试次数
	defaultLoginWindow      = 15 * time.Minute // 默认统计窗口
	defaultLockoutDuration  = 15 * time.Minute // 默认锁定时长
	cleanupInterval         = 10 * time.Minute // 清理间隔
)

type loginAttempt struct {
	count       int
	lastTime    time.Time
	locked      bool
	lockTime    time.Time
	lockoutDur  time.Duration
	window      time.Duration
	maxAttempts int
}

type LoginRateLimitMiddleware struct {
	mu          sync.Mutex
	attempts    map[string]*loginAttempt
	maxAttempts int
	window      time.Duration
	lockoutDur  time.Duration
}

func NewLoginRateLimitMiddleware() *LoginRateLimitMiddleware {
	return NewLoginRateLimitMiddlewareWithConfig(defaultMaxLoginAttempts, defaultLoginWindow, defaultLockoutDuration)
}

func NewLoginRateLimitMiddlewareWithConfig(maxAttempts int, window, lockoutDur time.Duration) *LoginRateLimitMiddleware {
	if maxAttempts <= 0 {
		maxAttempts = defaultMaxLoginAttempts
	}
	if window <= 0 {
		window = defaultLoginWindow
	}
	if lockoutDur <= 0 {
		lockoutDur = defaultLockoutDuration
	}
	m := &LoginRateLimitMiddleware{
		attempts:    make(map[string]*loginAttempt),
		maxAttempts: maxAttempts,
		window:      window,
		lockoutDur:  lockoutDur,
	}
	go m.cleanup()
	return m
}

func (m *LoginRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		m.mu.Lock()
		attempt, exists := m.attempts[ip]
		if !exists {
			attempt = &loginAttempt{
				lockoutDur:  m.lockoutDur,
				window:      m.window,
				maxAttempts: m.maxAttempts,
			}
			m.attempts[ip] = attempt
		}

		// 检查是否被锁定
		if attempt.locked {
			if time.Since(attempt.lockTime) < m.lockoutDur {
				remaining := m.lockoutDur - time.Since(attempt.lockTime)
				m.mu.Unlock()
				logx.Errorf("[LoginRateLimit] IP %s is locked out, remaining: %v", ip, remaining)
				w.Header().Set("Retry-After", time.Now().Add(remaining).Format(time.RFC1123))
				httpx.ErrorCtx(r.Context(), w, errLoginRateLimit)
				return
			}
			// 锁定已过期，重置
			attempt.locked = false
			attempt.count = 0
		}

		// 检查窗口是否过期，过期则重置
		if time.Since(attempt.lastTime) > m.window {
			attempt.count = 0
		}

		attempt.lastTime = time.Now()
		attempt.count++
		shouldLock := attempt.count > m.maxAttempts
		if shouldLock {
			attempt.locked = true
			attempt.lockTime = time.Now()
		}
		m.mu.Unlock()

		if shouldLock {
			logx.Errorf("[LoginRateLimit] IP %s exceeded max attempts, locking out", ip)
			w.Header().Set("Retry-After", time.Now().Add(m.lockoutDur).Format(time.RFC1123))
			httpx.ErrorCtx(r.Context(), w, errLoginRateLimit)
			return
		}

		next(w, r)
	}
}

// RecordLoginFailure 记录登录失败（由 login logic 调用）
// 当前实现已在 Handle 中预计数，此方法预留用于按用户名限流的扩展
func (m *LoginRateLimitMiddleware) RecordLoginFailure(ip string) {
	// 当前预计数方式已满足需求，预留扩展接口
}

// RecordLoginSuccess 记录登录成功，清除该 IP 的失败计数
func (m *LoginRateLimitMiddleware) RecordLoginSuccess(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.attempts, ip)
}

func (m *LoginRateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for ip, attempt := range m.attempts {
			// 锁定已过期或窗口已过，清除
			if (attempt.locked && now.Sub(attempt.lockTime) > lockoutDuration) ||
				(!attempt.locked && now.Sub(attempt.lastTime) > loginWindow) {
				delete(m.attempts, ip)
			}
		}
		m.mu.Unlock()
	}
}

var errLoginRateLimit = xerr.NewRateLimitError("登录尝试次数过多，请稍后再试")
