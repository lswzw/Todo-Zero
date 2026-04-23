package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OperationLogModel = (*customOperationLogModel)(nil)

type (
	// OperationLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOperationLogModel.
	OperationLogModel interface {
		operationLogModel
	}

	customOperationLogModel struct {
		*defaultOperationLogModel
	}
)

// NewOperationLogModel returns a model for the database table.
func NewOperationLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OperationLogModel {
	return &customOperationLogModel{
		defaultOperationLogModel: newOperationLogModel(conn, c, opts...),
	}
}
