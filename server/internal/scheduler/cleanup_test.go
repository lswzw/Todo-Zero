package scheduler

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"server/internal/config"
	"server/internal/model"
	"server/internal/svc"
)

// ---- mock models for scheduler tests ----

type mockTaskModel struct {
	hardDeleteCompletedBefore  int64
	hardDeleteSoftDeletedBefore int64
	err                        error
}

func (m *mockTaskModel) Insert(ctx context.Context, data *model.Task) (sql.Result, error) { return nil, nil }
func (m *mockTaskModel) Update(ctx context.Context, data *model.Task) error                { return nil }
func (m *mockTaskModel) Delete(ctx context.Context, id int64) error                        { return nil }
func (m *mockTaskModel) FindOne(ctx context.Context, id int64) (*model.Task, error)        { return nil, model.ErrNotFound }
func (m *mockTaskModel) FindOneIncludeDeleted(ctx context.Context, id int64) (*model.Task, error) { return nil, model.ErrNotFound }
func (m *mockTaskModel) FindList(ctx context.Context, userId int64, keyword string, status, priority, categoryId, page, pageSize int64) ([]*model.Task, int64, error) {
	return nil, 0, nil
}
func (m *mockTaskModel) FindDeletedList(ctx context.Context, userId int64, page, pageSize int64) ([]*model.Task, int64, error) {
	return nil, 0, nil
}
func (m *mockTaskModel) UpdateStatus(ctx context.Context, id, status int64) error { return nil }
func (m *mockTaskModel) Restore(ctx context.Context, id int64) error              { return nil }
func (m *mockTaskModel) PermanentDelete(ctx context.Context, id int64) error      { return nil }
func (m *mockTaskModel) CountStats(ctx context.Context, userId int64) (total, todo, done, overdue int64, err error) {
	return 0, 0, 0, 0, nil
}
func (m *mockTaskModel) HardDeleteCompletedBefore(ctx context.Context, beforeTime time.Time) (int64, error) {
	return m.hardDeleteCompletedBefore, m.err
}
func (m *mockTaskModel) HardDeleteSoftDeletedBefore(ctx context.Context, beforeTime time.Time) (int64, error) {
	return m.hardDeleteSoftDeletedBefore, m.err
}

type mockOperationLogModel struct {
	deletedCount int64
	err          error
}

func (m *mockOperationLogModel) Insert(ctx context.Context, data *model.OperationLog) (sql.Result, error) { return nil, nil }
func (m *mockOperationLogModel) FindList(ctx context.Context, action, username string, page, pageSize int64) ([]*model.OperationLog, int64, error) {
	return nil, 0, nil
}
func (m *mockOperationLogModel) Count(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockOperationLogModel) DeleteById(ctx context.Context, id int64) error { return nil }
func (m *mockOperationLogModel) DeleteBatch(ctx context.Context, ids []int64) error { return nil }
func (m *mockOperationLogModel) DeleteOlderThan(ctx context.Context, beforeTime time.Time) (int64, error) {
	return m.deletedCount, m.err
}

type mockLoginLogModel struct {
	deletedCount int64
	err          error
}

func (m *mockLoginLogModel) Insert(ctx context.Context, data *model.LoginLog) (sql.Result, error) { return nil, nil }
func (m *mockLoginLogModel) FindOne(ctx context.Context, id int64) (*model.LoginLog, error)      { return nil, nil }
func (m *mockLoginLogModel) Update(ctx context.Context, data *model.LoginLog) error               { return nil }
func (m *mockLoginLogModel) Delete(ctx context.Context, id int64) error                           { return nil }
func (m *mockLoginLogModel) DeleteBatch(ctx context.Context, ids []int64) error                   { return nil }
func (m *mockLoginLogModel) FindList(ctx context.Context, username string, page, pageSize int64) ([]*model.LoginLog, int64, error) {
	return nil, 0, nil
}
func (m *mockLoginLogModel) DeleteOlderThan(ctx context.Context, beforeTime time.Time) (int64, error) {
	return m.deletedCount, m.err
}

type mockSystemConfigModel struct {
	configs map[string]string
}

