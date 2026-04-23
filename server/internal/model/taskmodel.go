package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TaskModel = (*customTaskModel)(nil)

type (
	// TaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskModel.
	TaskModel interface {
		taskModel
		FindList(ctx context.Context, userId int64, status, categoryId, priority int64, keyword string, page, pageSize int64) ([]*Task, int64, error)
	}

	customTaskModel struct {
		*defaultTaskModel
	}
)

// NewTaskModel returns a model for the database table.
func NewTaskModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) TaskModel {
	return &customTaskModel{
		defaultTaskModel: newTaskModel(conn, c, opts...),
	}
}

func (m *customTaskModel) FindList(ctx context.Context, userId int64, status, categoryId, priority int64, keyword string, page, pageSize int64) ([]*Task, int64, error) {
	var where string
	var args []interface{}

	where += " WHERE user_id = ?"
	args = append(args, userId)

	if status != 0 {
		where += " AND status = ?"
		args = append(args, status)
	}
	if categoryId != 0 {
		where += " AND category_id = ?"
		args = append(args, categoryId)
	}
	if priority != 0 {
		where += " AND priority = ?"
		args = append(args, priority)
	}
	if keyword != "" {
		where += " AND title LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	// 查总数
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s%s", m.table, where)
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查列表
	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf("SELECT %s FROM %s%s ORDER BY id DESC LIMIT ? OFFSET ?", taskRows, m.table, where)
	listArgs := append(args, pageSize, offset)

	var list []*Task
	err = m.QueryRowsNoCacheCtx(ctx, &list, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
