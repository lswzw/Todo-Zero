# 📝 Todo App 开发路线图

## 阶段零：环境准备

- [x] 安装 Go (1.21+) — 已安装 go1.24.10 linux/amd64
- [x] 安装 goctl（go-zero 脚手架工具）— 已安装 goctl 1.10.1 linux/amd64
- [x] 安装 MySQL + 创建数据库 — 已搭建 192.168.100.192 DB:zero
- [x] 安装 Redis — 已搭建 192.168.102.77:6379 v6.2.21
- [x] 安装 Node.js (18+) — 已安装 v20.18.0
- [x] 配置 IDE（Go 插件 + Vue 插件）— 已安装 Go 插件

---

## 阶段一：需求分析 & 设计

| 步骤 | 产出 | 状态 |
|------|------|------|
| 功能清单 | 明确所有功能点 | [x] → `docs/requirements.md` |
| 原型设计 | 页面线框图、交互说明 | [x] → `docs/prototype.md` |
| 数据库设计 | 表结构、字段定义 | [x] → `docs/sql/init.sql` |
| API 设计 | 接口路径、请求/响应格式 | [x] → `docs/api/todo.api` |
| 前端页面规划 | 页面列表、交互流程 | [x] → `docs/prototype.md` |

### Todo App 功能清单

#### 用户模块
- [x] 用户注册/登录（JWT）
- [x] 修改密码（用户自助）
- [x] 注册开关（管理员控制）
- [x] 用户禁用/启用

#### 任务模块
- [x] 任务创建/编辑/删除
- [x] 任务状态切换（待办/已完成）
- [x] 任务内容（content 字段）
- [x] 任务分类（工作/生活/学习...）
- [x] 任务搜索/筛选
- [x] 批量操作（完成/取消/删除）

#### 管理员模块
- [x] 用户管理（列表/重置密码/禁用/删除）
- [x] 系统设置（注册开关等配置）
- [x] 操作日志
- [x] 登录日志

#### 统计模块
- [x] 统计概览（完成率等）

#### 原型页面（7个）
- [x] 登录页 → `docs/prototype/login.html`
- [x] 注册页 → `docs/prototype/register.html`
- [x] 任务主页 → `docs/prototype/index.html`
- [x] 用户管理 → `docs/prototype/admin-user.html`
- [x] 系统设置 → `docs/prototype/admin-config.html`
- [x] 操作日志 → `docs/prototype/admin-log.html`
- [x] 登录日志 → `docs/prototype/admin-login-log.html`

---

## 阶段二：后端开发

- [x] ① 编写 .api 文件（定义接口）→ `docs/api/todo.api`
- [x] ② goctl 生成项目骨架代码 → `server/` 目录
- [x] ③ 建表 + 编写 SQL → `docs/sql/init.sql`（需在数据库服务器执行）
- [x] ④ 编写业务逻辑（handler → logic → model）→ 24个接口全部完成
- [x] ⑤ 配置 etc 文件（MySQL/Redis/JWT 等）→ `server/etc/todo-api.yaml`
- [x] ⑥ 接口测试（用 curl / Postman 验证）→ 全部 23 个接口测试通过

---

## 阶段三：前端开发

- [x] ① 创建 Vue3 项目（Vite 脚手架）→ `web/` 目录
- [x] ② 安装 Element Plus + Axios + Vue Router + Pinia
- [x] ③ 页面布局（登录页 / 注册页 / 任务主页 / 管理后台4页）
- [x] ④ 封装 API 请求（对接后端 23 个接口）→ `web/src/api/`
- [x] ⑤ 状态管理（Pinia 管理用户登录态）→ `web/src/stores/`
- [x] ⑥ 联调测试 → 全部 23 个接口通过前端代理正常访问

---

## 阶段四：单体二进制改造

> **目标**：最终产物为一个独立二进制文件，无需外部依赖（MySQL/Redis/Nginx），`./todo-app` 即可运行。

