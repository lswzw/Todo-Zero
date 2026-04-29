package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Database struct {
		DataDir string // e.g. "data"
		DBFile  string // e.g. "todo.db"
	}
	Debug          bool     // enable debug mode: API docs, verbose logging, etc.
	TrustedProxies []string // list of trusted proxy IPs (e.g. "127.0.0.1", "10.0.0.0/8")
	AllowedOrigins []string // allowed CORS origins (empty = allow all)
	RateLimit      struct {
		APIRequestsPerMinute int // API 限流：每分钟最大请求数
		LoginAttempts        int // 登录限流：最大尝试次数
		LoginWindowMinutes   int // 登录限流：统计窗口（分钟）
		LoginLockoutMinutes  int // 登录限流：锁定时长（分钟）
	}
}
