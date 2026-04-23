package svc

import (
	"server/internal/config"
	"server/internal/model"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	Redis             *redis.Redis
	UserModel         model.UserModel
	TaskModel         model.TaskModel
	CategoryModel     model.CategoryModel
	SystemConfigModel model.SystemConfigModel
	OperationLogModel model.OperationLogModel
	LoginLogModel     model.LoginLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)
	rds := redis.MustNewRedis(redis.RedisConf{
		Host: c.Redis.Host,
		Pass: c.Redis.Pass,
		Type: c.Redis.Type,
	})

	// 带 Redis 缓存的 model
	cacheConf := cache.CacheConf{
		{
			RedisConf: redis.RedisConf{
				Host: c.Redis.Host,
				Pass: c.Redis.Pass,
				Type: c.Redis.Type,
			},
		},
	}

	return &ServiceContext{
		Config:            c,
		Redis:             rds,
		UserModel:         model.NewUserModel(conn, cacheConf),
		TaskModel:         model.NewTaskModel(conn, cacheConf),
		CategoryModel:     model.NewCategoryModel(conn, cacheConf),
		SystemConfigModel: model.NewSystemConfigModel(conn, cacheConf),
		OperationLogModel: model.NewOperationLogModel(conn, cacheConf),
		LoginLogModel:     model.NewLoginLogModel(conn, cacheConf),
	}
}
