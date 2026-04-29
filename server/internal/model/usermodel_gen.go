package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
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
		UpdateStatus(ctx context.Context, id, status int64) error
		UpdatePassword(ctx context.Context, id int64, password string) error
		IncrementFailedAttempts(ctx context.Context, id int64, maxAttempts int, lockDurationMinutes int) error
		ResetFailedAttempts(ctx context.Context, id int64) error
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
	query := fmt.Sprintf(`SELECT id, username, password, nickname, email, phone, avatar, role, status, is_deleted, create_time, update_time FROM %s WHERE id = ? AND is_deleted = 0 LIMIT 1`, m.tableName())
	var resp User
	err := m.db.QueryRowContext(ctx, query, id).Scan(&resp.Id, &resp.Username, &resp.Password, &resp.Nickname, &resp.Email, &resp.Phone, &resp.Avatar, &resp.Role, &resp.Status, &resp.IsDeleted, &resp.CreateTime, &resp.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &resp, err
}

func (m *defaultUserModel) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	query := fmt.Sprintf(`SELECT id, username, password, nickname, email, phone, avatar, role, status, is_deleted, failed_attempts, locked_until, create_time, update_time FROM %s WHERE username = ? AND is_deleted = 0 LIMIT 1`, m.tableName())
	var resp User
	err := m.db.QueryRowContext(ctx, query, username).Scan(&resp.Id, &resp.Username, &resp.Password, &resp.Nickname, &resp.Email, &resp.Phone, &resp.Avatar, &resp.Role, &resp.Status, &resp.IsDeleted, &resp.FailedAttempts, &resp.LockedUntil, &resp.CreateTime, &resp.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &resp, err
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO %s (username, password, nickname, email, phone, avatar, role, status, is_deleted, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`, m.tableName())
	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	data.Password = string(hashedPassword)

	return m.db.ExecContext(ctx, query, data.Username, data.Password, data.Nickname, data.Email, data.Phone, data.Avatar, data.Role, data.Status, data.CreateTime, data.UpdateTime)
}

func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	query := fmt.Sprintf(`UPDATE %s SET nickname = ?, email = ?, phone = ?, avatar = ?, role = ?, status = ?, update_time = ? WHERE id = ?`, m.tableName())
	data.UpdateTime = time.Now()
	_, err := m.db.ExecContext(ctx, query, data.Nickname, data.Email, data.Phone, data.Avatar, data.Role, data.Status, data.UpdateTime, data.Id)
	return err
}

func (m *defaultUserModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`UPDATE %s SET is_deleted = 1, update_time = ? WHERE id = ?`, m.tableName())
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
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, m.tableName()) + where
	err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf(`SELECT id, username, nickname, email, phone, avatar, role, status, is_deleted, create_time, update_time FROM %s`, m.tableName()) + where + ` ORDER BY id DESC LIMIT ? OFFSET ?`
	listArgs := append(args, pageSize, offset)

	rows, err := m.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Username, &u.Nickname, &u.Email, &u.Phone, &u.Avatar, &u.Role, &u.Status, &u.IsDeleted, &u.CreateTime, &u.UpdateTime)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, &u)
	}
	return list, total, rows.Err()
}

func (m *defaultUserModel) UpdateStatus(ctx context.Context, id, status int64) error {
	query := fmt.Sprintf(`UPDATE %s SET status = ?, update_time = ? WHERE id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (m *defaultUserModel) UpdatePassword(ctx context.Context, id int64, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`UPDATE %s SET password = ?, update_time = ? WHERE id = ?`, m.tableName())
	_, err = m.db.ExecContext(ctx, query, string(hashedPassword), time.Now(), id)
	return err
}

func (m *defaultUserModel) IncrementFailedAttempts(ctx context.Context, id int64, maxAttempts int, lockDurationMinutes int) error {
	query := fmt.Sprintf(`UPDATE %s SET failed_attempts = failed_attempts + 1, locked_until = CASE WHEN failed_attempts + 1 >= ? THEN ? ELSE NULL END, update_time = ? WHERE id = ? AND (locked_until IS NULL OR locked_until < ?)`, m.tableName())
	lockedUntil := time.Now().Add(time.Duration(lockDurationMinutes) * time.Minute)
	_, err := m.db.ExecContext(ctx, query, maxAttempts, lockedUntil, time.Now(), id, time.Now())
	return err
}

func (m *defaultUserModel) ResetFailedAttempts(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`UPDATE %s SET failed_attempts = 0, locked_until = NULL, update_time = ? WHERE id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
