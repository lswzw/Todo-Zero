package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoginLogModel = (*customLoginLogModel)(nil)

type (
	LoginLogModel interface {
		loginLogModel
		FindList(ctx context.Context, username string, page, pageSize int64) ([]*LoginLog, int64, error)
	}

	customLoginLogModel struct {
		*defaultLoginLogModel
	}
)

func NewLoginLogModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LoginLogModel {
	return &customLoginLogModel{
		defaultLoginLogModel: newLoginLogModel(conn, c, opts...),
	}
}

func (m *customLoginLogModel) FindList(ctx context.Context, username string, page, pageSize int64) ([]*LoginLog, int64, error) {
	var where string
	var args []interface{}

	where = " WHERE 1=1"
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
	listQuery := fmt.Sprintf("SELECT %s FROM %s%s ORDER BY id DESC LIMIT ? OFFSET ?", loginLogRows, m.table, where)
	listArgs := append(args, pageSize, offset)

	var list []*LoginLog
	err = m.QueryRowsNoCacheCtx(ctx, &list, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
