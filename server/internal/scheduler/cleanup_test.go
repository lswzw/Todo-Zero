package scheduler

import (
	"context"
	"testing"
	"time"

	"server/internal/config"
	"server/internal/model"
	"server/internal/svc"
)

// ---- test helpers (internal use only) ----

type mockTaskModel struct {
	model.TaskModel
	hardDeleteCompletedBefore   int64
	hardDeleteSoftDeletedBefore int64
	err                         error
}

func (m *mockTaskModel) HardDeleteCompletedBefore(ctx context.Context, beforeTime time.Time) (int64, error) {
	return m.hardDeleteCompletedBefore, m.err
}

func (m *mockTaskModel) HardDeleteSoftDeletedBefore(ctx context.Context, beforeTime time.Time) (int64, error) {
	return m.hardDeleteSoftDeletedBefore, m.err
}

type mockOperationLogModel struct {
	model.OperationLogModel
	deletedCount int64
	err          error
}

func (m *mockOperationLogModel) DeleteOlderThan(ctx context.Context, beforeTime time.Time) (int64, error) {
	return m.deletedCount, m.err
}

type mockLoginLogModel struct {
	model.LoginLogModel
	deletedCount int64
	err          error
}

func (m *mockLoginLogModel) DeleteOlderThan(ctx context.Context, beforeTime time.Time) (int64, error) {
	return m.deletedCount, m.err
}

type mockSystemConfigModel struct {
	model.SystemConfigModel
	configs map[string]string
}

func (m *mockSystemConfigModel) FindByKey(ctx context.Context, key string) (*model.SystemConfig, error) {
	if val, ok := m.configs[key]; ok {
		return &model.SystemConfig{ConfigKey: key, ConfigValue: val}, nil
	}
	return nil, model.ErrNotFound
}

// ---- helper ----

func newTestSvcCtx(configs map[string]string) *svc.ServiceContext {
	return &svc.ServiceContext{
		Config:            config.Config{},
		TaskModel:         &mockTaskModel{},
		OperationLogModel: &mockOperationLogModel{},
		LoginLogModel:     &mockLoginLogModel{},
		SystemConfigModel: &mockSystemConfigModel{configs: configs},
	}
}

// ---- tests ----

func TestGetConfigInt(t *testing.T) {
	svcCtx := newTestSvcCtx(map[string]string{
		"task_auto_delete_days": "90",
		"missing_key":           "abc",
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
		Config:            config.Config{},
		TaskModel:         &mockTaskModel{hardDeleteCompletedBefore: 5, hardDeleteSoftDeletedBefore: 3},
		OperationLogModel: &mockOperationLogModel{},
		LoginLogModel:     &mockLoginLogModel{},
		SystemConfigModel: &mockSystemConfigModel{configs: map[string]string{
			"task_auto_delete_days":     "90",
			"task_trash_retention_days": "30",
			"log_auto_delete_days":      "0",
		}},
	}

	// Should not panic
	runCleanup(svcCtx)
}

func TestRunCleanup_LogAutoDelete(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		Config:            config.Config{},
		TaskModel:         &mockTaskModel{},
		OperationLogModel: &mockOperationLogModel{deletedCount: 10},
		LoginLogModel:     &mockLoginLogModel{deletedCount: 7},
		SystemConfigModel: &mockSystemConfigModel{configs: map[string]string{
			"task_auto_delete_days":     "0",
			"task_trash_retention_days": "30",
			"log_auto_delete_days":      "180",
		}},
	}

	// Should not panic
	runCleanup(svcCtx)
}

func TestRunCleanup_AllDisabled(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		Config:            config.Config{},
		TaskModel:         &mockTaskModel{},
		OperationLogModel: &mockOperationLogModel{},
		LoginLogModel:     &mockLoginLogModel{},
		SystemConfigModel: &mockSystemConfigModel{configs: map[string]string{
			"task_auto_delete_days":     "0",
			"task_trash_retention_days": "0",
			"log_auto_delete_days":      "0",
		}},
	}

	// All disabled, should be a no-op
	runCleanup(svcCtx)
}

func TestRunCleanup_ConfigMissing(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		Config:            config.Config{},
		TaskModel:         &mockTaskModel{},
		OperationLogModel: &mockOperationLogModel{},
		LoginLogModel:     &mockLoginLogModel{},
		SystemConfigModel: &mockSystemConfigModel{configs: map[string]string{}},
	}

	// Configs missing, should use defaults (0 for task/log, 30 for trash)
	runCleanup(svcCtx)
}
