package admin

import (
	"context"

	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigListLogic {
	return &ConfigListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigListLogic) ConfigList() (resp *types.ConfigListResp, err error) {
	if err := l.checkAdmin(); err != nil {
		return nil, err
	}

	configs, err := l.svcCtx.SystemConfigModel.FindAll(l.ctx)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var list []types.ConfigItem
	for _, c := range configs {
		list = append(list, types.ConfigItem{
			Key:    c.ConfigKey,
			Value:  c.ConfigValue,
			Remark: c.Remark.String,
		})
	}

	return &types.ConfigListResp{List: list}, nil
}

func (l *ConfigListLogic) checkAdmin() error {
	isAdmin, ok := l.ctx.Value("isAdmin").(float64)
	if !ok || isAdmin != 1 {
		return xerr.NewCodeError(xerr.AdminRequired)
	}
	return nil
}
