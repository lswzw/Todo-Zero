package svc

import (
	"server/internal/config"
	"server/internal/model"

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

	// 从 CacheRedis 配置中取第一个节点创建 Redis 实例
	rds := redis.MustNewRedis(c.CacheRedis[0].RedisConf)

	return &ServiceContext{
		Config:            c,
		Redis:             rds,
		UserModel:         model.NewUserModel(conn, c.CacheRedis),
		TaskModel:         model.NewTaskModel(conn, c.CacheRedis),
		CategoryModel:     model.NewCategoryModel(conn, c.CacheRedis),
		SystemConfigModel: model.NewSystemConfigModel(conn, c.CacheRedis),
		OperationLogModel: model.NewOperationLogModel(conn, c.CacheRedis),
		LoginLogModel:     model.NewLoginLogModel(conn, c.CacheRedis),
	}
}
