package admin

import (
	"context"
	"strings"

	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var sensitiveConfigPrefixes = []string{
	"secret",
	"key",
	"password",
	"token",
	"credential",
	"private",
}

func isSensitiveConfigKey(key string) bool {
	lowerKey := strings.ToLower(key)
	for _, prefix := range sensitiveConfigPrefixes {
		if strings.HasPrefix(lowerKey, prefix) {
			return true
		}
	}
	return false
}

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
	configs, err := l.svcCtx.SystemConfigModel.FindAll(l.ctx)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ServerCommonError)
	}

	var list []types.ConfigItem
	for _, c := range configs {
		if isSensitiveConfigKey(c.ConfigKey) {
			list = append(list, types.ConfigItem{
				Key:    c.ConfigKey,
				Value:  "******",
				Remark: c.Description,
			})
			continue
		}
		list = append(list, types.ConfigItem{
			Key:    c.ConfigKey,
			Value:  c.ConfigValue,
			Remark: c.Description,
		})
	}

	return &types.ConfigListResp{List: list}, nil
}
