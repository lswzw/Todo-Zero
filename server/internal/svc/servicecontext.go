package svc

import (
	"database/sql"

	"server/internal/config"
	"server/internal/middleware"
	"server/internal/model"
)

type ServiceContext struct {
	Config config.Config

	UserModel         model.UserModel
	TaskModel         model.TaskModel
	CategoryModel     model.CategoryModel
	TagModel          model.TagModel
	TaskTagModel      model.TaskTagModel
	SystemConfigModel model.SystemConfigModel
	OperationLogModel model.OperationLogModel
	LoginLogModel     model.LoginLogModel
	DB                *sql.DB

	AdminMiddleware           *middleware.AdminMiddleware
	OperationLogMiddleware    *middleware.OperationLogMiddleware
	SecurityHeadersMiddleware *middleware.SecurityHeadersMiddleware
	LoginRateLimitMiddleware  *middleware.LoginRateLimitMiddleware
	APIRateLimitMiddleware    *middleware.APIRateLimitMiddleware
	LocaleMiddleware          *middleware.LocaleMiddleware
}

func NewServiceContext(c config.Config, db *sql.DB) *ServiceContext {
	userModel := model.NewUserModel(db)
	opLogModel := model.NewOperationLogModel(db)

	return &ServiceContext{
		Config:                    c,
		DB:                        db,
		UserModel:                 userModel,
		TaskModel:                 model.NewTaskModel(db),
		CategoryModel:             model.NewCategoryModel(db),
		TagModel:                  model.NewTagModel(db),
		TaskTagModel:              model.NewTaskTagModel(db),
		SystemConfigModel:         model.NewSystemConfigModel(db),
		OperationLogModel:         opLogModel,
		LoginLogModel:             model.NewLoginLogModel(db),
		AdminMiddleware:           middleware.NewAdminMiddleware(),
		OperationLogMiddleware:    middleware.NewOperationLogMiddleware(userModel, opLogModel),
		SecurityHeadersMiddleware: middleware.NewSecurityHeadersMiddleware(),
		LoginRateLimitMiddleware:  middleware.NewLoginRateLimitMiddleware(),
		APIRateLimitMiddleware:    middleware.NewAPIRateLimitMiddleware(),
		LocaleMiddleware:          middleware.NewLocaleMiddleware(),
	}
}
