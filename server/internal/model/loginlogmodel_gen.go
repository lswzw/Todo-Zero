package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

var _ LoginLogModel = (*defaultLoginLogModel)(nil)

type (
	LoginLogModel interface {
		Insert(ctx context.Context, data *LoginLog) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*LoginLog, error)
		Update(ctx context.Context, data *LoginLog) error
		Delete(ctx context.Context, id int64) error
		DeleteBatch(ctx context.Context, ids []int64) error
		FindList(ctx context.Context, username string, page, pageSize int64) ([]*LoginLog, int64, error)
		DeleteOlderThan(ctx context.Context, beforeTime time.Time) (int64, error)
	}

	defaultLoginLogModel struct {
		db *sql.DB
	}
)

func NewLoginLogModel(db *sql.DB) LoginLogModel {
	return &defaultLoginLogModel{db: db}
}

func (m *defaultLoginLogModel) tableName() string { return "`login_log`" }

func (m *defaultLoginLogModel) Insert(ctx context.Context, data *LoginLog) (sql.Result, error) {
	data.Ip = maskIP(data.Ip)
	query := fmt.Sprintf(`INSERT INTO %s (user_id, username, ip, user_agent, status, remark) VALUES (?, ?, ?, ?, ?, ?)`, m.tableName())
	return m.db.ExecContext(ctx, query, data.UserId, data.Username, data.Ip, data.UserAgent, data.Status, data.Remark)
}

func maskIP(ip string) string {
	if ip == "" {
		return ""
	}
	if strings.Contains(ip, ".") {
		parts := strings.Split(ip, ".")
		if len(parts) >= 4 {
			return parts[0] + "." + parts[1] + "." + parts[2] + ".x"
		}
	}
	if strings.Contains(ip, ":") {
		return "::1"
	}
	return "***.***.***.x"
}

func (m *defaultLoginLogModel) FindOne(ctx context.Context, id int64) (*LoginLog, error) {
	query := fmt.Sprintf(`SELECT id, user_id, username, ip, user_agent, status, remark, create_time FROM %s WHERE id = ? LIMIT 1`, m.tableName())
	var l LoginLog
	err := m.db.QueryRowContext(ctx, query, id).Scan(&l.Id, &l.UserId, &l.Username, &l.Ip, &l.UserAgent, &l.Status, &l.Remark, &l.CreateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &l, err
}

func (m *defaultLoginLogModel) Update(ctx context.Context, data *LoginLog) error {
	query := fmt.Sprintf(`UPDATE %s SET user_id = ?, username = ?, ip = ?, user_agent = ?, status = ?, remark = ? WHERE id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, data.UserId, data.Username, data.Ip, data.UserAgent, data.Status, data.Remark, data.Id)
	return err
}

func (m *defaultLoginLogModel) DeleteOlderThan(ctx context.Context, beforeTime time.Time) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE create_time < ?`, m.tableName())
	result, err := m.db.ExecContext(ctx, query, beforeTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultLoginLogModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m *defaultLoginLogModel) DeleteBatch(ctx context.Context, ids []int64) error {
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

func (m *defaultLoginLogModel) FindList(ctx context.Context, username string, page, pageSize int64) ([]*LoginLog, int64, error) {
	var where string
	var args []interface{}

	where = " WHERE 1=1"
	if username != "" {
		where += " AND username LIKE ?"
		args = append(args, "%"+username+"%")
	}

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, m.tableName()) + where
	err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf(`SELECT id, user_id, username, ip, user_agent, status, remark, create_time FROM %s`, m.tableName()) + where + ` ORDER BY id DESC LIMIT ? OFFSET ?`
	listArgs := append(args, pageSize, offset)
	rows, err := m.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*LoginLog
	for rows.Next() {
		var l LoginLog
		err := rows.Scan(&l.Id, &l.UserId, &l.Username, &l.Ip, &l.UserAgent, &l.Status, &l.Remark, &l.CreateTime)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, &l)
	}
	return list, total, rows.Err()
}
