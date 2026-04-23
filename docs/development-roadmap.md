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

- [ ] ① 创建 Vue3 项目（Vite 脚手架）
- [ ] ② 安装 Element Plus + Axios + Vue Router + Pinia
- [ ] ③ 页面布局（登录页 / 主页 / 任务列表）
- [ ] ④ 封装 API 请求（对接后端接口）
- [ ] ⑤ 状态管理（Pinia 管理用户/任务数据）
- [ ] ⑥ 联调测试

---

## 阶段四：联调 & 部署

- [ ] ① 前后端联调，解决跨域等问题
- [ ] ② 后端编译部署（Docker / 直接编译）
- [ ] ③ 前端构建 + 部署（Nginx 托管静态资源）
- [ ] ④ 完整功能测试

---

## 整体时间线（AI 辅助开发）

| 阶段 | 内容 | AI 协作方式 | 完成状态 |
|------|------|------------|---------|
| Day 1 | 环境搭建 + 项目初始化 | AI 生成安装命令和项目骨架 | [ ] |
| Day 2 | 需求设计 + API 设计 + 建表 | AI 编写 .api 文件和 SQL | [ ] |
| Day 3-4 | 后端业务逻辑开发 | AI 编写 CRUD 逻辑代码 | [ ] |
| Day 5-6 | 前端页面开发 | AI 生成页面和组件 | [ ] |
| Day 7 | 联调 + 测试 + 修复 | AI 辅助定位和修复 Bug | [ ] |

---

## 技术选型

| 维度 | 选择 | 理由 |
|------|------|------|
| 后端框架 | go-zero | 高性能微服务框架 |
| 前端框架 | Vue3 | 上手快，中文生态好 |
| UI 组件库 | Element Plus | 企业级组件库，文档齐全 |
| 数据库 | MySQL | 最通用，学习资料最多 |
| 缓存 | Redis | go-zero 原生支持 |
| ORM | 原生 SQL / sqlx | go-zero 官方推荐 |
| API 风格 | RESTful | go-zero 原生支持 |
| 认证方式 | JWT | go-zero 内置支持 |
