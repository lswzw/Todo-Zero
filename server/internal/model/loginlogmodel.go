package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoginLogModel = (*customLoginLogModel)(nil)

type (
	// LoginLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoginLogModel.
	LoginLogModel interface {
		loginLogModel
	}

	customLoginLogModel struct {
		*defaultLoginLogModel
	}
)

// NewLoginLogModel returns a model for the database table.
func NewLoginLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LoginLogModel {
	return &customLoginLogModel{
		defaultLoginLogModel: newLoginLogModel(conn, c, opts...),
	}
}
