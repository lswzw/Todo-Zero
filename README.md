# Todo-Zero

零依赖全栈待办引擎 — Go 后端 + Vue3 前端编译为**单一二进制**，一行命令即可部署。无需 MySQL、Redis、Nginx，开箱即用。

## ✨ 特性

- **单体部署** — 一个二进制文件包含后端 API + 前端页面，零外部依赖
- **SQLite 存储** — 使用纯 Go 实现的 SQLite 驱动，无需 CGO，数据存为本地文件
- **自动初始化** — 首次启动自动建库建表，创建默认管理员账号
- **JWT 认证** — 安全的无状态认证机制
- **SPA 前端** — Vue3 + Element Plus + TypeScript，嵌入二进制提供服务
- **定时任务** — 自动清理已删除的任务、定期备份数据库

## 🚀 快速开始

### 编译构建

```bash
# 1. 构建前端（输出到 server/dist/）
cd web
npm install
npm run build

# 2. 编译后端（前端已嵌入）
cd ../server
go build -o todo-api .
```

### 运行方式

```bash
# 直接运行（零配置）
./todo-api

# 自定义参数
./todo-api -port 9090 -data-dir /var/todo

# 使用配置文件
./todo-api -f /etc/todo-api.yaml

# 查看帮助
./todo-api -h
```

首次启动会自动创建 `data/todo.db`，并初始化表结构和默认数据。

### 默认管理员

| 用户名 | 密码 |
|--------|------|
| admin | admin123 |

访问 http://localhost:8888 即可使用。

## 📁 项目结构

```
.
├── server/                  # Go 后端
│   ├── todo.go              # 入口：SQLite初始化 + embed静态文件 + API注册
│   ├── etc/todo-api.yaml    # 配置文件
│   ├── internal/
│   │   ├── config/          # 配置定义
│   │   ├── db/              # SQLite初始化 + 建表SQL
│   │   ├── handler/         # HTTP处理器（goctl生成）
│   │   ├── logic/           # 业务逻辑
│   │   ├── model/           # 数据访问层（*sql.DB）
│   │   ├── middleware/      # 中间件（认证、限流、日志）
│   │   ├── pkg/xerr/        # 统一错误处理
│   │   ├── scheduler/       # 定时任务（备份、清理）
│   │   └── svc/             # 服务上下文（依赖注入）
│   └── dist/                # 前端构建产物（go:embed嵌入）
├── web/                     # Vue3 前端
│   ├── src/
│   │   ├── api/             # Axios请求封装 + API函数
│   │   ├── router/          # Vue Router路由配置
│   │   ├── stores/          # Pinia状态管理
│   │   ├── locales/         # 国际化（中英文）
│   │   └── views/           # 页面组件
│   └── vite.config.ts       # Vite配置（输出到server/dist）
├── docs/                    # 文档
│   ├── api/                 # API定义文件
│   ├── prototype/           # 页面原型
│   └── development-roadmap.md
└── data/                    # 运行时数据（SQLite数据库，.gitignore）
```

## 🎯 功能模块

### 用户管理

- 注册 / 登录（JWT）
- 修改密码
- 注册开关（管理员控制）

### 任务管理

- 创建 / 编辑 / 删除任务
- 状态切换（待办 / 已完成）
- 按分类、优先级、关键词筛选
- 批量操作（完成 / 取消 / 删除）
- 回收站（软删除、恢复、永久删除）

### 分类管理

- 创建 / 编辑 / 删除分类
- 任务按分类组织

### 管理员功能

- 用户管理（列表 / 重置密码 / 禁用 / 删除）
- 系统配置（注册开关、站点名称等）
- 操作日志查看
- 登录日志查看
- 数据库备份 / 恢复

### 统计概览

- 任务总数、完成数、待办数、完成率

## 📄 页面路由

| 路径 | 页面 | 认证 |
|------|------|------|
| `/login` | 登录 | - |
| `/register` | 注册 | - |
| `/` | 任务主页 | 需要 |
| `/task/:id` | 任务详情 | 需要 |
| `/trash` | 回收站 | 需要 |
| `/admin/user` | 用户管理 | 管理员 |
| `/admin/config` | 系统配置 | 管理员 |
| `/admin/log` | 操作日志 | 管理员 |
| `/admin/login-log` | 登录日志 | 管理员 |
| `/admin/backup` | 数据备份 | 管理员 |

