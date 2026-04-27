package model

import (
	"context"
	"database/sql"
	"testing"
)

func TestTaskModel_Insert(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	m := NewTaskModel(db)
	ctx := context.Background()

	task := &Task{
		Title:    "Test Task",
		Content:  sql.NullString{String: "hello", Valid: true},
		Priority: 1,
		Status:   0,
		UserId:   userId,
		Tags:     "go,test",
	}

	result, err := m.Insert(ctx, task)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	id, _ := result.LastInsertId()
	if id <= 0 {
		t.Errorf("expected positive id, got %d", id)
	}
	if task.CreateTime.IsZero() {
		t.Error("expected CreateTime to be set")
	}
}

func TestTaskModel_FindOne(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	m := NewTaskModel(db)
	ctx := context.Background()

	// 查找不存在的任务
	_, err := m.FindOne(ctx, 999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	// 插入后查找
	result, _ := m.Insert(ctx, &Task{Title: "FindMe", UserId: userId})
	id, _ := result.LastInsertId()

	found, err := m.FindOne(ctx, id)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if found.Title != "FindMe" {
		t.Errorf("expected title FindMe, got %s", found.Title)
	}
}

func TestTaskModel_SoftDelete(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	m := NewTaskModel(db)
	ctx := context.Background()

	result, _ := m.Insert(ctx, &Task{Title: "ToDelete", UserId: userId})
	id, _ := result.LastInsertId()

	// 软删除
	if err := m.Delete(ctx, id); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// 软删除后 FindOne 应返回 ErrNotFound
	_, err := m.FindOne(ctx, id)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after soft delete, got %v", err)
	}
}

func TestTaskModel_Update(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	m := NewTaskModel(db)
	ctx := context.Background()

	result, _ := m.Insert(ctx, &Task{Title: "Old", UserId: userId, Priority: 1})
	id, _ := result.LastInsertId()

	task, _ := m.FindOne(ctx, id)
	task.Title = "New"
	task.Priority = 2

	if err := m.Update(ctx, task); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := m.FindOne(ctx, id)
	if updated.Title != "New" {
		t.Errorf("expected title New, got %s", updated.Title)
	}
	if updated.Priority != 2 {
		t.Errorf("expected priority 2, got %d", updated.Priority)
	}
}

func TestTaskModel_UpdateStatus(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	m := NewTaskModel(db)
	ctx := context.Background()

	result, _ := m.Insert(ctx, &Task{Title: "Toggle", UserId: userId, Status: 0})
	id, _ := result.LastInsertId()

	if err := m.UpdateStatus(ctx, id, 2); err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	task, _ := m.FindOne(ctx, id)
	if task.Status != 2 {
		t.Errorf("expected status 2, got %d", task.Status)
	}
}

func TestTaskModel_FindList(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	otherUserId := insertTestUser(db, "otheruser")
	m := NewTaskModel(db)
	ctx := context.Background()

	// 插入多条任务
	m.Insert(ctx, &Task{Title: "Task A", UserId: userId, Status: 0, Priority: 1})
	m.Insert(ctx, &Task{Title: "Task B", UserId: userId, Status: 2, Priority: 3})
	m.Insert(ctx, &Task{Title: "Go Task", UserId: userId, Status: 0, Priority: 2, Tags: "go"})
	m.Insert(ctx, &Task{Title: "Other User", UserId: otherUserId, Status: 0})

	// 基本查询
	list, total, err := m.FindList(ctx, userId, "", -1, -1, 0, 1, 10)
	if err != nil {
		t.Fatalf("FindList failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 items, got %d", len(list))
	}

	// 按 status 过滤
	list, total, _ = m.FindList(ctx, userId, "", 2, -1, 0, 1, 10)
	if total != 1 {
		t.Errorf("expected 1 done task, got %d", total)
	}

	// 按关键字搜索
	list, total, _ = m.FindList(ctx, userId, "Go", -1, -1, 0, 1, 10)
	if total != 1 {
		t.Errorf("expected 1 keyword match, got %d", total)
	}

	// 按 priority 过滤
	list, total, _ = m.FindList(ctx, userId, "", -1, 1, 0, 1, 10)
	if total != 1 {
		t.Errorf("expected 1 priority=1 task, got %d", total)
	}

	// 分页
	list, total, _ = m.FindList(ctx, userId, "", -1, -1, 0, 2, 2)
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 item on page 2, got %d", len(list))
	}
}

func TestTaskModel_CountStats(t *testing.T) {
	db := testDB(t)
	defer db.Close()

	userId := insertTestUser(db, "testuser")
	m := NewTaskModel(db)
	ctx := context.Background()

	m.Insert(ctx, &Task{Title: "Todo1", UserId: userId, Status: 0})
	m.Insert(ctx, &Task{Title: "Todo2", UserId: userId, Status: 0})
	m.Insert(ctx, &Task{Title: "Done1", UserId: userId, Status: 2})

	total, todo, done, overdue, err := m.CountStats(ctx, userId)
	if err != nil {
		t.Fatalf("CountStats failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if todo != 2 {
		t.Errorf("expected todo 2, got %d", todo)
	}
	if done != 1 {
		t.Errorf("expected done 1, got %d", done)
	}
	if overdue != 0 {
		t.Errorf("expected overdue 0, got %d", overdue)
	}
}
