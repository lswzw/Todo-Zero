package model

import (
	"context"
	"database/sql"
	"strings"
	"time"
)

var _ TaskModel = (*defaultTaskModel)(nil)

type (
	TaskModel interface {
		Insert(ctx context.Context, data *Task) (sql.Result, error)
		Update(ctx context.Context, data *Task) error
		Delete(ctx context.Context, id int64) error
		FindOne(ctx context.Context, id int64) (*Task, error)
		FindByUserId(ctx context.Context, userId int64) ([]*Task, error)
		FindByCategoryId(ctx context.Context, categoryId int64) ([]*Task, error)
		CountByStatus(ctx context.Context, userId int64, status int64) (int64, error)
		BatchDelete(ctx context.Context, ids []int64) error
		FindList(ctx context.Context, userId int64, keyword string, status, priority, categoryId, page, pageSize int64) ([]*Task, int64, error)
		UpdateStatus(ctx context.Context, id, status int64) error
		CountStats(ctx context.Context, userId int64) (total, todo, done, overdue int64, err error)
	}

	defaultTaskModel struct {
		db *sql.DB
	}
)

func NewTaskModel(db *sql.DB) TaskModel {
	return &defaultTaskModel{db: db}
}

func (m *defaultTaskModel) tableName() string { return "`tasks`" }

func (m *defaultTaskModel) Insert(ctx context.Context, data *Task) (sql.Result, error) {
	query := `INSERT INTO ` + m.tableName() + ` (title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, is_deleted, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`
	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now
	return m.db.ExecContext(ctx, query, data.Title, data.Content, data.Priority, data.Status, data.CategoryId, data.UserId, data.StartTime, data.EndTime, data.Reminder, data.Tags, data.CreateTime, data.UpdateTime)
}

func (m *defaultTaskModel) Update(ctx context.Context, data *Task) error {
	query := `UPDATE ` + m.tableName() + ` SET title = ?, content = ?, priority = ?, status = ?, category_id = ?, start_time = ?, end_time = ?, reminder = ?, tags = ?, update_time = ? WHERE id = ?`
	data.UpdateTime = time.Now()
	_, err := m.db.ExecContext(ctx, query, data.Title, data.Content, data.Priority, data.Status, data.CategoryId, data.StartTime, data.EndTime, data.Reminder, data.Tags, data.UpdateTime, data.Id)
	return err
}

func (m *defaultTaskModel) Delete(ctx context.Context, id int64) error {
	query := `UPDATE ` + m.tableName() + ` SET is_deleted = 1, update_time = ? WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (m *defaultTaskModel) FindOne(ctx context.Context, id int64) (*Task, error) {
	query := `SELECT id, title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, is_deleted, create_time, update_time FROM ` + m.tableName() + ` WHERE id = ? AND is_deleted = 0 LIMIT 1`
	var resp Task
	err := m.db.QueryRowContext(ctx, query, id).Scan(&resp.Id, &resp.Title, &resp.Content, &resp.Priority, &resp.Status, &resp.CategoryId, &resp.UserId, &resp.StartTime, &resp.EndTime, &resp.Reminder, &resp.Tags, &resp.IsDeleted, &resp.CreateTime, &resp.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &resp, err
}

func (m *defaultTaskModel) FindByUserId(ctx context.Context, userId int64) ([]*Task, error) {
	query := `SELECT id, title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, is_deleted, create_time, update_time FROM ` + m.tableName() + ` WHERE user_id = ? AND is_deleted = 0 ORDER BY id DESC`
	rows, err := m.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Content, &t.Priority, &t.Status, &t.CategoryId, &t.UserId, &t.StartTime, &t.EndTime, &t.Reminder, &t.Tags, &t.IsDeleted, &t.CreateTime, &t.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	return list, rows.Err()
}

func (m *defaultTaskModel) FindByCategoryId(ctx context.Context, categoryId int64) ([]*Task, error) {
	query := `SELECT id, title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, is_deleted, create_time, update_time FROM ` + m.tableName() + ` WHERE category_id = ? AND is_deleted = 0 ORDER BY id DESC`
	rows, err := m.db.QueryContext(ctx, query, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Content, &t.Priority, &t.Status, &t.CategoryId, &t.UserId, &t.StartTime, &t.EndTime, &t.Reminder, &t.Tags, &t.IsDeleted, &t.CreateTime, &t.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	return list, rows.Err()
}

func (m *defaultTaskModel) CountByStatus(ctx context.Context, userId int64, status int64) (int64, error) {
	query := `SELECT COUNT(*) FROM ` + m.tableName() + ` WHERE user_id = ? AND status = ? AND is_deleted = 0`
	var count int64
	err := m.db.QueryRowContext(ctx, query, userId, status).Scan(&count)
	return count, err
}

func (m *defaultTaskModel) BatchDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `UPDATE ` + m.tableName() + ` SET is_deleted = 1, update_time = ? WHERE id IN (` + strings.Join(placeholders, ",") + `)`
	args = append([]interface{}{time.Now()}, args...)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

func (m *defaultTaskModel) FindList(ctx context.Context, userId int64, keyword string, status, priority, categoryId, page, pageSize int64) ([]*Task, int64, error) {
	where := " WHERE user_id = ? AND is_deleted = 0"
	args := []interface{}{userId}
	if keyword != "" {
		where += " AND (title LIKE ? OR content LIKE ?)"
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}
	if status >= 0 {
		where += " AND status = ?"
		args = append(args, status)
	}
	if priority >= 0 {
		where += " AND priority = ?"
		args = append(args, priority)
	}
	if categoryId > 0 {
		where += " AND category_id = ?"
		args = append(args, categoryId)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM ` + m.tableName() + where
	err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := `SELECT id, title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, is_deleted, create_time, update_time FROM ` + m.tableName() + where + ` ORDER BY id DESC LIMIT ? OFFSET ?`
	listArgs := append(args, pageSize, offset)
	rows, err := m.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Content, &t.Priority, &t.Status, &t.CategoryId, &t.UserId, &t.StartTime, &t.EndTime, &t.Reminder, &t.Tags, &t.IsDeleted, &t.CreateTime, &t.UpdateTime)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, &t)
	}
	return list, total, rows.Err()
}

func (m *defaultTaskModel) UpdateStatus(ctx context.Context, id, status int64) error {
	query := `UPDATE ` + m.tableName() + ` SET status = ?, update_time = ? WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (m *defaultTaskModel) CountStats(ctx context.Context, userId int64) (total, todo, done, overdue int64, err error) {
	if err = m.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+m.tableName()+` WHERE user_id = ? AND is_deleted = 0`, userId).Scan(&total); err != nil {
		return
	}
	if err = m.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+m.tableName()+` WHERE user_id = ? AND status != 2 AND is_deleted = 0`, userId).Scan(&todo); err != nil {
		return
	}
	if err = m.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+m.tableName()+` WHERE user_id = ? AND status = 2 AND is_deleted = 0`, userId).Scan(&done); err != nil {
		return
	}
	if err = m.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+m.tableName()+` WHERE user_id = ? AND end_time < ? AND status != 2 AND is_deleted = 0`, userId, time.Now()).Scan(&overdue); err != nil {
		return
	}
	return
}
