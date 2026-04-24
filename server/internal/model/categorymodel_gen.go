package model

import (
	"context"
	"database/sql"
	"time"
)

var _ CategoryModel = (*defaultCategoryModel)(nil)

type (
	CategoryModel interface {
		Insert(ctx context.Context, data *Category) (sql.Result, error)
		Update(ctx context.Context, data *Category) error
		Delete(ctx context.Context, id int64) error
		FindAll(ctx context.Context, userId int64) ([]*Category, error)
		FindSystem(ctx context.Context) ([]*Category, error)
		FindOne(ctx context.Context, id int64) (*Category, error)
		CountByUser(ctx context.Context, userId int64) (int64, error)
	}

	defaultCategoryModel struct {
		db *sql.DB
	}
)

func NewCategoryModel(db *sql.DB) CategoryModel {
	return &defaultCategoryModel{db: db}
}

func (m *defaultCategoryModel) tableName() string { return "`categories`" }

func (m *defaultCategoryModel) Insert(ctx context.Context, data *Category) (sql.Result, error) {
	query := `INSERT INTO ` + m.tableName() + ` (name, color, icon, sort, user_id, is_system, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now
	return m.db.ExecContext(ctx, query, data.Name, data.Color, data.Icon, data.Sort, data.UserId, data.IsSystem, data.CreateTime, data.UpdateTime)
}

func (m *defaultCategoryModel) Update(ctx context.Context, data *Category) error {
	query := `UPDATE ` + m.tableName() + ` SET name = ?, color = ?, icon = ?, sort = ?, update_time = ? WHERE id = ?`
	data.UpdateTime = time.Now()
	_, err := m.db.ExecContext(ctx, query, data.Name, data.Color, data.Icon, data.Sort, data.UpdateTime, data.Id)
	return err
}

func (m *defaultCategoryModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM ` + m.tableName() + ` WHERE id = ? AND is_system = 0`
	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m *defaultCategoryModel) FindAll(ctx context.Context, userId int64) ([]*Category, error) {
	query := `SELECT id, name, color, icon, sort, user_id, is_system, create_time, update_time FROM ` + m.tableName() + ` WHERE (user_id = ? OR is_system = 1) ORDER BY is_system DESC, sort ASC, id ASC`
	rows, err := m.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.Id, &c.Name, &c.Color, &c.Icon, &c.Sort, &c.UserId, &c.IsSystem, &c.CreateTime, &c.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, rows.Err()
}

func (m *defaultCategoryModel) FindSystem(ctx context.Context) ([]*Category, error) {
	query := `SELECT id, name, color, icon, sort, user_id, is_system, create_time, update_time FROM ` + m.tableName() + ` WHERE is_system = 1 ORDER BY sort ASC`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.Id, &c.Name, &c.Color, &c.Icon, &c.Sort, &c.UserId, &c.IsSystem, &c.CreateTime, &c.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, rows.Err()
}

func (m *defaultCategoryModel) FindOne(ctx context.Context, id int64) (*Category, error) {
	query := `SELECT id, name, color, icon, sort, user_id, is_system, create_time, update_time FROM ` + m.tableName() + ` WHERE id = ? LIMIT 1`
	var c Category
	err := m.db.QueryRowContext(ctx, query, id).Scan(&c.Id, &c.Name, &c.Color, &c.Icon, &c.Sort, &c.UserId, &c.IsSystem, &c.CreateTime, &c.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &c, err
}

func (m *defaultCategoryModel) CountByUser(ctx context.Context, userId int64) (int64, error) {
	query := `SELECT COUNT(*) FROM ` + m.tableName() + ` WHERE user_id = ?`
	var count int64
	err := m.db.QueryRowContext(ctx, query, userId).Scan(&count)
	return count, err
}
