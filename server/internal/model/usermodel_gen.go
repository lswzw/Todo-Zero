package model

import (
	"context"
	"database/sql"
	"time"
)

var _ UserModel = (*defaultUserModel)(nil)

type (
	UserModel interface {
		FindOne(ctx context.Context, id int64) (*User, error)
		FindOneByUsername(ctx context.Context, username string) (*User, error)
		Insert(ctx context.Context, data *User) (sql.Result, error)
		Update(ctx context.Context, data *User) error
		Delete(ctx context.Context, id int64) error
		FindList(ctx context.Context, username string, status, page, pageSize int64) ([]*User, int64, error)
		FindById(ctx context.Context, id int64) (*User, error)
		UpdateStatus(ctx context.Context, id, status int64) error
		UpdatePassword(ctx context.Context, id int64, password string) error
	}

	defaultUserModel struct {
		db *sql.DB
	}
)

func NewUserModel(db *sql.DB) UserModel {
	return &defaultUserModel{db: db}
}

func (m *defaultUserModel) tableName() string { return "`users`" }

func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, username, password, nickname, email, phone, avatar, role, status, is_deleted, create_time, update_time FROM ` + m.tableName() + ` WHERE id = ? AND is_deleted = 0 LIMIT 1`
	var resp User
	err := m.db.QueryRowContext(ctx, query, id).Scan(&resp.Id, &resp.Username, &resp.Password, &resp.Nickname, &resp.Email, &resp.Phone, &resp.Avatar, &resp.Role, &resp.Status, &resp.IsDeleted, &resp.CreateTime, &resp.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &resp, err
}

func (m *defaultUserModel) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, password, nickname, email, phone, avatar, role, status, is_deleted, create_time, update_time FROM ` + m.tableName() + ` WHERE username = ? AND is_deleted = 0 LIMIT 1`
	var resp User
	err := m.db.QueryRowContext(ctx, query, username).Scan(&resp.Id, &resp.Username, &resp.Password, &resp.Nickname, &resp.Email, &resp.Phone, &resp.Avatar, &resp.Role, &resp.Status, &resp.IsDeleted, &resp.CreateTime, &resp.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &resp, err
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	query := `INSERT INTO ` + m.tableName() + ` (username, password, nickname, email, phone, avatar, role, status, is_deleted, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`
	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now
	return m.db.ExecContext(ctx, query, data.Username, data.Password, data.Nickname, data.Email, data.Phone, data.Avatar, data.Role, data.Status, data.CreateTime, data.UpdateTime)
}

func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	query := `UPDATE ` + m.tableName() + ` SET nickname = ?, email = ?, phone = ?, avatar = ?, role = ?, status = ?, update_time = ? WHERE id = ?`
	data.UpdateTime = time.Now()
	_, err := m.db.ExecContext(ctx, query, data.Nickname, data.Email, data.Phone, data.Avatar, data.Role, data.Status, data.UpdateTime, data.Id)
	return err
}

func (m *defaultUserModel) Delete(ctx context.Context, id int64) error {
	query := `UPDATE ` + m.tableName() + ` SET is_deleted = 1, update_time = ? WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (m *defaultUserModel) FindList(ctx context.Context, username string, status, page, pageSize int64) ([]*User, int64, error) {
	where := " WHERE is_deleted = 0"
	var args []interface{}
	if username != "" {
		where += " AND username LIKE ?"
		args = append(args, "%"+username+"%")
	}
	if status >= 0 {
		where += " AND status = ?"
		args = append(args, status)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM ` + m.tableName() + where
	err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := `SELECT id, username, password, nickname, email, phone, avatar, role, status, is_deleted, create_time, update_time FROM ` + m.tableName() + where + ` ORDER BY id DESC LIMIT ? OFFSET ?`
	listArgs := append(args, pageSize, offset)

	rows, err := m.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.Nickname, &u.Email, &u.Phone, &u.Avatar, &u.Role, &u.Status, &u.IsDeleted, &u.CreateTime, &u.UpdateTime)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, &u)
	}
	return list, total, rows.Err()
}

func (m *defaultUserModel) FindById(ctx context.Context, id int64) (*User, error) {
	return m.FindOne(ctx, id)
}

func (m *defaultUserModel) UpdateStatus(ctx context.Context, id, status int64) error {
	query := `UPDATE ` + m.tableName() + ` SET status = ?, update_time = ? WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (m *defaultUserModel) UpdatePassword(ctx context.Context, id int64, password string) error {
	query := `UPDATE ` + m.tableName() + ` SET password = ?, update_time = ? WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, password, time.Now(), id)
	return err
}
