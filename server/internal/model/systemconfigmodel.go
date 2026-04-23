package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SystemConfigModel = (*customSystemConfigModel)(nil)

type (
	// SystemConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSystemConfigModel.
	SystemConfigModel interface {
		systemConfigModel
	}

	customSystemConfigModel struct {
		*defaultSystemConfigModel
	}
)

// NewSystemConfigModel returns a model for the database table.
func NewSystemConfigModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SystemConfigModel {
	return &customSystemConfigModel{
		defaultSystemConfigModel: newSystemConfigModel(conn, c, opts...),
	}
}
