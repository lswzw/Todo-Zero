package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	UserModel interface {
		userModel
		FindList(ctx context.Context, keyword string, page, pageSize int64) ([]*User, int64, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *customUserModel) FindList(ctx context.Context, keyword string, page, pageSize int64) ([]*User, int64, error) {
	var where string
	var args []interface{}

	if keyword != "" {
		where = " WHERE username LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s%s", m.table, where)
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf("SELECT %s FROM %s%s ORDER BY id DESC LIMIT ? OFFSET ?", userRows, m.table, where)
	listArgs := append(args, pageSize, offset)

	var list []*User
	err = m.QueryRowsNoCacheCtx(ctx, &list, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