## 🔌 API 接口

所有接口前缀 `/api/v1`，需要认证的接口需携带 `Authorization: Bearer <token>` 头。

### 用户接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /user/check-register | 检查是否允许注册 | - |
| POST | /user/login | 登录 | - |
| POST | /user/register | 注册 | - |
| GET | /user/info | 当前用户信息 | 需要 |
| PUT | /user/password | 修改密码 | 需要 |

### 任务接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /task | 任务列表 | 需要 |
| POST | /task | 创建任务 | 需要 |
| GET | /task/:id | 任务详情 | 需要 |
| PUT | /task/:id | 更新任务 | 需要 |
| DELETE | /task/:id | 删除任务 | 需要 |
| PATCH | /task/:id/toggle | 切换状态 | 需要 |
| POST | /task/batch | 批量操作 | 需要 |
| GET | /task/trash | 回收站列表 | 需要 |
| POST | /task/:id/restore | 恢复任务 | 需要 |
| DELETE | /task/:id/permanent | 永久删除 | 需要 |
| POST | /task/sort | 排序 | 需要 |

### 分类接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /category | 分类列表 | 需要 |
| POST | /category | 创建分类 | 需要 |
| PUT | /category/:id | 更新分类 | 需要 |
| DELETE | /category/:id | 删除分类 | 需要 |

### 统计接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /stat | 统计概览 | 需要 |

### 管理员接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /admin/user | 用户列表 | 管理员 |
| DELETE | /admin/user/:id | 删除用户 | 管理员 |
| PUT | /admin/user/:id/password | 重置密码 | 管理员 |
| PATCH | /admin/user/:id/toggle | 禁用/启用 | 管理员 |
| GET | /admin/config | 系统配置 | 管理员 |
| PUT | /admin/config | 更新配置 | 管理员 |
| GET | /admin/log/operation | 操作日志 | 管理员 |
| GET | /admin/log/login | 登录日志 | 管理员 |
| POST | /admin/backup | 触发备份 | 管理员 |
| GET | /admin/backup/list | 备份列表 | 管理员 |
| POST | /admin/backup/restore | 恢复备份 | 管理员 |
| GET | /admin/backup/download | 下载备份 | 管理员 |

## ⚙️ 配置说明

支持命令行参数和配置文件两种方式，命令行参数优先级更高。

### 命令行参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-host` | `0.0.0.0` | 监听地址 |
| `-port` | `8888` | 监听端口 |
| `-data-dir` | `data` | 数据存储目录 |
| `-db-file` | `todo.db` | SQLite 数据库文件名 |
| `-jwt-secret` | `todo-app-jwt-secret-key-2024` | JWT 签名密钥 |
| `-jwt-expire` | `86400` | Token 有效期（秒） |
| `-f` | - | 配置文件路径（指定后忽略命令行参数） |

### 配置文件

```yaml
Name: todo-api
Host: 0.0.0.0
Port: 8888

Auth:
  AccessSecret: "your-jwt-secret"    # JWT签名密钥，生产环境请修改
  AccessExpire: 86400                 # Token有效期（秒）

Database:
  DataDir: "data"                     # 数据目录
  DBFile: "todo.db"                   # SQLite文件名
```

## 🛠️ 技术栈

| 维度 | 选择 | 说明 |
|------|------|------|
| 后端框架 | go-zero | 高性能Go微服务框架 |
| 前端框架 | Vue3 + TypeScript | 现代前端技术栈 |
| UI 组件库 | Element Plus | 基于Vue3的UI库 |
| 数据库 | SQLite | 嵌入式数据库，零依赖 |
| 认证 | JWT | 无状态认证机制 |
| 静态资源 | Go embed | 编译时嵌入静态文件 |
| 构建工具 | Vite | 快速前端构建工具 |
| 状态管理 | Pinia | Vue官方状态管理 |
| 路由 | Vue Router | Vue官方路由 |

## 👨‍💻 开发指南

### 前端开发

```bash
cd web
npm install
npm run dev          # 启动开发服务器（代理API到localhost:8888）
```

### 后端开发

```bash
cd server
go run todo.go          # 默认配置启动
go run todo.go -port 9090  # 自定义端口
```

### 重新构建前端到二进制

```bash
cd web && npm run build && cd ../server && go build -o todo-api .
```

### 运行测试

```bash
cd server
go test ./...           # 运行所有测试
```

## 📝 License

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
