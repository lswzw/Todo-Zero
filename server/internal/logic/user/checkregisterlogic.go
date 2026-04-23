package user

import (
	"context"

	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckRegisterLogic {
	return &CheckRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckRegisterLogic) CheckRegister() (resp *types.CheckRegisterResp, err error) {
	// 默认允许注册
	allowRegister := true

	// 查询系统配置
	config, err := l.svcCtx.SystemConfigModel.FindOneByKey(l.ctx, "allow_register")
	if err == nil && config.ConfigValue == "false" {
		allowRegister = false
	}

	return &types.CheckRegisterResp{
		AllowRegister: allowRegister,
	}, nil
}
