package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	maxLoginAttempts = 10              // 最大尝试次数
	loginWindow      = 15 * time.Minute // 统计窗口
	lockoutDuration  = 15 * time.Minute // 锁定时长
	cleanupInterval  = 10 * time.Minute // 清理间隔
)

type loginAttempt struct {
	count    int
	lastTime time.Time
	locked   bool
	lockTime time.Time
}

type LoginRateLimitMiddleware struct {
	mu       sync.Mutex
	attempts map[string]*loginAttempt
}

func NewLoginRateLimitMiddleware() *LoginRateLimitMiddleware {
	m := &LoginRateLimitMiddleware{
		attempts: make(map[string]*loginAttempt),
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
			attempt = &loginAttempt{}
			m.attempts[ip] = attempt
		}

		// 检查是否被锁定
		if attempt.locked {
			if time.Since(attempt.lockTime) < lockoutDuration {
				remaining := lockoutDuration - time.Since(attempt.lockTime)
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
		if time.Since(attempt.lastTime) > loginWindow {
			attempt.count = 0
		}

		attempt.lastTime = time.Now()
		attempt.count++
		shouldLock := attempt.count > maxLoginAttempts
		if shouldLock {
			attempt.locked = true
			attempt.lockTime = time.Now()
		}
		m.mu.Unlock()

		if shouldLock {
			logx.Errorf("[LoginRateLimit] IP %s exceeded max attempts, locking out", ip)
			w.Header().Set("Retry-After", time.Now().Add(lockoutDuration).Format(time.RFC1123))
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

func getClientIP(r *http.Request) string {
	// 优先从 X-Forwarded-For / X-Real-IP 获取（反向代理场景）
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For 可能包含多个 IP，取第一个
		for _, ip := range splitByComma(xff) {
			if ip := trimSpace(ip); ip != "" {
				return ip
			}
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return trimSpace(xri)
	}
	// RemoteAddr 包含端口，如 192.168.1.1:12345
	addr := r.RemoteAddr
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i]
		}
	}
	return addr
}

func splitByComma(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

type loginRateLimitError struct{}

func (e *loginRateLimitError) Error() string { return "登录尝试次数过多，请稍后再试" }

var errLoginRateLimit = &loginRateLimitError{}
