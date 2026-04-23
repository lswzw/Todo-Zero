package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SystemConfigModel = (*customSystemConfigModel)(nil)

type (
	SystemConfigModel interface {
		systemConfigModel
		FindOneByKey(ctx context.Context, key string) (*SystemConfig, error)
		FindAll(ctx context.Context) ([]*SystemConfig, error)
	}

	customSystemConfigModel struct {
		*defaultSystemConfigModel
	}
)

func NewSystemConfigModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SystemConfigModel {
	return &customSystemConfigModel{
		defaultSystemConfigModel: newSystemConfigModel(conn, c, opts...),
	}
}

func (m *customSystemConfigModel) FindOneByKey(ctx context.Context, key string) (*SystemConfig, error) {
	var resp SystemConfig
	query := fmt.Sprintf("select %s from %s where `config_key` = ? limit 1", systemConfigRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, key)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (m *customSystemConfigModel) FindAll(ctx context.Context) ([]*SystemConfig, error) {
	var list []*SystemConfig
	query := fmt.Sprintf("select %s from %s order by id asc", systemConfigRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query)
	if err != nil {
		return nil, err
	}
	return list, nil
}
