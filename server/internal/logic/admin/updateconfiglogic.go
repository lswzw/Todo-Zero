package admin

import (
	"context"

	"server/internal/pkg/jwtx"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConfigLogic {
	return &UpdateConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateConfigLogic) UpdateConfig(req *types.UpdateConfigReq) (resp *types.UpdateConfigResp, err error) {
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}

	config, err := l.svcCtx.SystemConfigModel.FindOneByKey(l.ctx, req.Key)
	if err != nil {
		return nil, xerr.NewCodeErrFromMsg("配置项不存在")
	}

	config.ConfigValue = req.Value
	if err := l.svcCtx.SystemConfigModel.Update(l.ctx, config); err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	return &types.UpdateConfigResp{}, nil
}

func (l *UpdateConfigLogic) checkAdmin() error {
	isAdmin, err := jwtx.GetIsAdminFromCtx(l.ctx)
	if err != nil || isAdmin != 1 {
		return xerr.NewCodeError(xerr.AdminRequired)
	}
	return nil
}
