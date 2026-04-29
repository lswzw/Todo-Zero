package model

import (
	"database/sql"
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

// testDB 创建一个内存 SQLite 数据库并初始化表结构，用于单元测试。
func testDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}

	initStmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			username varchar(50) NOT NULL UNIQUE,
			password varchar(255) NOT NULL,
			nickname varchar(50) DEFAULT '',
			email varchar(100) DEFAULT '',
			phone varchar(20) DEFAULT '',
			avatar varchar(255) DEFAULT '',
			role tinyint NOT NULL DEFAULT 0,
			status tinyint NOT NULL DEFAULT 1,
			is_deleted tinyint NOT NULL DEFAULT 0,
			failed_attempts integer NOT NULL DEFAULT 0,
			locked_until datetime DEFAULT NULL,
			create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			name varchar(50) NOT NULL,
			color varchar(20) NOT NULL DEFAULT '#1890ff',
			icon varchar(50) DEFAULT '',
			sort integer NOT NULL DEFAULT 0,
			user_id integer DEFAULT NULL,
			is_system tinyint NOT NULL DEFAULT 0,
			create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			title varchar(200) NOT NULL,
			content text,
			priority tinyint NOT NULL DEFAULT 0,
			status tinyint NOT NULL DEFAULT 0,
			category_id integer DEFAULT NULL,
			user_id integer NOT NULL,
			start_time datetime DEFAULT NULL,
			end_time datetime DEFAULT NULL,
			reminder datetime DEFAULT NULL,
			tags varchar(500) DEFAULT '',
sort_order integer NOT NULL DEFAULT 0,
			is_deleted tinyint NOT NULL DEFAULT 0,
			create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS system_configs (
			id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			config_key varchar(100) NOT NULL UNIQUE,
			config_value text NOT NULL,
			group_name varchar(50) NOT NULL DEFAULT 'default',
			description varchar(255) DEFAULT '',
			create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS operation_logs (
			id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			user_id integer DEFAULT NULL,
			username varchar(50) DEFAULT '',
			module varchar(50) NOT NULL,
			action varchar(100) NOT NULL,
			method varchar(20) NOT NULL,
			ip varchar(50) DEFAULT '',
			location varchar(255) DEFAULT '',
			params text,
			status tinyint NOT NULL DEFAULT 1,
			error_msg text,
			duration integer DEFAULT 0,
			created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS login_log (
			id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
			user_id integer DEFAULT NULL,
			username varchar(50) NOT NULL,
			ip varchar(50) DEFAULT '',
			user_agent varchar(500) DEFAULT '',
			status tinyint NOT NULL,
			remark varchar(255) DEFAULT '',
			create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for i, stmt := range initStmts {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("exec init stmt %d: %v", i, err)
		}
	}
	return db
}

// insertTestUser 插入一条测试用户并返回其 ID。
func insertTestUser(db *sql.DB, username string) int64 {
	result, err := db.Exec(
		"INSERT INTO users (username, password, role, status, is_deleted) VALUES (?, ?, 0, 1, 0)",
		username, fmt.Sprintf("hash_%s", username),
	)
	if err != nil {
		panic(fmt.Sprintf("insert test user: %v", err))
	}
	id, _ := result.LastInsertId()
	return id
}

// insertTestCategory 插入一条测试分类并返回其 ID。
func insertTestCategory(db *sql.DB, name string, isSystem int64) int64 {
	result, err := db.Exec(
		"INSERT INTO categories (name, color, is_system, sort) VALUES (?, '#409eff', ?, 0)",
		name, isSystem,
	)
	if err != nil {
		panic(fmt.Sprintf("insert test category: %v", err))
	}
	id, _ := result.LastInsertId()
	return id
}
