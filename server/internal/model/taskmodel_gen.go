package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var _ TaskModel = (*defaultTaskModel)(nil)

type (
	TaskModel interface {
		Insert(ctx context.Context, data *Task) (sql.Result, error)
		Update(ctx context.Context, data *Task) error
		Delete(ctx context.Context, id int64) error
		FindOne(ctx context.Context, id int64) (*Task, error)
		FindOneIncludeDeleted(ctx context.Context, id int64) (*Task, error)
		FindList(ctx context.Context, userId int64, keyword string, status, priority, categoryId, page, pageSize int64) ([]*Task, int64, error)
		FindAllForExport(ctx context.Context, userId int64, keyword string, status, priority, categoryId int64) ([]*Task, error)
		FindDeletedList(ctx context.Context, userId int64, page, pageSize int64) ([]*Task, int64, error)
		UpdateStatus(ctx context.Context, id, status int64) error
		UpdateSortOrder(ctx context.Context, userId int64, orders []SortOrderItem) error
		Restore(ctx context.Context, id int64) error
		PermanentDelete(ctx context.Context, id int64) error
		CountStats(ctx context.Context, userId int64) (total, todo, done, overdue int64, err error)
		HardDeleteCompletedBefore(ctx context.Context, beforeTime time.Time) (int64, error)
		HardDeleteSoftDeletedBefore(ctx context.Context, beforeTime time.Time) (int64, error)
	}

	defaultTaskModel struct {
		db *sql.DB
	}
)

func NewTaskModel(db *sql.DB) TaskModel {
	return &defaultTaskModel{db: db}
}

// SortOrderItem 排序项
type SortOrderItem struct {
	Id        int64 `json:"id"`
	SortOrder int64 `json:"sortOrder"`
}

func (m *defaultTaskModel) tableName() string { return "`tasks`" }

func (m *defaultTaskModel) Insert(ctx context.Context, data *Task) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO %s (title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, sort_order, is_deleted, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`, m.tableName())
	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now
	return m.db.ExecContext(ctx, query, data.Title, data.Content, data.Priority, data.Status, data.CategoryId, data.UserId, data.StartTime, data.EndTime, data.Reminder, data.Tags, data.SortOrder, data.CreateTime, data.UpdateTime)
}

func (m *defaultTaskModel) Update(ctx context.Context, data *Task) error {
	query := fmt.Sprintf(`UPDATE %s SET title = ?, content = ?, priority = ?, status = ?, category_id = ?, start_time = ?, end_time = ?, reminder = ?, tags = ?, sort_order = ?, update_time = ? WHERE id = ?`, m.tableName())
	data.UpdateTime = time.Now()
	_, err := m.db.ExecContext(ctx, query, data.Title, data.Content, data.Priority, data.Status, data.CategoryId, data.StartTime, data.EndTime, data.Reminder, data.Tags, data.SortOrder, data.UpdateTime, data.Id)
	return err
}

func (m *defaultTaskModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`UPDATE %s SET is_deleted = 1, update_time = ? WHERE id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (m *defaultTaskModel) FindOne(ctx context.Context, id int64) (*Task, error) {
	query := fmt.Sprintf(`SELECT id, title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, sort_order, is_deleted, create_time, update_time FROM %s WHERE id = ? AND is_deleted = 0 LIMIT 1`, m.tableName())
	var resp Task
	err := m.db.QueryRowContext(ctx, query, id).Scan(&resp.Id, &resp.Title, &resp.Content, &resp.Priority, &resp.Status, &resp.CategoryId, &resp.UserId, &resp.StartTime, &resp.EndTime, &resp.Reminder, &resp.Tags, &resp.SortOrder, &resp.IsDeleted, &resp.CreateTime, &resp.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &resp, err
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
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, m.tableName()) + where
	err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf(`SELECT id, title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, sort_order, is_deleted, create_time, update_time FROM %s`, m.tableName()) + where + ` ORDER BY sort_order ASC, id DESC LIMIT ? OFFSET ?`
	listArgs := append(args, pageSize, offset)
	rows, err := m.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Content, &t.Priority, &t.Status, &t.CategoryId, &t.UserId, &t.StartTime, &t.EndTime, &t.Reminder, &t.Tags, &t.SortOrder, &t.IsDeleted, &t.CreateTime, &t.UpdateTime)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, &t)
	}
	return list, total, rows.Err()
}

