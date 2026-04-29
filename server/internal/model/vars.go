package model

import (
	"database/sql"
	"errors"
	"time"
)

// Common errors
var ErrNotFound = errors.New("record not found")

// User represents a user record in the database.
type User struct {
	Id             int64          `db:"id"`              // 用户ID
	Username       string         `db:"username"`        // 用户名
	Password       string         `db:"password"`        // 密码(存储时由 Model 层使用 bcrypt 加密)
	Nickname       string         `db:"nickname"`        // 昵称
	Email          string         `db:"email"`           // 邮箱
	Phone          string         `db:"phone"`           // 手机号
	Avatar         string         `db:"avatar"`          // 头像URL
	Role           int64          `db:"role"`            // 角色: 0=普通用户 1=管理员
	Status         int64          `db:"status"`          // 状态: 0=禁用 1=正常
	IsDeleted      int64          `db:"is_deleted"`      // 是否删除: 0=否 1=是
	FailedAttempts int64          `db:"failed_attempts"` // 登录失败次数
	LockedUntil    sql.NullTime   `db:"locked_until"`    // 账户锁定时间
	CreateTime     time.Time      `db:"create_time"`     // 创建时间
	UpdateTime     time.Time      `db:"update_time"`     // 更新时间
}

// Task represents a task record in the database.
type Task struct {
	Id         int64          `db:"id"`          // 任务ID
	Title      string         `db:"title"`       // 任务标题
	Content    sql.NullString `db:"content"`     // 任务内容
	Priority   int64          `db:"priority"`    // 优先级: 1=重要 2=紧急 3=普通
	Status     int64          `db:"status"`      // 状态: 0=待办 2=已完成
	CategoryId sql.NullInt64  `db:"category_id"` // 分类ID
	UserId     int64          `db:"user_id"`     // 用户ID
	StartTime  sql.NullTime   `db:"start_time"`  // 开始时间
	EndTime    sql.NullTime   `db:"end_time"`    // 截止时间
	Reminder   sql.NullTime   `db:"reminder"`    // 提醒时间
	Tags       string         `db:"tags"`        // 标签(逗号分隔)
	SortOrder  int64          `db:"sort_order"`  // 排序顺序
	IsDeleted  int64          `db:"is_deleted"`  // 是否删除: 0=否 1=是
	CreateTime time.Time      `db:"create_time"` // 创建时间
	UpdateTime time.Time      `db:"update_time"` // 更新时间
}

// Category represents a category record in the database.
type Category struct {
	Id         int64          `db:"id"`          // 分类ID
	Name       string         `db:"name"`        // 分类名称
	Color      string         `db:"color"`       // 分类颜色
	Icon       string         `db:"icon"`        // 分类图标
	Sort       int64          `db:"sort"`        // 排序
	UserId     sql.NullInt64  `db:"user_id"`     // 用户ID(为NULL表示系统预置)
	IsSystem   int64          `db:"is_system"`   // 是否系统预置: 0=否 1=是
	CreateTime time.Time      `db:"create_time"` // 创建时间
	UpdateTime time.Time      `db:"update_time"` // 更新时间
}

// SystemConfig represents a system configuration record.
type SystemConfig struct {
	Id          int64     `db:"id"`          // 配置ID
	ConfigKey   string    `db:"config_key"` // 配置键
	ConfigValue string    `db:"config_value"` // 配置值
	GroupName   string    `db:"group_name"` // 配置分组
	Description string    `db:"description"` // 配置描述
	CreateTime  time.Time `db:"create_time"` // 创建时间
	UpdateTime  time.Time `db:"update_time"` // 更新时间
}

// OperationLog represents an operation log record.
type OperationLog struct {
	Id        int64          `db:"id"`         // 日志ID
	UserId    sql.NullInt64  `db:"user_id"`    // 用户ID
	Username  string         `db:"username"`  // 用户名
	Module    string         `db:"module"`    // 操作模块
	Action    string         `db:"action"`    // 操作动作
	Method    string         `db:"method"`    // 请求方法
	Ip        string         `db:"ip"`        // IP地址
	Location  string         `db:"location"`  // 地理位置
	Params    string         `db:"params"`    // 请求参数
	Status    int64          `db:"status"`    // 状态: 1=成功 0=失败
	ErrorMsg  string         `db:"error_msg"` // 错误信息
	Duration  int64          `db:"duration"`  // 耗时(毫秒)
	CreatedAt time.Time      `db:"created_at"` // 创建时间
}

// LoginLog represents a login log record.
type LoginLog struct {
	Id        int64          `db:"id"`         // 日志ID
	UserId    sql.NullInt64  `db:"user_id"`    // 用户ID
	Username  string         `db:"username"`  // 用户名
	Ip        string         `db:"ip"`        // IP地址
	UserAgent string         `db:"user_agent"` // 浏览器UA
	Status    int64          `db:"status"`    // 状态: 1=成功 0=失败
	Remark    string         `db:"remark"`    // 备注(失败原因等)
	CreateTime time.Time     `db:"create_time"` // 创建时间
}
