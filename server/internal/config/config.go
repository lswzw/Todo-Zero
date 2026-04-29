package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire  int64
	}
	Database struct {
		DataDir string // e.g. "data"
		DBFile  string // e.g. "todo.db"
	}
	Debug           bool     // enable debug mode: API docs, verbose logging, etc.
	TrustedProxies  []string // list of trusted proxy IPs (e.g. "127.0.0.1", "10.0.0.0/8")
}
