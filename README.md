# Todo App

一个基于 Go + Vue3 的全栈待办事项管理应用，**编译为单个二进制文件**，无需 MySQL、Redis、Nginx 等外部依赖，`./todo-api` 即可运行。

## 特性

- **单体部署** — 一个二进制文件包含后端 API + 前端页面，零外部依赖
- **SQLite 存储** — 使用纯 Go 实现的 SQLite 驱动，无需 CGO，数据存为本地文件
- **自动初始化** — 首次启动自动建库建表，创建默认管理员账号
- **JWT 认证** — 安全的无状态认证机制
- **SPA 前端** — Vue3 + Element Plus + TypeScript，嵌入二进制提供服务

## 快速开始

### 编译

```bash
# 1. 构建前端
cd web
npm install
npm run build        # 输出到 ../server/dist/

# 2. 编译后端（前端已嵌入）
cd ../server
go build -o todo-api .

# 3. 运行
./todo-api -f etc/todo-api.yaml
```

首次启动会自动创建 `data/todo.db`，并初始化表结构和默认数据。

### 默认管理员

| 用户名 | 密码 |
|--------|------|
| admin | admin123 |

访问 http://localhost:8888 即可使用。

## 项目结构

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
│   │   ├── pkg/xerr/        # 统一错误处理
│   │   └── svc/             # 服务上下文（依赖注入）
│   └── dist/                # 前端构建产物（go:embed嵌入）
├── web/                     # Vue3 前端
│   ├── src/
│   │   ├── api/             # Axios请求封装 + 23个API函数
│   │   ├── router/          # Vue Router路由配置
│   │   ├── stores/          # Pinia状态管理
│   │   └── views/           # 页面组件
│   └── vite.config.ts       # Vite配置（输出到server/dist）
├── docs/                    # 文档
│   ├── development-roadmap.md
│   ├── frontend-guide.md
│   ├── api/                 # API定义文件
│   ├── sql/                 # MySQL原始建表SQL
│   └── prototype/           # 页面原型
└── data/                    # 运行时数据（SQLite数据库，.gitignore）
```

## 功能模块

### 用户

- 注册 / 登录（JWT）
- 修改密码
- 注册开关（管理员控制）

### 任务

- 创建 / 编辑 / 删除任务
- 状态切换（待办 / 已完成）
- 按分类、优先级、关键词筛选
- 批量操作（完成 / 取消 / 删除）

### 管理员

- 用户管理（列表 / 重置密码 / 禁用 / 删除）
- 系统配置（注册开关、站点名称等）
- 操作日志查看
- 登录日志查看

### 统计

- 任务总数、完成数、待办数、完成率

## 页面路由

| 路径 | 页面 | 认证 |
|------|------|------|
| `/login` | 登录 | - |
| `/register` | 注册 | - |
| `/` | 任务主页 | 需要 |
| `/admin/user` | 用户管理 | 管理员 |
| `/admin/config` | 系统配置 | 管理员 |
| `/admin/log` | 操作日志 | 管理员 |
| `/admin/login-log` | 登录日志 | 管理员 |

## API 接口

所有接口前缀 `/api/v1`，需要认证的接口需携带 `Authorization: Bearer <token>` 头。

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /user/check-register | 检查是否允许注册 | - |
| POST | /user/login | 登录 | - |
| POST | /user/register | 注册 | - |
| GET | /user/info | 当前用户信息 | 需要 |
| PUT | /user/password | 修改密码 | 需要 |
| GET | /task | 任务列表 | 需要 |
| POST | /task | 创建任务 | 需要 |
| GET | /task/:id | 任务详情 | 需要 |
| PUT | /task/:id | 更新任务 | 需要 |
| DELETE | /task/:id | 删除任务 | 需要 |
| PATCH | /task/:id/toggle | 切换状态 | 需要 |
| POST | /task/batch | 批量操作 | 需要 |
| GET | /category | 分类列表 | 需要 |
| POST | /category | 创建分类 | 需要 |
| GET | /stat | 统计概览 | 需要 |
| GET | /admin/user | 用户列表 | 管理员 |
| DELETE | /admin/user/:id | 删除用户 | 管理员 |
| PUT | /admin/user/:id/password | 重置密码 | 管理员 |
| PATCH | /admin/user/:id/toggle | 禁用/启用 | 管理员 |
| GET | /admin/config | 系统配置 | 管理员 |
| PUT | /admin/config | 更新配置 | 管理员 |
| GET | /admin/log/operation | 操作日志 | 管理员 |
| GET | /admin/log/login | 登录日志 | 管理员 |

## 配置说明

配置文件 `server/etc/todo-api.yaml`：

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

## 技术栈

| 维度 | 选择 |
|------|------|
| 后端框架 | go-zero |
| 前端框架 | Vue3 + TypeScript |
| UI 组件库 | Element Plus |
| 数据库 | SQLite（modernc.org/sqlite） |
| 认证 | JWT（go-zero 内置） |
| 静态资源 | Go embed |
| 构建工具 | Vite |

## 开发

### 前端开发

```bash
cd web
npm install
npm run dev          # 启动开发服务器（代理API到localhost:8888）
```

### 后端开发

```bash
cd server
go run todo.go -f etc/todo-api.yaml
```

### 重新构建前端到二进制

```bash
cd web && npm run build && cd ../server && go build -o todo-api .
```

## License

MIT