### 设计思路

- **SQLite 替代 MySQL**：使用 `modernc.org/sqlite`（纯 Go 实现，无需 CGO），数据存储为本地文件 `data/todo.db`
- **去掉 Redis**：SQLite 本地文件访问速度快，缓存意义不大；所有查询直接走数据库
- **Go embed 嵌入前端**：`//go:embed dist/` 将 Vue 构建产物嵌入二进制，通过 `http.FS` 提供静态文件服务
- **自动初始化**：启动时检测 SQLite 文件是否存在，不存在则自动建库建表、创建默认管理员

### 改造清单

- [x] ① MySQL → SQLite：替换驱动、调整 SQL 语法（`AUTO_INCREMENT` → `AUTOINCREMENT`、去掉反引号等）
- [x] ② 去掉 Redis：移除缓存配置和逻辑，Model 层改为纯数据库查询
- [x] ③ 配置精简：`todo-api.yaml` 去掉 Redis/MySQL 配置，仅保留 JWT 等必要项
- [x] ④ 启动自动初始化：检测并创建 SQLite 数据库 + 建表 + 默认数据
- [x] ⑤ 嵌入前端静态文件：`embed dist/` + go-zero NotFoundHandler 兜底返回 `index.html`
- [x] ⑥ 编译测试：`go build -o todo-api` 生成单体二进制，全部接口功能验证通过

### 最终效果

```bash
# 编译
go build -o todo-app

# 运行（一个文件搞定）
./todo-app
# → 自动创建 data/todo.db，自动建表
# → 浏览器访问 http://localhost:8888 直接使用
```

---

## 阶段五：完整测试 & 发布

- [x] ① 完整功能测试（所有 CRUD + 管理员功能）→ 23 个接口全部通过，修复密码修改 bug
- [x] ② 多平台编译（Linux / macOS / Windows）→ 4 个平台版本已生成
- [x] ③ 打标签发布 → `v1.0.0`

---

## 整体时间线（AI 辅助开发）

| 阶段 | 内容 | AI 协作方式 | 完成状态 |
|------|------|------------|---------|
| Day 1 | 环境搭建 + 项目初始化 | AI 生成安装命令和项目骨架 | [ ] |
| Day 2 | 需求设计 + API 设计 + 建表 | AI 编写 .api 文件和 SQL | [ ] |
| Day 3-4 | 后端业务逻辑开发 | AI 编写 CRUD 逻辑代码 | [ ] |
| Day 5-6 | 前端页面开发 | AI 生成页面和组件 | [ ] |
| Day 7 | 单体二进制改造 + 联调 | MySQL→SQLite，去Redis，embed前端 | [ ] |
| Day 8 | 完整测试 + 发布 | AI 辅助定位和修复 Bug | [ ] |

---

## 技术选型

| 维度 | 选择 | 理由 |
|------|------|------|
| 后端框架 | go-zero | 高性能微服务框架 |
| 前端框架 | Vue3 | 上手快，中文生态好 |
| UI 组件库 | Element Plus | 企业级组件库，文档齐全 |
| 数据库 | SQLite（`modernc.org/sqlite`） | 纯 Go 实现，无需 CGO，单文件部署 |
| ~~缓存~~ | ~~Redis~~ → 去掉 | SQLite 本地访问快，无需缓存层 |
| ORM | 原生 SQL / sqlx | go-zero 官方推荐 |
| API 风格 | RESTful | go-zero 原生支持 |
| 认证方式 | JWT | go-zero 内置支持 |
| 静态资源 | Go embed | 前端产物嵌入二进制，单体部署 |

### 架构演进说明

| 阶段 | 数据库 | 缓存 | 前端部署 | 部署方式 |
|------|--------|------|----------|----------|
| 开发期（阶段二） | MySQL | Redis | — | 前后端分离 |
| **最终目标（阶段四）** | **SQLite** | **无** | **Go embed** | **单体二进制** |
