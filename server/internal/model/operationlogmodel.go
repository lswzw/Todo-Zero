package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OperationLogModel = (*customOperationLogModel)(nil)

type (
	OperationLogModel interface {
		operationLogModel
		FindList(ctx context.Context, action, username string, page, pageSize int64) ([]*OperationLog, int64, error)
	}

	customOperationLogModel struct {
		*defaultOperationLogModel
	}
)

func NewOperationLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OperationLogModel {
	return &customOperationLogModel{
		defaultOperationLogModel: newOperationLogModel(conn, c, opts...),
	}
}

func (m *customOperationLogModel) FindList(ctx context.Context, action, username string, page, pageSize int64) ([]*OperationLog, int64, error) {
	var where string
	var args []interface{}

	where = " WHERE 1=1"
	if action != "" {
		where += " AND action = ?"
		args = append(args, action)
	}
	if username != "" {
		where += " AND username LIKE ?"
		args = append(args, "%"+username+"%")
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s%s", m.table, where)
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf("SELECT %s FROM %s%s ORDER BY id DESC LIMIT ? OFFSET ?", operationLogRows, m.table, where)
	listArgs := append(args, pageSize, offset)

	var list []*OperationLog
	err = m.QueryRowsNoCacheCtx(ctx, &list, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
