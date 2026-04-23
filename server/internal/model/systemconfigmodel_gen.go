package model

import (
	"context"
	"database/sql"
	"time"
)

var _ SystemConfigModel = (*defaultSystemConfigModel)(nil)

type (
	SystemConfigModel interface {
		Insert(ctx context.Context, data *SystemConfig) (sql.Result, error)
		Update(ctx context.Context, data *SystemConfig) error
		Delete(ctx context.Context, id int64) error
		FindAll(ctx context.Context) ([]*SystemConfig, error)
		FindByKey(ctx context.Context, key string) (*SystemConfig, error)
		FindByGroup(ctx context.Context, group string) ([]*SystemConfig, error)
		FindOneByKey(ctx context.Context, key string) (*SystemConfig, error)
	}

	defaultSystemConfigModel struct {
		db *sql.DB
	}
)

func NewSystemConfigModel(db *sql.DB) SystemConfigModel {
	return &defaultSystemConfigModel{db: db}
}

func (m *defaultSystemConfigModel) tableName() string { return "`system_configs`" }

func (m *defaultSystemConfigModel) Insert(ctx context.Context, data *SystemConfig) (sql.Result, error) {
	query := `INSERT INTO ` + m.tableName() + ` (config_key, config_value, group_name, description, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now
	return m.db.ExecContext(ctx, query, data.ConfigKey, data.ConfigValue, data.GroupName, data.Description, data.CreateTime, data.UpdateTime)
}

func (m *defaultSystemConfigModel) Update(ctx context.Context, data *SystemConfig) error {
	query := `UPDATE ` + m.tableName() + ` SET config_value = ?, update_time = ? WHERE id = ?`
	data.UpdateTime = time.Now()
	_, err := m.db.ExecContext(ctx, query, data.ConfigValue, data.UpdateTime, data.Id)
	return err
}

func (m *defaultSystemConfigModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM ` + m.tableName() + ` WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m *defaultSystemConfigModel) FindAll(ctx context.Context) ([]*SystemConfig, error) {
	query := `SELECT id, config_key, config_value, group_name, description, create_time, update_time FROM ` + m.tableName() + ` ORDER BY group_name ASC, id ASC`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*SystemConfig
	for rows.Next() {
		var c SystemConfig
		err := rows.Scan(&c.Id, &c.ConfigKey, &c.ConfigValue, &c.GroupName, &c.Description, &c.CreateTime, &c.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, rows.Err()
}

func (m *defaultSystemConfigModel) FindByKey(ctx context.Context, key string) (*SystemConfig, error) {
	query := `SELECT id, config_key, config_value, group_name, description, create_time, update_time FROM ` + m.tableName() + ` WHERE config_key = ? LIMIT 1`
	var c SystemConfig
	err := m.db.QueryRowContext(ctx, query, key).Scan(&c.Id, &c.ConfigKey, &c.ConfigValue, &c.GroupName, &c.Description, &c.CreateTime, &c.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &c, err
}

func (m *defaultSystemConfigModel) FindOneByKey(ctx context.Context, key string) (*SystemConfig, error) {
	return m.FindByKey(ctx, key)
}

func (m *defaultSystemConfigModel) FindByGroup(ctx context.Context, group string) ([]*SystemConfig, error) {
	query := `SELECT id, config_key, config_value, group_name, description, create_time, update_time FROM ` + m.tableName() + ` WHERE group_name = ? ORDER BY id ASC`
	rows, err := m.db.QueryContext(ctx, query, group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*SystemConfig
	for rows.Next() {
		var c SystemConfig
		err := rows.Scan(&c.Id, &c.ConfigKey, &c.ConfigValue, &c.GroupName, &c.Description, &c.CreateTime, &c.UpdateTime)
		if err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, rows.Err()
}
