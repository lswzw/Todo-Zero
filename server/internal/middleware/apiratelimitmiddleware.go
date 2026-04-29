package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"server/internal/pkg/xerr"
)

const (
	maxAPIRequests      = 100              // 每分钟最大请求数
	apiRateLimitWindow  = 1 * time.Minute  // 统计窗口
	apiCleanupInterval  = 5 * time.Minute  // 清理间隔
)

type apiRequest struct {
	count    int
	lastTime time.Time
}

type APIRateLimitMiddleware struct {
	mu       sync.Mutex
	requests map[string]*apiRequest
}

func NewAPIRateLimitMiddleware() *APIRateLimitMiddleware {
	m := &APIRateLimitMiddleware{
		requests: make(map[string]*apiRequest),
	}
	go m.cleanup()
	return m
}

func (m *APIRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		m.mu.Lock()
		req, exists := m.requests[ip]
		if !exists {
			req = &apiRequest{}
			m.requests[ip] = req
		}

		// 检查窗口是否过期，过期则重置计数
		if time.Since(req.lastTime) > apiRateLimitWindow {
			req.count = 0
		}

		req.lastTime = time.Now()
		req.count++

		// 检查是否超过限制
		if req.count > maxAPIRequests {
			m.mu.Unlock()
			logx.Errorf("[APIRateLimit] IP %s exceeded rate limit: %d requests/min", ip, req.count)
			w.Header().Set("Retry-After", time.Now().Add(apiRateLimitWindow).Format(time.RFC1123))
			w.Header().Set("X-RateLimit-Limit", "100")
			w.Header().Set("X-RateLimit-Remaining", "0")
			httpx.ErrorCtx(r.Context(), w, errAPIRateLimit)
			return
		}

		// 设置响应头
		remaining := maxAPIRequests - req.count
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		m.mu.Unlock()

		next(w, r)
	}
}

func (m *APIRateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(apiCleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for ip, req := range m.requests {
			if now.Sub(req.lastTime) > apiRateLimitWindow {
				delete(m.requests, ip)
			}
		}
		m.mu.Unlock()
	}
}

var errAPIRateLimit = xerr.NewRateLimitError("API 请求过于频繁，请稍后再试")
