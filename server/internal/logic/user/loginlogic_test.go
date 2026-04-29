package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"server/internal/config"
	"server/internal/model"
	"server/internal/pkg/xerr"
	"server/internal/svc"
	"server/internal/types"

	"golang.org/x/crypto/bcrypt"
)

// ---- mock models ----

type mockUserModel struct {
	user *model.User
	err  error
}

func (m *mockUserModel) FindOne(ctx context.Context, id int64) (*model.User, error)         { return m.user, nil }
func (m *mockUserModel) FindOneByUsername(ctx context.Context, username string) (*model.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}
func (m *mockUserModel) Insert(ctx context.Context, data *model.User) (sql.Result, error)    { return nil, nil }
func (m *mockUserModel) Update(ctx context.Context, data *model.User) error                   { return nil }
func (m *mockUserModel) Delete(ctx context.Context, id int64) error                           { return nil }
func (m *mockUserModel) FindList(ctx context.Context, username string, status, page, pageSize int64) ([]*model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserModel) UpdateStatus(ctx context.Context, id, status int64) error      { return nil }
func (m *mockUserModel) UpdatePassword(ctx context.Context, id int64, password string) error { return nil }
func (m *mockUserModel) IncrementFailedAttempts(ctx context.Context, id int64, maxAttempts int, lockDurationMinutes int) error {
	return nil
}
func (m *mockUserModel) ResetFailedAttempts(ctx context.Context, id int64) error { return nil }

type mockLoginLogModel struct{}

func (m *mockLoginLogModel) Insert(ctx context.Context, data *model.LoginLog) (sql.Result, error) {
	return nil, nil
}
func (m *mockLoginLogModel) FindOne(ctx context.Context, id int64) (*model.LoginLog, error) { return nil, nil }
func (m *mockLoginLogModel) Update(ctx context.Context, data *model.LoginLog) error          { return nil }
func (m *mockLoginLogModel) Delete(ctx context.Context, id int64) error                      { return nil }
func (m *mockLoginLogModel) DeleteBatch(ctx context.Context, ids []int64) error                { return nil }
func (m *mockLoginLogModel) FindList(ctx context.Context, username string, page, pageSize int64) ([]*model.LoginLog, int64, error) {
	return nil, 0, nil
}
func (m *mockLoginLogModel) DeleteOlderThan(ctx context.Context, beforeTime time.Time) (int64, error) {
	return 0, nil
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes)
}

func testSvcCtx(user *model.User, userErr error) *svc.ServiceContext {
	c := config.Config{}
	c.Auth.AccessSecret = "test-secret-key-for-unit-test"
	c.Auth.AccessExpire = 3600

	return &svc.ServiceContext{
		Config:        c,
		UserModel:     &mockUserModel{user: user, err: userErr},
		LoginLogModel: &mockLoginLogModel{},
	}
}

func TestLogin_Success(t *testing.T) {
	svcCtx := testSvcCtx(&model.User{
		Id:       1,
		Username: "alice",
		Password: hashPassword("pass123"),
		Role:     0,
		Status:   1,
	}, nil)

	logic := NewLoginLogic(context.Background(), svcCtx)
	resp, err := logic.Login(&types.LoginReq{Username: "alice", Password: "pass123"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	if resp.IsAdmin != 0 {
		t.Errorf("expected isAdmin 0, got %d", resp.IsAdmin)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	svcCtx := testSvcCtx(nil, model.ErrNotFound)

	logic := NewLoginLogic(context.Background(), svcCtx)
	_, err := logic.Login(&types.LoginReq{Username: "nobody", Password: "pass"})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		t.Fatalf("expected *xerr.CodeError, got %T", err)
	}
	// 不再返回 UserNotFoundError，统一返回 UserOrPasswordError 防止用户名枚举
	if codeErr.Code != xerr.UserOrPasswordError {
		t.Errorf("expected code %d, got %d", xerr.UserOrPasswordError, codeErr.Code)
	}
}

func TestLogin_UserDisabled(t *testing.T) {
	svcCtx := testSvcCtx(&model.User{
		Id:       1,
		Username: "disabled",
		Password: hashPassword("pass"),
		Role:     0,
		Status:   0, // 禁用
	}, nil)

	logic := NewLoginLogic(context.Background(), svcCtx)
	_, err := logic.Login(&types.LoginReq{Username: "disabled", Password: "pass"})

	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		t.Fatalf("expected *xerr.CodeError, got %T", err)
	}
	// 禁用用户也统一返回 UserOrPasswordError，防止信息泄露
	if codeErr.Code != xerr.UserOrPasswordError {
		t.Errorf("expected code %d, got %d", xerr.UserOrPasswordError, codeErr.Code)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svcCtx := testSvcCtx(&model.User{
		Id:       1,
		Username: "alice",
		Password: hashPassword("correct"),
		Role:     0,
		Status:   1,
	}, nil)

	logic := NewLoginLogic(context.Background(), svcCtx)
	_, err := logic.Login(&types.LoginReq{Username: "alice", Password: "wrong"})

	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		t.Fatalf("expected *xerr.CodeError, got %T", err)
	}
	// 密码错误统一返回 UserOrPasswordError
	if codeErr.Code != xerr.UserOrPasswordError {
		t.Errorf("expected code %d, got %d", xerr.UserOrPasswordError, codeErr.Code)
	}
}

func TestLogin_AdminUser(t *testing.T) {
	svcCtx := testSvcCtx(&model.User{
		Id:       1,
		Username: "admin",
		Password: hashPassword("adminpass"),
		Role:     1, // 管理员
		Status:   1,
	}, nil)

	logic := NewLoginLogic(context.Background(), svcCtx)
	resp, err := logic.Login(&types.LoginReq{Username: "admin", Password: "adminpass"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.IsAdmin != 1 {
		t.Errorf("expected isAdmin 1, got %d", resp.IsAdmin)
	}
}
