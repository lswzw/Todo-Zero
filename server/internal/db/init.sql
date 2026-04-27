-- ================================================
-- SQLite initialization script for Todo App
-- Run automatically on first startup
-- ================================================

-- Users table
CREATE TABLE IF NOT EXISTS `users` (
    `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `username` varchar(50) NOT NULL UNIQUE,
    `password` varchar(255) NOT NULL,
    `nickname` varchar(50) DEFAULT '',
    `email` varchar(100) DEFAULT '',
    `phone` varchar(20) DEFAULT '',
    `avatar` varchar(255) DEFAULT '',
    `role` tinyint NOT NULL DEFAULT 0,
    `status` tinyint NOT NULL DEFAULT 1,
    `is_deleted` tinyint NOT NULL DEFAULT 0,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Categories table
CREATE TABLE IF NOT EXISTS `categories` (
    `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `name` varchar(50) NOT NULL,
    `color` varchar(20) NOT NULL DEFAULT '#1890ff',
    `icon` varchar(50) DEFAULT '',
    `sort` integer NOT NULL DEFAULT 0,
    `user_id` integer DEFAULT NULL,
    `is_system` tinyint NOT NULL DEFAULT 0,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE SET NULL
);

-- Tasks table
CREATE TABLE IF NOT EXISTS `tasks` (
    `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `title` varchar(200) NOT NULL,
    `content` text,
    `priority` tinyint NOT NULL DEFAULT 0,
    `status` tinyint NOT NULL DEFAULT 0,
    `category_id` integer DEFAULT NULL,
    `user_id` integer NOT NULL,
    `start_time` datetime DEFAULT NULL,
    `end_time` datetime DEFAULT NULL,
    `reminder` datetime DEFAULT NULL,
    `tags` varchar(500) DEFAULT '',
    `is_deleted` tinyint NOT NULL DEFAULT 0,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

-- System configs table
CREATE TABLE IF NOT EXISTS `system_configs` (
    `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `config_key` varchar(100) NOT NULL UNIQUE,
    `config_value` text NOT NULL,
    `group_name` varchar(50) NOT NULL DEFAULT 'default',
    `description` varchar(255) DEFAULT '',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Operation logs table
CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `user_id` integer DEFAULT NULL,
    `username` varchar(50) DEFAULT '',
    `module` varchar(50) NOT NULL,
    `action` varchar(100) NOT NULL,
    `method` varchar(20) NOT NULL,
    `ip` varchar(50) DEFAULT '',
    `location` varchar(255) DEFAULT '',
    `params` text,
    `status` tinyint NOT NULL DEFAULT 1,
    `error_msg` text,
    `duration` integer DEFAULT 0,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE SET NULL
);

-- Login logs table
CREATE TABLE IF NOT EXISTS `login_log` (
    `id` integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    `user_id` integer DEFAULT NULL,
    `username` varchar(50) NOT NULL,
    `ip` varchar(50) DEFAULT '',
    `user_agent` varchar(500) DEFAULT '',
    `status` tinyint NOT NULL,
    `remark` varchar(255) DEFAULT '',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ================================================
-- Default admin account
-- Default admin account (password is hashed, please change after first login)
-- ================================================
INSERT OR IGNORE INTO `users` (`username`, `password`, `nickname`, `role`, `status`)
VALUES ('admin', '$2a$10$ktZcnwXDUmrcYXI.gOkjOuDeeivJNiuVLgW696nNDV7zYHzaNqgaW', '管理员', 1, 1);

-- ================================================
-- System preset categories
-- ================================================
INSERT OR IGNORE INTO `categories` (`name`, `color`, `icon`, `sort`, `user_id`, `is_system`) VALUES
('工作', '#f5222d', 'briefcase', 1, NULL, 1),
('生活', '#faad14', 'home', 2, NULL, 1),
('学习', '#52c41a', 'book', 3, NULL, 1),
('健康', '#13c2c2', 'heart', 4, NULL, 1);

-- ================================================
-- Default system configs
-- ================================================
INSERT OR IGNORE INTO `system_configs` (`config_key`, `config_value`, `group_name`, `description`) VALUES
('site_name', 'Todo 管理平台', 'basic', '站点名称'),
('allow_register', 'false', 'basic', '是否允许新用户注册'),
('task_default_priority', '0', 'task', '新建任务默认优先级'),
('task_auto_delete_days', '0', 'task', '自动清理已完成任务天数（0=不清理）'),
('task_trash_retention_days', '30', 'task', '回收站保留天数，超过后永久删除（0=不清理）'),
('log_auto_delete_days', '0', 'log', '自动清理操作日志和登录日志天数（0=不清理）'),
('db_backup_enabled', '0', 'backup', '是否启用数据库自动备份（0=关闭 1=开启）'),
('db_backup_interval_hours', '24', 'backup', '自动备份间隔小时数'),
('db_backup_max_count', '7', 'backup', '最大备份数量，超过后自动清理最旧的');

-- ================================================
-- Indexes for performance
-- ================================================
CREATE INDEX IF NOT EXISTS `idx_tasks_user_id` ON `tasks` (`user_id`, `is_deleted`);
CREATE INDEX IF NOT EXISTS `idx_tasks_status` ON `tasks` (`user_id`, `status`, `is_deleted`);
CREATE INDEX IF NOT EXISTS `idx_tasks_category_id` ON `tasks` (`category_id`);
CREATE INDEX IF NOT EXISTS `idx_users_username` ON `users` (`username`, `is_deleted`);
CREATE INDEX IF NOT EXISTS `idx_login_log_username` ON `login_log` (`username`);
CREATE INDEX IF NOT EXISTS `idx_tasks_completed_cleanup` ON `tasks` (`status`, `is_deleted`, `update_time`);
CREATE INDEX IF NOT EXISTS `idx_tasks_soft_deleted` ON `tasks` (`is_deleted`, `update_time`);
CREATE INDEX IF NOT EXISTS `idx_operation_logs_created_at` ON `operation_logs` (`created_at`);
CREATE INDEX IF NOT EXISTS `idx_login_log_create_time` ON `login_log` (`create_time`);
