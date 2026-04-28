package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var _ TaskTagModel = (*defaultTaskTagModel)(nil)

type (
	TaskTagModel interface {
		Insert(ctx context.Context, taskId, tagId int64) (sql.Result, error)
		DeleteByTaskId(ctx context.Context, taskId int64) error
		DeleteByTagId(ctx context.Context, tagId int64) error
		DeleteByTaskIdAndTagId(ctx context.Context, taskId, tagId int64) error
		FindByTaskId(ctx context.Context, taskId int64) ([]*TaskTag, error)
		FindByTagId(ctx context.Context, tagId int64) ([]*TaskTag, error)
		BatchInsert(ctx context.Context, taskId int64, tagIds []int64) error
	}

	defaultTaskTagModel struct {
		db *sql.DB
	}

	TaskTag struct {
		Id         int64     `json:"id"`
		TaskId     int64     `json:"taskId"`
		TagId      int64     `json:"tagId"`
		CreateTime time.Time `json:"createTime"`
	}
)

func NewTaskTagModel(db *sql.DB) TaskTagModel {
	return &defaultTaskTagModel{db: db}
}

func (m *defaultTaskTagModel) tableName() string { return "`task_tags`" }

func (m *defaultTaskTagModel) Insert(ctx context.Context, taskId, tagId int64) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT OR IGNORE INTO %s (task_id, tag_id, create_time) VALUES (?, ?, ?)`, m.tableName())
	return m.db.ExecContext(ctx, query, taskId, tagId, time.Now())
}

func (m *defaultTaskTagModel) DeleteByTaskId(ctx context.Context, taskId int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE task_id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, taskId)
	return err
}

func (m *defaultTaskTagModel) DeleteByTagId(ctx context.Context, tagId int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE tag_id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, tagId)
	return err
}

func (m *defaultTaskTagModel) DeleteByTaskIdAndTagId(ctx context.Context, taskId, tagId int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE task_id = ? AND tag_id = ?`, m.tableName())
	_, err := m.db.ExecContext(ctx, query, taskId, tagId)
	return err
}

func (m *defaultTaskTagModel) FindByTaskId(ctx context.Context, taskId int64) ([]*TaskTag, error) {
	query := fmt.Sprintf(`SELECT id, task_id, tag_id, create_time FROM %s WHERE task_id = ?`, m.tableName())
	rows, err := m.db.QueryContext(ctx, query, taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var resp []*TaskTag
	for rows.Next() {
		var item TaskTag
		if err := rows.Scan(&item.Id, &item.TaskId, &item.TagId, &item.CreateTime); err != nil {
			return nil, err
		}
		resp = append(resp, &item)
	}
	return resp, nil
}

func (m *defaultTaskTagModel) FindByTagId(ctx context.Context, tagId int64) ([]*TaskTag, error) {
	query := fmt.Sprintf(`SELECT id, task_id, tag_id, create_time FROM %s WHERE tag_id = ?`, m.tableName())
	rows, err := m.db.QueryContext(ctx, query, tagId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var resp []*TaskTag
	for rows.Next() {
		var item TaskTag
		if err := rows.Scan(&item.Id, &item.TaskId, &item.TagId, &item.CreateTime); err != nil {
			return nil, err
		}
		resp = append(resp, &item)
	}
	return resp, nil
}

func (m *defaultTaskTagModel) BatchInsert(ctx context.Context, taskId int64, tagIds []int64) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	query := fmt.Sprintf(`INSERT OR IGNORE INTO %s (task_id, tag_id, create_time) VALUES (?, ?, ?)`, m.tableName())
	now := time.Now()
	
	for _, tagId := range tagIds {
		if _, err := tx.ExecContext(ctx, query, taskId, tagId, now); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit()
}
