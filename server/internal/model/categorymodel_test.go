package model

import (
	"context"
	"database/sql"
	"testing"
)

func TestCategoryModel_Insert(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	m := NewCategoryModel(db)
	ctx := context.Background()

	cat := &Category{
		Name:     "MyCat",
		Color:    "#ff0000",
		UserId:   sql.NullInt64{Int64: userId, Valid: true},
		IsSystem: 0,
		Sort:     1,
	}

	result, err := m.Insert(ctx, cat)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	id, _ := result.LastInsertId()
	if id <= 0 {
		t.Errorf("expected positive id, got %d", id)
	}
}

func TestCategoryModel_FindOne(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewCategoryModel(db)
	ctx := context.Background()

	// 不存在
	_, err := m.FindOne(ctx, 999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	// 存在
	id := insertTestCategory(db, "Work", 1)
	cat, err := m.FindOne(ctx, id)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if cat.Name != "Work" {
		t.Errorf("expected name Work, got %s", cat.Name)
	}
	if cat.IsSystem != 1 {
		t.Errorf("expected is_system 1, got %d", cat.IsSystem)
	}
}

func TestCategoryModel_Delete_SystemCategory(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewCategoryModel(db)
	ctx := context.Background()

	// 系统分类不可删除（Delete 限定 is_system=0）
	id := insertTestCategory(db, "System", 1)
	if err := m.Delete(ctx, id); err != nil {
		t.Fatalf("Delete should not error, got %v", err)
	}

	// 系统分类仍在
	cat, err := m.FindOne(ctx, id)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if cat == nil {
		t.Error("system category should not be deleted")
	}
}

func TestCategoryModel_Delete_UserCategory(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewCategoryModel(db)
	ctx := context.Background()

	// 用户分类可以删除
	id := insertTestCategory(db, "MyCat", 0)
	if err := m.Delete(ctx, id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := m.FindOne(ctx, id)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestCategoryModel_FindAll(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	m := NewCategoryModel(db)
	ctx := context.Background()

	// 插入系统分类和用户分类
	insertTestCategory(db, "SysCat", 1)
	db.Exec("INSERT INTO categories (name, color, is_system, sort, user_id) VALUES (?, '#000', 0, 0, ?)", "UserCat", userId)

	// 查询该用户的分类（应包含系统 + 自己的）
	list, err := m.FindAll(ctx, userId)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	// 至少有 2 个（1 系统 + 1 用户）
	if len(list) < 2 {
		t.Errorf("expected at least 2 categories, got %d", len(list))
	}
}

func TestCategoryModel_FindSystem(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewCategoryModel(db)
	ctx := context.Background()

	insertTestCategory(db, "SysA", 1)
	insertTestCategory(db, "SysB", 1)
	insertTestCategory(db, "UserC", 0)

	list, err := m.FindSystem(ctx)
	if err != nil {
		t.Fatalf("FindSystem failed: %v", err)
	}
	if len(list) < 2 {
		t.Errorf("expected at least 2 system categories, got %d", len(list))
	}
	for _, c := range list {
		if c.IsSystem != 1 {
			t.Errorf("FindSystem returned non-system category: %s", c.Name)
		}
	}
}

func TestCategoryModel_Update(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	m := NewCategoryModel(db)
	ctx := context.Background()

	id := insertTestCategory(db, "Old", 0)
	cat, _ := m.FindOne(ctx, id)
	cat.Name = "New"
	cat.Color = "#00ff00"

	if err := m.Update(ctx, cat); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := m.FindOne(ctx, id)
	if updated.Name != "New" {
		t.Errorf("expected name New, got %s", updated.Name)
	}
	if updated.Color != "#00ff00" {
		t.Errorf("expected color #00ff00, got %s", updated.Color)
	}
}

func TestCategoryModel_CountByUser(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	m := NewCategoryModel(db)
	ctx := context.Background()

	// 用户还没自定义分类
	count, err := m.CountByUser(ctx, userId)
	if err != nil {
		t.Fatalf("CountByUser failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}

	// 插入用户分类
	db.Exec("INSERT INTO categories (name, color, is_system, sort, user_id) VALUES (?, '#000', 0, 0, ?)", "Cat1", userId)
	db.Exec("INSERT INTO categories (name, color, is_system, sort, user_id) VALUES (?, '#000', 0, 0, ?)", "Cat2", userId)

	count, _ = m.CountByUser(ctx, userId)
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}
