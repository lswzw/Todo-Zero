package svc

import (
	"database/sql"

	"server/internal/config"
	"server/internal/model"
)

type ServiceContext struct {
	Config config.Config

	UserModel         model.UserModel
	TaskModel         model.TaskModel
	CategoryModel     model.CategoryModel
	SystemConfigModel model.SystemConfigModel
	OperationLogModel model.OperationLogModel
	LoginLogModel     model.LoginLogModel
}

func NewServiceContext(c config.Config, db *sql.DB) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		UserModel:         model.NewUserModel(db),
		TaskModel:         model.NewTaskModel(db),
		CategoryModel:     model.NewCategoryModel(db),
		SystemConfigModel: model.NewSystemConfigModel(db),
		OperationLogModel: model.NewOperationLogModel(db),
		LoginLogModel:     model.NewLoginLogModel(db),
	}
}
