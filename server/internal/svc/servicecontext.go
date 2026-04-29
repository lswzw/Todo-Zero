package svc

import (
	"database/sql"
	"time"

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

	middleware.SetTrustedProxies(c.TrustedProxies)

	// 配置化限流参数
	apiRateLimit := middleware.NewAPIRateLimitMiddleware()
	loginRateLimit := middleware.NewLoginRateLimitMiddleware()

	if c.RateLimit.APIRequestsPerMinute > 0 {
		apiRateLimit = middleware.NewAPIRateLimitMiddlewareWithConfig(
			c.RateLimit.APIRequestsPerMinute,
			1*time.Minute,
		)
	}
	if c.RateLimit.LoginAttempts > 0 {
		loginRateLimit = middleware.NewLoginRateLimitMiddlewareWithConfig(
			c.RateLimit.LoginAttempts,
			time.Duration(c.RateLimit.LoginWindowMinutes)*time.Minute,
			time.Duration(c.RateLimit.LoginLockoutMinutes)*time.Minute,
		)
	}

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
		LoginRateLimitMiddleware:  loginRateLimit,
		APIRateLimitMiddleware:    apiRateLimit,
		LocaleMiddleware:          middleware.NewLocaleMiddleware(),
	}
}
