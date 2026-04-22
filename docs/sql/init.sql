-- =============================================
-- Todo App 数据库设计
-- 数据库: zero
-- MySQL 版本: 5.7
-- =============================================

CREATE DATABASE IF NOT EXISTS zero DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE zero;

-- -------------------------------------------
-- 用户表
-- -------------------------------------------
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `password` varchar(200) NOT NULL COMMENT '密码(bcrypt)',
  `is_admin` tinyint(4) NOT NULL DEFAULT 0 COMMENT '是否管理员: 0=否 1=是',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态: 0=禁用 1=启用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- -------------------------------------------
-- 分类表
-- -------------------------------------------
CREATE TABLE `category` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(20) NOT NULL COMMENT '分类名称',
  `user_id` bigint(20) DEFAULT NULL COMMENT '用户ID(NULL表示系统预置)',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分类表';

-- -------------------------------------------
-- 任务表
-- -------------------------------------------
CREATE TABLE `task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `title` varchar(100) NOT NULL COMMENT '任务标题',
  `content` varchar(1000) DEFAULT NULL COMMENT '任务内容',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态: 0=待办 1=已完成',
  `priority` tinyint(4) NOT NULL DEFAULT 2 COMMENT '优先级: 1=高 2=中 3=低',
  `category_id` bigint(20) DEFAULT NULL COMMENT '分类ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_user_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务表';

-- -------------------------------------------
-- 系统配置表
-- -------------------------------------------
CREATE TABLE `system_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `config_key` varchar(50) NOT NULL COMMENT '配置键',
  `config_value` varchar(500) NOT NULL COMMENT '配置值',
  `remark` varchar(200) DEFAULT NULL COMMENT '备注说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- -------------------------------------------
-- 操作日志表
-- -------------------------------------------
CREATE TABLE `operation_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint(20) NOT NULL COMMENT '操作用户ID',
  `username` varchar(20) NOT NULL COMMENT '操作用户名',
  `action` varchar(50) NOT NULL COMMENT '操作类型',
  `target_type` varchar(50) DEFAULT NULL COMMENT '操作对象类型(user/task/category/config)',
  `target_id` bigint(20) DEFAULT NULL COMMENT '操作对象ID',
  `detail` varchar(500) DEFAULT NULL COMMENT '操作详情',
  `ip` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- -------------------------------------------
-- 登录日志表
-- -------------------------------------------
CREATE TABLE `login_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint(20) DEFAULT NULL COMMENT '用户ID',
  `username` varchar(20) NOT NULL COMMENT '用户名',
  `ip` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` varchar(500) DEFAULT NULL COMMENT '浏览器UA',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态: 0=失败 1=成功',
  `remark` varchar(200) DEFAULT NULL COMMENT '备注(失败原因等)',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';

-- -------------------------------------------
-- 初始数据
-- -------------------------------------------

-- 系统预置分类
INSERT INTO `category` (`name`, `sort_order`) VALUES
('工作', 1),
('生活', 2),
('学习', 3);

-- 默认管理员 (密码: admin123)
INSERT INTO `user` (`username`, `password`, `is_admin`, `status`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJHd1F6d0m2', 1, 1);

-- 系统默认配置
INSERT INTO `system_config` (`config_key`, `config_value`, `remark`) VALUES
('allow_register', 'true', '是否允许公开注册: true/false'),
('site_name', 'Todo App', '站点名称');
