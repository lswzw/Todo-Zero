package model

import (
	"context"
	"database/sql"
)

var _ LoginLogModel = (*defaultLoginLogModel)(nil)

type (
	LoginLogModel interface {
		Insert(ctx context.Context, data *LoginLog) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*LoginLog, error)
		Update(ctx context.Context, data *LoginLog) error
		Delete(ctx context.Context, id int64) error
		FindList(ctx context.Context, username string, page, pageSize int64) ([]*LoginLog, int64, error)
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
	query := `INSERT INTO ` + m.tableName() + ` (user_id, username, ip, user_agent, status, remark) VALUES (?, ?, ?, ?, ?, ?)`
	return m.db.ExecContext(ctx, query, data.UserId, data.Username, data.Ip, data.UserAgent, data.Status, data.Remark)
}

func (m *defaultLoginLogModel) FindOne(ctx context.Context, id int64) (*LoginLog, error) {
	query := `SELECT id, user_id, username, ip, user_agent, status, remark, create_time FROM ` + m.tableName() + ` WHERE id = ? LIMIT 1`
	var l LoginLog
	err := m.db.QueryRowContext(ctx, query, id).Scan(&l.Id, &l.UserId, &l.Username, &l.Ip, &l.UserAgent, &l.Status, &l.Remark, &l.CreateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &l, err
}

func (m *defaultLoginLogModel) Update(ctx context.Context, data *LoginLog) error {
	query := `UPDATE ` + m.tableName() + ` SET user_id = ?, username = ?, ip = ?, user_agent = ?, status = ?, remark = ? WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, data.UserId, data.Username, data.Ip, data.UserAgent, data.Status, data.Remark, data.Id)
	return err
}

func (m *defaultLoginLogModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM ` + m.tableName() + ` WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, id)
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
	countQuery := `SELECT COUNT(*) FROM ` + m.tableName() + where
	err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := `SELECT id, user_id, username, ip, user_agent, status, remark, create_time FROM ` + m.tableName() + where + ` ORDER BY id DESC LIMIT ? OFFSET ?`
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
