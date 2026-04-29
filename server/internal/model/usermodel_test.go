package model

import (
	"context"
	"testing"
)

func TestUserModel_Insert(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewUserModel(db)
	ctx := context.Background()

	user := &User{
		Username: "alice",
		Password: "hashed_password",
		Nickname: "Alice",
		Role:     0,
		Status:   1,
	}

	result, err := m.Insert(ctx, user)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	id, _ := result.LastInsertId()
	if id <= 0 {
		t.Errorf("expected positive id, got %d", id)
	}
	if user.CreateTime.IsZero() {
		t.Error("expected CreateTime to be set")
	}
}

func TestUserModel_FindOne(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewUserModel(db)
	ctx := context.Background()

	// 查找不存在的用户
	_, err := m.FindOne(ctx, 999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	// 插入后查找
	m.Insert(ctx, &User{Username: "bob", Password: "hash", Status: 1, Role: 0})
	found, _ := m.FindOneByUsername(ctx, "bob")

	user, err := m.FindOne(ctx, found.Id)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if user.Username != "bob" {
		t.Errorf("expected username bob, got %s", user.Username)
	}
}

func TestUserModel_FindOneByUsername(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewUserModel(db)
	ctx := context.Background()

	// 不存在
	_, err := m.FindOneByUsername(ctx, "nobody")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	// 存在
	m.Insert(ctx, &User{Username: "charlie", Password: "hash", Status: 1})
	user, err := m.FindOneByUsername(ctx, "charlie")
	if err != nil {
		t.Fatalf("FindOneByUsername failed: %v", err)
	}
	if user.Username != "charlie" {
		t.Errorf("expected username charlie, got %s", user.Username)
	}
}

func TestUserModel_SoftDelete(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewUserModel(db)
	ctx := context.Background()

	result, _ := m.Insert(ctx, &User{Username: "deleteme", Password: "hash", Status: 1})
	id, _ := result.LastInsertId()

	if err := m.Delete(ctx, id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// 软删除后 FindOne 应返回 ErrNotFound
	_, err := m.FindOne(ctx, id)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after soft delete, got %v", err)
	}

	// FindOneByUsername 也应返回 ErrNotFound
	_, err = m.FindOneByUsername(ctx, "deleteme")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after soft delete, got %v", err)
	}
}

func TestUserModel_Update(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewUserModel(db)
	ctx := context.Background()

	result, _ := m.Insert(ctx, &User{Username: "updateme", Password: "hash", Nickname: "Old", Status: 1})
	id, _ := result.LastInsertId()

	user, _ := m.FindOne(ctx, id)
	user.Nickname = "New"
	user.Email = "new@test.com"

	if err := m.Update(ctx, user); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := m.FindOne(ctx, id)
	if updated.Nickname != "New" {
		t.Errorf("expected nickname New, got %s", updated.Nickname)
	}
	if updated.Email != "new@test.com" {
		t.Errorf("expected email new@test.com, got %s", updated.Email)
	}
}

func TestUserModel_UpdateStatus(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewUserModel(db)
	ctx := context.Background()

	result, _ := m.Insert(ctx, &User{Username: "toggle", Password: "hash", Status: 1})
	id, _ := result.LastInsertId()

	if err := m.UpdateStatus(ctx, id, 0); err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	user, _ := m.FindOne(ctx, id)
	if user.Status != 0 {
		t.Errorf("expected status 0, got %d", user.Status)
	}
}

func TestUserModel_UpdatePassword(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewUserModel(db)
	ctx := context.Background()

	result, _ := m.Insert(ctx, &User{Username: "pwdtest", Password: "oldhash", Status: 1})
	id, _ := result.LastInsertId()

	if err := m.UpdatePassword(ctx, id, "newhash"); err != nil {
		t.Fatalf("UpdatePassword failed: %v", err)
	}

	user, _ := m.FindOne(ctx, id)
	// UpdatePassword 现在会用 bcrypt 加密，所以密码不等于明文
	if user.Password == "newhash" {
		t.Error("expected bcrypt-hashed password, got plaintext")
	}
}

func TestUserModel_FindList(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewUserModel(db)
	ctx := context.Background()

	m.Insert(ctx, &User{Username: "user_a", Password: "hash", Status: 1, Role: 0})
	m.Insert(ctx, &User{Username: "user_b", Password: "hash", Status: 0, Role: 0})
	m.Insert(ctx, &User{Username: "admin_x", Password: "hash", Status: 1, Role: 1})

	// 全部查询
	list, total, err := m.FindList(ctx, "", -1, 1, 10)
	if err != nil {
		t.Fatalf("FindList failed: %v", err)
	}
	// 包含 init.sql 插入的 admin 默认账户
	if total < 3 {
		t.Errorf("expected at least 3 users, got %d", total)
	}
	if len(list) < 3 {
		t.Errorf("expected at least 3 items, got %d", len(list))
	}

	// 按用户名搜索
	list, total, _ = m.FindList(ctx, "admin", -1, 1, 10)
	if total < 1 {
		t.Errorf("expected at least 1 admin match, got %d", total)
	}

	// 按状态过滤
	list, total, _ = m.FindList(ctx, "", 1, 1, 10)
	if total < 2 {
		t.Errorf("expected at least 2 active users, got %d", total)
	}
}

func TestUserModel_UsernameUnique(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewUserModel(db)
	ctx := context.Background()

	m.Insert(ctx, &User{Username: "unique", Password: "hash", Status: 1})

	// 重复用户名应报错
	_, err := m.Insert(ctx, &User{Username: "unique", Password: "hash2", Status: 1})
	if err == nil {
		t.Error("expected UNIQUE constraint error, got nil")
	}
}
