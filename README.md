# Todo-Zero

> **一个命令，开箱即用。** 零外部依赖的全栈待办管理应用 — Go 后端 + Vue3 前端打包成单一二进制，自带 SQLite 数据库，无需 MySQL、Redis、Nginx。

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vuedotjs)
![License](https://img.shields.io/badge/License-MIT-blue)

## 它能做什么

- 📝 **任务管理** — 创建、编辑、删除、标记完成，支持分类、优先级、标签、截止时间
- 🗂️ **分类 & 标签** — 系统预置 + 自定义分类，支持颜色标记；多标签自由组合
- ♻️ **回收站** — 软删除保护，支持恢复和永久删除
- 📊 **统计概览** — 任务总数、完成率一目了然
- 📥 **数据导出** — 一键导出任务为 CSV 或 JSON
- 🔍 **搜索筛选** — 按关键词、状态、分类、优先级组合筛选
- 👥 **多用户** — 注册登录，数据严格隔离
- 🛡️ **管理后台** — 用户管理、系统配置、操作日志、登录日志、数据备份恢复
- 🌍 **国际化** — 中英文一键切换
- 📱 **响应式** — 桌面、平板、手机自适应
- 🔒 **安全可靠** — bcrypt 加密、JWT 认证、RBAC 权限、API 限流、安全响应头

## 快速开始

### 一条命令运行

```bash
# 编译（前端自动嵌入）
cd web && npm install && npm run build && cd ../server && go build -o todo-api .

# 运行
./todo-api
# → 浏览器打开 http://localhost:8888
```

首次启动自动完成：创建 SQLite 数据库 → 建表 → 创建默认管理员账号。

### Docker

```bash
docker run -d -p 8888:8888 -v todo-data:/app/data todo-zero
```

### 默认账号

| 用户名 | 密码 |
|--------|------|
| admin | admin123 |

## 项目结构

```
.
├── server/                  # Go 后端 (go-zero 框架)
│   ├── todo.go              # 入口：SQLite 初始化 + 嵌入静态文件 + 路由注册
│   ├── etc/todo-api.yaml    # 配置文件
│   └── internal/
│       ├── config/          # 配置加载（命令行 flag 覆盖配置文件）
│       ├── db/              # SQLite 初始化和建表
│       ├── handler/         # HTTP 处理器（goctl 生成 + 自定义）
│       ├── logic/           # 业务逻辑（user / task / category / stat / tag / admin）
│       ├── middleware/      # 中间件（认证、管理员 RBAC、限流、安全头、操作日志）
│       ├── model/           # 数据访问层（原生 SQL）
│       ├── pkg/xerr/        # 统一错误码和响应格式
│       ├── scheduler/       # 定时任务（自动清理过期数据、定期备份）
│       └── svc/             # 服务上下文（依赖注入）
├── web/                     # Vue3 前端 (TypeScript)
│   ├── src/
│   │   ├── api/             # Axios 封装 + API 函数（自动解包响应）
│   │   ├── composables/     # 组合式函数（语言切换等）
│   │   ├── locales/         # 国际化翻译文件（中/英）
│   │   ├── router/          # Vue Router 路由 + 守卫
│   │   ├── stores/          # Pinia 状态管理（用户、语言）
│   │   ├── styles/          # 全局样式变量 + 毛玻璃组件
│   │   ├── types/           # TypeScript 类型定义
│   │   └── views/           # 页面组件
│   │       ├── login.vue           # 登录页
│   │       ├── register.vue        # 注册页
│   │       ├── home.vue            # 任务主页
│   │       ├── task-detail.vue     # 任务详情
│   │       ├── trash.vue           # 回收站
│   │       └── admin/              # 管理后台（用户/配置/日志/备份）
│   └── vite.config.ts       # 构建输出到 server/dist（被 Go embed 嵌入）
├── docs/                    # 文档
│   ├── api/                 # API 定义和 OpenAPI 规范
│   └── development-roadmap.md
└── test.sh                  # 集成测试脚本（覆盖全部 API）
```

## 技术栈

| 层次 | 技术 | 说明 |
|------|------|------|
| 后端 | go-zero | 高性能 Go 框架，内置中间件、代码生成 |
| 前端 | Vue3 + TypeScript | 组合式 API，完整类型定义 |
| UI | Element Plus | Vue3 组件库，国际化联动 |
| 数据库 | SQLite (modernc.org) | 纯 Go 实现，无需 CGO，单文件存储 |
| 认证 | JWT + bcrypt | 无状态认证，密码安全哈希 |
| 部署 | Go embed + Docker | 前端编译产物嵌入二进制；31MB Docker 镜像 |

## API 概览

所有接口前缀 `/api/v1`，需认证的接口携带 `Authorization: Bearer <token>`。

| 模块 | 接口数 | 说明 |
|------|--------|------|
| 用户 | 5 | 注册、登录、获取信息、改密、检查注册开关 |
| 任务 | 11 | CRUD、状态切换、批量操作、回收站、排序、导出、详情 |
| 分类 | 4 | CRUD，系统分类受保护 |
| 标签 | 4 | CRUD，用户隔离 |
| 统计 | 1 | 任务概览（总数/完成率等） |
| 管理 | 12 | 用户管理、系统配置、操作日志、登录日志、备份恢复 |

完整 API 文档：[docs/api/api-docs.md](docs/api/api-docs.md)（项目也可在 `-debug` 模式下通过 `/api-docs` 查看交互式文档）

## 配置

支持**命令行参数**和**配置文件**两种方式，命令行优先级更高。

```bash
./todo-api -port 9090 -data-dir /var/todo    # 命令行
./todo-api -f /etc/todo-api.yaml            # 配置文件
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-host` | `0.0.0.0` | 监听地址 |
| `-port` | `8888` | 监听端口 |
| `-data-dir` | `data` | 数据存储目录 |
| `-jwt-secret` | 自动生成 | JWT 密钥（首次启动自动生成并持久化） |
| `-jwt-expire` | `86400` | Token 有效期（秒） |
| `-debug` | `false` | 调试模式（开启 API 文档等） |

## 开发指南

```bash
# 前端开发（热更新，代理 API 到 localhost:8888）
cd web && npm install && npm run dev

# 后端开发
cd server && go run todo.go

# 重新构建（前端嵌入后端）
cd web && npm run build && cd ../server && go build -o todo-api .

# 运行测试
cd server && go test ./...          # 单元测试（25+ model 测试 + logic 测试）
./test.sh                           # 集成测试（覆盖全部 API）
```

## 版本历史

| 版本 | 亮点 |
|------|------|
| v1.0 | 基础 CRUD、JWT 认证、单体二进制 |
| v1.1 | P0 bug 修复（状态映射、N+1 查询、索引、密码泄露） |
| v1.2 | RBAC 中间件、JWT 自动生成、操作日志自动化 |
| v1.3 | 输入验证全覆盖、安全头、连接池 |
| v1.4 | 前端安全加固、类型化、API 响应统一解包 |
| v1.5 | 验证框架修复、竞态修复、密码复杂度 |
| v1.6 | 状态筛选 bug 修复、安全头补全 |
| v1.7 | 任务时间/标签字段全链路暴露 |
| v1.8 | 任务详情页、Go 单元测试 |
| v1.9 | 定时清理、ESLint/Prettier |
| v2.0 | 环境变量管理、SQL 规范化、配置热更新 |
| v2.1 | 数据导出、拖拽排序、Docker 镜像、API 文档 |
| v2.2 | 国际化 (i18n 中英文) |
| v2.3 | 并发优化、构建分包 |
| v2.4 | 标签管理、API 限流、数据库索引优化、懒加载 |
| v2.5 | 安全评审 54 项全部修复 (HIGH 13 + MEDIUM 27 + LOW 14) |
| v3.0 | 界面现代化 — 毛玻璃/渐变/动画，22 项全部完成 |
| v3.1 | 国际化增强、悬浮语言切换按钮、bug 清理 |

## License

MIT
