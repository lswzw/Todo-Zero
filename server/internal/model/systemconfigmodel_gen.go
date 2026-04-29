package model

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
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
	}

	defaultSystemConfigModel struct {
		db    *sql.DB
		cache sync.Map // key -> *cacheEntry
	}

	cacheEntry struct {
		value     *SystemConfig
		expiredAt time.Time
	}
)

// cacheTTL 配置缓存有效期，更新后自动失效
const cacheTTL = 30 * time.Second

func NewSystemConfigModel(db *sql.DB) SystemConfigModel {
	return &defaultSystemConfigModel{db: db}
}

func (m *defaultSystemConfigModel) tableName() string { return "`system_configs`" }

func (m *defaultSystemConfigModel) Insert(ctx context.Context, data *SystemConfig) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO %s (config_key, config_value, group_name, description, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?)`, m.tableName())
	now := time.Now()
	data.CreateTime = now
	data.UpdateTime = now
	result, err := m.db.ExecContext(ctx, query, data.ConfigKey, data.ConfigValue, data.GroupName, data.Description, data.CreateTime, data.UpdateTime)
	if err != nil {
		return result, err
	}
	// 新增后清除缓存，确保下次读取拿到最新值
	m.cache.Delete(data.ConfigKey)
	return result, nil
}

func (m *defaultSystemConfigModel) Update(ctx context.Context, data *SystemConfig) error {
	query := fmt.Sprintf(`UPDATE %s SET config_value = ?, update_time = ? WHERE id = ?`, m.tableName())
	data.UpdateTime = time.Now()
	_, err := m.db.ExecContext(ctx, query, data.ConfigValue, data.UpdateTime, data.Id)
	if err != nil {
		return err
	}
	// 更新后清除缓存，确保下次读取拿到最新值
	m.cache.Delete(data.ConfigKey)
	return nil
}

func (m *defaultSystemConfigModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, id)
	return err
}

func (m *defaultSystemConfigModel) FindAll(ctx context.Context) ([]*SystemConfig, error) {
	query := fmt.Sprintf(`SELECT id, config_key, config_value, group_name, description, create_time, update_time FROM %s ORDER BY group_name ASC, id ASC`, m.tableName())
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
	// 先查缓存
	if v, ok := m.cache.Load(key); ok {
		entry := v.(*cacheEntry)
		if time.Now().Before(entry.expiredAt) {
			return entry.value, nil
		}
		m.cache.Delete(key)
	}

	// 缓存未命中，查数据库
	query := fmt.Sprintf(`SELECT id, config_key, config_value, group_name, description, create_time, update_time FROM %s WHERE config_key = ? LIMIT 1`, m.tableName())
	var c SystemConfig
	err := m.db.QueryRowContext(ctx, query, key).Scan(&c.Id, &c.ConfigKey, &c.ConfigValue, &c.GroupName, &c.Description, &c.CreateTime, &c.UpdateTime)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// 写入缓存
	m.cache.Store(key, &cacheEntry{
		value:     &c,
		expiredAt: time.Now().Add(cacheTTL),
	})
	return &c, nil
}

func (m *defaultSystemConfigModel) FindByGroup(ctx context.Context, group string) ([]*SystemConfig, error) {
	query := fmt.Sprintf(`SELECT id, config_key, config_value, group_name, description, create_time, update_time FROM %s WHERE group_name = ? ORDER BY id ASC`, m.tableName())
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

// ClearCache 清除所有缓存，用于配置批量更新后调用
func (m *defaultSystemConfigModel) ClearCache() {
	m.cache.Range(func(key, value interface{}) bool {
		m.cache.Delete(key)
		return true
	})
}

// CleanupExpiredCache 清理过期缓存条目
func (m *defaultSystemConfigModel) CleanupExpiredCache() {
	now := time.Now()
	m.cache.Range(func(key, value interface{}) bool {
		entry := value.(*cacheEntry)
		if now.After(entry.expiredAt) {
			m.cache.Delete(key)
		}
		return true
	})
}
