package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	MySQL struct {
		DataSource string
	}
	Redis struct {
		Host string
		Pass string
		Type string
	}
}
