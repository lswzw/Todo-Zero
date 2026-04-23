package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CategoryModel = (*customCategoryModel)(nil)

type (
	CategoryModel interface {
		categoryModel
		FindAll(ctx context.Context) ([]*Category, error)
	}

	customCategoryModel struct {
		*defaultCategoryModel
	}
)

func NewCategoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CategoryModel {
	return &customCategoryModel{
		defaultCategoryModel: newCategoryModel(conn, c, opts...),
	}
}

func (m *customCategoryModel) FindAll(ctx context.Context) ([]*Category, error) {
	var list []*Category
	query := fmt.Sprintf("select %s from %s order by sort_order asc", categoryRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query)
	if err != nil {
		return nil, err
	}
	return list, nil
}
