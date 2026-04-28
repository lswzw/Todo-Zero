package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var _ TagModel = (*defaultTagModel)(nil)

type (
	TagModel interface {
		Insert(ctx context.Context, data *Tag) (sql.Result, error)
		Update(ctx context.Context, data *Tag) error
		Delete(ctx context.Context, id int64) error
		FindOne(ctx context.Context, id int64) (*Tag, error)
		FindList(ctx context.Context, userId int64, keyword string) ([]*Tag, error)
		FindByName(ctx context.Context, userId int64, name string) (*Tag, error)
	}

	defaultTagModel struct {
		db *sql.DB
	}

	Tag struct {
		Id       int64     `json:"id"`
		Name     string    `json:"name"`
		Color    string    `json:"color"`
		UserId   int64     `json:"userId"`
		IsSystem int64     `json:"isSystem"`
		CreateTime time.Time `json:"createTime"`
		UpdateTime time.Time `json:"updateTime"`
	}
)

func NewTagModel(db *sql.DB) TagModel {
	return &defaultTagModel{db: db}
}

func (m *defaultTagModel) tableName() string { return "`tags`" }

func (m *defaultTagModel) Insert(ctx context.Context, data *Tag) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO %s (name, color, user_id, is_system, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?)`, m.tableName())
	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now
	return m.db.ExecContext(ctx, query, data.Name, data.Color, data.UserId, data.IsSystem, data.CreateTime, data.UpdateTime)
}

func (m *defaultTagModel) Update(ctx context.Context, data *Tag) error {
	query := fmt.Sprintf(`UPDATE %s SET name = ?, color = ?, update_time = ? WHERE id = ? AND user_id = ?`, m.tableName())
	data.UpdateTime = time.Now()
	_, err := m.db.ExecContext(ctx, query, data.Name, data.Color, data.UpdateTime, data.Id, data.UserId)
	return err
}

func (m *defaultTagModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m *defaultTagModel) FindOne(ctx context.Context, id int64) (*Tag, error) {
	query := fmt.Sprintf(`SELECT id, name, color, user_id, is_system, create_time, update_time FROM %s WHERE id = ? LIMIT 1`, m.tableName())
	var resp Tag
	err := m.db.QueryRowContext(ctx, query, id).Scan(&resp.Id, &resp.Name, &resp.Color, &resp.UserId, &resp.IsSystem, &resp.CreateTime, &resp.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &resp, err
}

func (m *defaultTagModel) FindList(ctx context.Context, userId int64, keyword string) ([]*Tag, error) {
	query := fmt.Sprintf(`SELECT id, name, color, user_id, is_system, create_time, update_time FROM %s WHERE user_id = ?`, m.tableName())
	args := []interface{}{userId}
	
	if keyword != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+keyword+"%")
	}
	
	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var resp []*Tag
	for rows.Next() {
		var item Tag
		if err := rows.Scan(&item.Id, &item.Name, &item.Color, &item.UserId, &item.IsSystem, &item.CreateTime, &item.UpdateTime); err != nil {
			return nil, err
		}
		resp = append(resp, &item)
	}
	return resp, nil
}

func (m *defaultTagModel) FindByName(ctx context.Context, userId int64, name string) (*Tag, error) {
	query := fmt.Sprintf(`SELECT id, name, color, user_id, is_system, create_time, update_time FROM %s WHERE user_id = ? AND name = ? LIMIT 1`, m.tableName())
	var resp Tag
	err := m.db.QueryRowContext(ctx, query, userId, name).Scan(&resp.Id, &resp.Name, &resp.Color, &resp.UserId, &resp.IsSystem, &resp.CreateTime, &resp.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &resp, err
}