func (m *mockSystemConfigModel) Insert(ctx context.Context, data *model.SystemConfig) (sql.Result, error) { return nil, nil }
func (m *mockSystemConfigModel) Update(ctx context.Context, data *model.SystemConfig) error                { return nil }
func (m *mockSystemConfigModel) Delete(ctx context.Context, id int64) error                                { return nil }
func (m *mockSystemConfigModel) FindAll(ctx context.Context) ([]*model.SystemConfig, error)                { return nil, nil }
func (m *mockSystemConfigModel) FindByKey(ctx context.Context, key string) (*model.SystemConfig, error) {
	if val, ok := m.configs[key]; ok {
		return &model.SystemConfig{ConfigKey: key, ConfigValue: val}, nil
	}
	return nil, model.ErrNotFound
}
func (m *mockSystemConfigModel) FindByGroup(ctx context.Context, group string) ([]*model.SystemConfig, error) {
	return nil, nil
}

// ---- helper ----

func newTestSvcCtx(configs map[string]string) *svc.ServiceContext {
	return &svc.ServiceContext{
		Config:             config.Config{},
		TaskModel:          &mockTaskModel{},
		OperationLogModel:  &mockOperationLogModel{},
		LoginLogModel:      &mockLoginLogModel{},
		SystemConfigModel:  &mockSystemConfigModel{configs: configs},
	}
}

// ---- tests ----

func TestGetConfigInt(t *testing.T) {
	svcCtx := newTestSvcCtx(map[string]string{
		"task_auto_delete_days": "90",
		"missing_key":          "abc",
	})

	tests := []struct {
		key        string
		defaultVal int64
		want       int64
	}{
		{"task_auto_delete_days", 0, 90},
		{"missing_key", 42, 42},
		{"nonexistent", 0, 0},
		{"missing_key", 0, 0}, // "abc" can't parse, returns default
	}

	for _, tt := range tests {
		got := getConfigInt(svcCtx, tt.key, tt.defaultVal)
		if got != tt.want {
			t.Errorf("getConfigInt(%q, %d) = %d, want %d", tt.key, tt.defaultVal, got, tt.want)
		}
	}
}

func TestRunCleanup_TaskAutoDelete(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		Config: config.Config{},
		TaskModel: &mockTaskModel{hardDeleteCompletedBefore: 5, hardDeleteSoftDeletedBefore: 3},
		OperationLogModel: &mockOperationLogModel{},
		LoginLogModel:     &mockLoginLogModel{},
		SystemConfigModel: &mockSystemConfigModel{configs: map[string]string{
			"task_auto_delete_days":      "90",
			"task_trash_retention_days":  "30",
			"log_auto_delete_days":       "0",
		}},
	}

	// Should not panic
	runCleanup(svcCtx)
}

func TestRunCleanup_LogAutoDelete(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		Config:    config.Config{},
		TaskModel: &mockTaskModel{},
		OperationLogModel: &mockOperationLogModel{deletedCount: 10},
		LoginLogModel:     &mockLoginLogModel{deletedCount: 7},
		SystemConfigModel: &mockSystemConfigModel{configs: map[string]string{
			"task_auto_delete_days":      "0",
			"task_trash_retention_days":  "30",
			"log_auto_delete_days":       "180",
		}},
	}

	// Should not panic
	runCleanup(svcCtx)
}

func TestRunCleanup_AllDisabled(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		Config:    config.Config{},
		TaskModel: &mockTaskModel{},
		OperationLogModel: &mockOperationLogModel{},
		LoginLogModel:     &mockLoginLogModel{},
		SystemConfigModel: &mockSystemConfigModel{configs: map[string]string{
			"task_auto_delete_days":      "0",
			"task_trash_retention_days":  "0",
			"log_auto_delete_days":       "0",
		}},
	}

	// All disabled, should be a no-op
	runCleanup(svcCtx)
}

func TestRunCleanup_ConfigMissing(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		Config:    config.Config{},
		TaskModel: &mockTaskModel{},
		OperationLogModel: &mockOperationLogModel{},
		LoginLogModel:     &mockLoginLogModel{},
		SystemConfigModel: &mockSystemConfigModel{configs: map[string]string{}},
	}

	// Configs missing, should use defaults (0 for task/log, 30 for trash)
	runCleanup(svcCtx)
}
