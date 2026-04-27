package task

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"server/internal/model"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ---- mock models ----

type mockTaskModel struct {
	task *model.Task
	err  error
}

func (m *mockTaskModel) Insert(ctx context.Context, data *model.Task) (sql.Result, error) { return nil, nil }
func (m *mockTaskModel) Update(ctx context.Context, data *model.Task) error                { return nil }
func (m *mockTaskModel) Delete(ctx context.Context, id int64) error                        { return nil }
func (m *mockTaskModel) FindOne(ctx context.Context, id int64) (*model.Task, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.task, nil
}
func (m *mockTaskModel) FindOneIncludeDeleted(ctx context.Context, id int64) (*model.Task, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.task, nil
}
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
	return 0, nil
}
func (m *mockTaskModel) HardDeleteSoftDeletedBefore(ctx context.Context, beforeTime time.Time) (int64, error) {
	return 0, nil
}

type mockCategoryModel struct {
	category *model.Category
	err      error
}

func (m *mockCategoryModel) Insert(ctx context.Context, data *model.Category) (sql.Result, error) { return nil, nil }
func (m *mockCategoryModel) Update(ctx context.Context, data *model.Category) error                { return nil }
func (m *mockCategoryModel) Delete(ctx context.Context, id int64) error                            { return nil }
func (m *mockCategoryModel) FindAll(ctx context.Context, userId int64) ([]*model.Category, error)  { return nil, nil }
func (m *mockCategoryModel) FindSystem(ctx context.Context) ([]*model.Category, error)             { return nil, nil }
func (m *mockCategoryModel) CountByUser(ctx context.Context, userId int64) (int64, error)          { return 0, nil }
func (m *mockCategoryModel) FindOne(ctx context.Context, id int64) (*model.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.category, nil
}

// ctxWithUserId 创建包含 userId 的 context（模拟 JWT 中间件注入）
func ctxWithUserId(userId int64) context.Context {
	return context.WithValue(context.Background(), "userId", userId)
}

func TestTaskDetail_Success(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		TaskModel: &mockTaskModel{
			task: &model.Task{
				Id:         1,
				Title:      "Test Task",
				Content:    sql.NullString{String: "Content here", Valid: true},
				Status:     0,
				Priority:   1,
				CategoryId: sql.NullInt64{Int64: 10, Valid: true},
				UserId:     42,
				Tags:       "go",
			},
		},
		CategoryModel: &mockCategoryModel{
			category: &model.Category{Id: 10, Name: "Work"},
		},
	}

	logic := NewTaskDetailLogic(ctxWithUserId(42), svcCtx)
	resp, err := logic.TaskDetail(&types.TaskDetailReq{Id: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Title != "Test Task" {
		t.Errorf("expected title Test Task, got %s", resp.Title)
	}
	if resp.Content != "Content here" {
		t.Errorf("expected content, got %s", resp.Content)
	}
	if resp.CategoryName != "Work" {
		t.Errorf("expected category Work, got %s", resp.CategoryName)
	}
	if resp.Tags != "go" {
		t.Errorf("expected tags go, got %s", resp.Tags)
	}
}

func TestTaskDetail_NoPermission(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		TaskModel: &mockTaskModel{
			task: &model.Task{
				Id:     1,
				Title:  "Secret",
				UserId: 42,
			},
		},
		CategoryModel: &mockCategoryModel{},
	}

	// 用 userId=99 去访问 userId=42 的任务
	logic := NewTaskDetailLogic(ctxWithUserId(99), svcCtx)
	_, err := logic.TaskDetail(&types.TaskDetailReq{Id: 1})

	if err == nil {
		t.Fatal("expected permission error, got nil")
	}
	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		t.Fatalf("expected *xerr.CodeError, got %T", err)
	}
	if codeErr.Code != xerr.NoPermission {
		t.Errorf("expected code %d, got %d", xerr.NoPermission, codeErr.Code)
	}
}

func TestTaskDetail_NotFound(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		TaskModel:     &mockTaskModel{err: model.ErrNotFound},
		CategoryModel: &mockCategoryModel{},
	}

	logic := NewTaskDetailLogic(ctxWithUserId(1), svcCtx)
	_, err := logic.TaskDetail(&types.TaskDetailReq{Id: 999})

	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		t.Fatalf("expected *xerr.CodeError, got %T", err)
	}
	if codeErr.Code != xerr.TaskNotFoundError {
		t.Errorf("expected code %d, got %d", xerr.TaskNotFoundError, codeErr.Code)
	}
}

func TestTaskDetail_NoUserId(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		TaskModel:     &mockTaskModel{},
		CategoryModel: &mockCategoryModel{},
	}

	// context 中没有 userId
	logic := NewTaskDetailLogic(context.Background(), svcCtx)
	_, err := logic.TaskDetail(&types.TaskDetailReq{Id: 1})

	if err == nil {
		t.Fatal("expected permission error for missing userId, got nil")
	}
}

func TestTaskDetail_NoCategory(t *testing.T) {
	svcCtx := &svc.ServiceContext{
		TaskModel: &mockTaskModel{
			task: &model.Task{
				Id:         1,
				Title:      "No Category",
				UserId:     1,
				CategoryId: sql.NullInt64{Valid: false},
			},
		},
		CategoryModel: &mockCategoryModel{},
	}

	logic := NewTaskDetailLogic(ctxWithUserId(1), svcCtx)
	resp, err := logic.TaskDetail(&types.TaskDetailReq{Id: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.CategoryName != "未分类" {
		t.Errorf("expected category 未分类, got %s", resp.CategoryName)
	}
	if resp.CategoryId != 0 {
		t.Errorf("expected categoryId 0, got %d", resp.CategoryId)
	}
}

// 确保编译时 logx.Logger 被使用（避免 unused import）
var _ logx.Logger
