package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

var _ OperationLogModel = (*defaultOperationLogModel)(nil)

type (
	OperationLogModel interface {
		Insert(ctx context.Context, data *OperationLog) (sql.Result, error)
		FindList(ctx context.Context, action, username string, page, pageSize int64) ([]*OperationLog, int64, error)
		Count(ctx context.Context) (int64, error)
		DeleteById(ctx context.Context, id int64) error
		DeleteBatch(ctx context.Context, ids []int64) error
		DeleteOlderThan(ctx context.Context, beforeTime time.Time) (int64, error)
	}

	defaultOperationLogModel struct {
		db *sql.DB
	}
)

func NewOperationLogModel(db *sql.DB) OperationLogModel {
	return &defaultOperationLogModel{db: db}
}

func (m *defaultOperationLogModel) tableName() string { return "`operation_logs`" }

func (m *defaultOperationLogModel) Insert(ctx context.Context, data *OperationLog) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, username, module, action, method, ip, location, params, status, error_msg, duration, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, m.tableName())
	data.CreatedAt = time.Now()
	return m.db.ExecContext(ctx, query, data.UserId, data.Username, data.Module, data.Action, data.Method, data.Ip, data.Location, data.Params, data.Status, data.ErrorMsg, data.Duration, data.CreatedAt)
}

func (m *defaultOperationLogModel) FindList(ctx context.Context, action, username string, page, pageSize int64) ([]*OperationLog, int64, error) {
	where := " WHERE 1=1"
	var args []interface{}
	if action != "" {
		where += " AND action = ?"
		args = append(args, action)
	}
	if username != "" {
		where += " AND username LIKE ?"
		args = append(args, "%"+username+"%")
	}

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, m.tableName()) + where
	if err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`SELECT id, user_id, username, module, action, method, ip, location, params, status, error_msg, duration, created_at FROM %s`, m.tableName()) + where + ` ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*OperationLog
	for rows.Next() {
		var l OperationLog
		err := rows.Scan(&l.Id, &l.UserId, &l.Username, &l.Module, &l.Action, &l.Method, &l.Ip, &l.Location, &l.Params, &l.Status, &l.ErrorMsg, &l.Duration, &l.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, &l)
	}
	return list, total, rows.Err()
}

func (m *defaultOperationLogModel) Count(ctx context.Context) (int64, error) {
	var count int64
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, m.tableName())
	err := m.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (m *defaultOperationLogModel) DeleteById(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m *defaultOperationLogModel) DeleteOlderThan(ctx context.Context, beforeTime time.Time) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE created_at < ?`, m.tableName())
	result, err := m.db.ExecContext(ctx, query, beforeTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultOperationLogModel) DeleteBatch(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE id IN (%s)`, m.tableName(), strings.Join(placeholders, ","))
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