func (m *defaultTaskModel) FindAllForExport(ctx context.Context, userId int64, keyword string, status, priority, categoryId int64) ([]*Task, error) {
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

	listQuery := fmt.Sprintf(`SELECT id, title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, sort_order, is_deleted, create_time, update_time FROM %s`, m.tableName()) + where + ` ORDER BY sort_order ASC, id DESC`
	rows, err := m.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Content, &t.Priority, &t.Status, &t.CategoryId, &t.UserId, &t.StartTime, &t.EndTime, &t.Reminder, &t.Tags, &t.SortOrder, &t.IsDeleted, &t.CreateTime, &t.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	return list, rows.Err()
}

func (m *defaultTaskModel) UpdateStatus(ctx context.Context, id, status int64) error {
	query := fmt.Sprintf(`UPDATE %s SET status = ?, update_time = ? WHERE id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (m *defaultTaskModel) FindOneIncludeDeleted(ctx context.Context, id int64) (*Task, error) {
	query := fmt.Sprintf(`SELECT id, title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, sort_order, is_deleted, create_time, update_time FROM %s WHERE id = ? LIMIT 1`, m.tableName())
	var resp Task
	err := m.db.QueryRowContext(ctx, query, id).Scan(&resp.Id, &resp.Title, &resp.Content, &resp.Priority, &resp.Status, &resp.CategoryId, &resp.UserId, &resp.StartTime, &resp.EndTime, &resp.Reminder, &resp.Tags, &resp.SortOrder, &resp.IsDeleted, &resp.CreateTime, &resp.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &resp, err
}

func (m *defaultTaskModel) FindDeletedList(ctx context.Context, userId int64, page, pageSize int64) ([]*Task, int64, error) {
	where := " WHERE user_id = ? AND is_deleted = 1"
	args := []interface{}{userId}

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, m.tableName()) + where
	err := m.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	listQuery := fmt.Sprintf(`SELECT id, title, content, priority, status, category_id, user_id, start_time, end_time, reminder, tags, sort_order, is_deleted, create_time, update_time FROM %s`, m.tableName()) + where + ` ORDER BY update_time DESC LIMIT ? OFFSET ?`
	listArgs := append(args, pageSize, offset)
	rows, err := m.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Content, &t.Priority, &t.Status, &t.CategoryId, &t.UserId, &t.StartTime, &t.EndTime, &t.Reminder, &t.Tags, &t.SortOrder, &t.IsDeleted, &t.CreateTime, &t.UpdateTime)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, &t)
	}
	return list, total, rows.Err()
}

func (m *defaultTaskModel) Restore(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`UPDATE %s SET is_deleted = 0, update_time = ? WHERE id = ? AND is_deleted = 1`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (m *defaultTaskModel) PermanentDelete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ? AND is_deleted = 1`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m *defaultTaskModel) HardDeleteCompletedBefore(ctx context.Context, beforeTime time.Time) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE status = 2 AND is_deleted = 0 AND update_time < ?`, m.tableName())
	result, err := m.db.ExecContext(ctx, query, beforeTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultTaskModel) HardDeleteSoftDeletedBefore(ctx context.Context, beforeTime time.Time) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE is_deleted = 1 AND update_time < ?`, m.tableName())
	result, err := m.db.ExecContext(ctx, query, beforeTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m *defaultTaskModel) CountStats(ctx context.Context, userId int64) (total, todo, done, overdue int64, err error) {
	if err = m.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = ? AND is_deleted = 0`, m.tableName()), userId).Scan(&total); err != nil {
		return
	}
	if err = m.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = ? AND status != 2 AND is_deleted = 0`, m.tableName()), userId).Scan(&todo); err != nil {
		return
	}
	if err = m.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = ? AND status = 2 AND is_deleted = 0`, m.tableName()), userId).Scan(&done); err != nil {
		return
	}
	if err = m.db.QueryRowContext(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = ? AND end_time < ? AND status != 2 AND is_deleted = 0`, m.tableName()), userId, time.Now()).Scan(&overdue); err != nil {
		return
	}
	return
}

func (m *defaultTaskModel) UpdateSortOrder(ctx context.Context, userId int64, orders []SortOrderItem) error {
	query := fmt.Sprintf(`UPDATE %s SET sort_order = ?, update_time = ? WHERE id = ? AND user_id = ? AND is_deleted = 0`, m.tableName())
	now := time.Now()
	for _, item := range orders {
		_, err := m.db.ExecContext(ctx, query, item.SortOrder, now, item.Id, userId)
		if err != nil {
			return err
		}
	}
	return nil
}
