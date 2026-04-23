# Todo-Zero 下一步开发清单

> 基于 v1.0.0 代码审查，P0 已在 v1.1.0 修复。

---

## P0 — 必须修复（影响正确性/性能）

- [x] **前后端状态/优先级映射不一致**
  - Status 统一为: 0=待办, 2=已完成（去掉未使用的"进行中"状态）
  - Priority 统一为: 1=重要, 2=紧急, 3=普通
  - ToggleTask: 0↔2 切换（原来 0↔1，1 不是"已完成"）
  - BatchTask: complete→status=2, undo→status=0（原来 complete→1）
  - 前端筛选/显示/表单全部对齐

- [x] **任务列表 N+1 查询**
  - 改为批量收集 categoryId → 一次查出所有分类名称 → map 查找

- [x] **数据库缺少索引**
  - 添加 `idx_tasks_user_id`、`idx_tasks_status`、`idx_tasks_category_id`、`idx_users_username`、`idx_login_log_username`
  - init.sql 和 ensureIndexes() 双重保障（新库/旧库都有索引）

- [x] **用户列表返回密码哈希**
  - `UserModel.FindList` 查询排除 `password` 字段

---

## P1 — 重要改进（安全性/功能补全）

- [ ] **登录限流** — 防止暴力破解，同 IP / 同用户短时间多次失败则锁定
- [ ] **补全分类 CRUD** — 后端 model 已有 Update/Delete 方法，缺 API 端点和前端界面
- [ ] **补全任务时间/标签字段** — 数据库已有 `start_time`、`end_time`、`reminder`、`tags` 字段，但 API 和前端均未暴露
- [ ] **补全任务详情页** — 后端有 GET /task/:id，前端有 `getTaskDetail` 函数，但 UI 无入口
- [ ] **添加 Go 单元测试** — 项目零 `_test.go` 文件，至少覆盖 model 层和核心 logic
- [ ] **操作日志自动记录** — OperationLogModel 存在但 logic 层未自动调用，只有登录日志有记录
- [ ] **JWT Secret 安全提醒** — 默认值硬编码，启动时若未修改应打印警告

---

## P2 — 体验优化

- [ ] **前端 TypeScript 类型化** — 消除大量 `any` 类型，定义 Task、User 等接口
- [ ] **提取公共组件** — `components/` 目录为空，分页、搜索栏、空状态等应抽为组件
- [ ] **空 catch 修复** — 前端几乎所有 async 函数 `catch {}` 静默吞异常，应添加错误提示
- [ ] **实现定时清理任务** — 系统配置 `task_auto_delete_days=30` 已存在，但无后台清理逻辑
- [ ] **日志自动清理** — 操作日志和登录日志无限增长，需定时清理策略
- [ ] **数据库迁移机制** — 当前只有首次初始化，无版本化 schema 迁移，后续改表困难
- [ ] **健康检查端点** — 添加 `/health` 接口用于部署探针
- [ ] **ESLint + Prettier** — 前端无代码规范和格式化配置

---

## P3 — 锦上添花

- [ ] **暗黑模式** — Element Plus 原生支持，切换成本低
- [ ] **数据导出** — 导出任务为 CSV/JSON
- [ ] **回收站** — 软删除任务可恢复，已删除列表 + 恢复操作
- [ ] **用户资料编辑** — 数据库有 avatar 字段，支持修改昵称/头像
- [ ] **数据库自动备份** — SQLite 文件定时备份
- [ ] **移动端优化** — 管理后台侧边栏折叠、触控交互适配
- [ ] **PWA 离线支持** — Service Worker + 缓存策略
- [ ] **API 文档** — Swagger/OpenAPI 规范自动生成
- [ ] **Docker 镜像** — 多阶段构建，一键容器化部署
- [ ] **拖拽排序** — 任务列表支持拖拽调整顺序

---

## 代码审查发现（v1.1.0 全量审查）

> 详见 [code-review.md](./code-review.md)，共 39 个发现（3 HIGH / 15 MEDIUM / 18 LOW / 3 INFO）。

### HIGH — 必须修复

- [x] **#1 管理员路由无 RBAC 中间件** — 已创建 `AdminMiddleware`，在路由级别统一拦截，移除各 logic 中重复的 `checkAdmin()`（7 处）。
- [x] **#2 JWT Secret 硬编码默认值** — 首次启动自动生成 64 位随机密钥持久化到 `system_configs`，后续启动从数据库加载。
- [x] **#3 操作日志系统形同虚设** — 已创建 `OperationLogMiddleware`，自动记录 admin 路由的所有写操作（POST/PUT/PATCH/DELETE）。

### MEDIUM — 应该修复

- [x] **#4 批量操作静默吞掉错误** — 收集失败 ID 并记录日志，改用 `UpdateStatus` 原子更新。
- [x] **#5 ToggleUserStatus TOCTOU 竞态** — `toggleuserstatuslogic.go:33-47`，已改用 `UpdateStatus` 原子更新。
- [x] **#6 ToggleTask TOCTOU 竞态** — 改用 `UpdateStatus` 原子更新。
- [x] **#7 UpdateTask 无法清空可选字段** — `UpdateTaskReq` 改用指针类型 `*string/*int64`，区分"未提供"(nil)和"清空"。
- [x] **#8 操作日志过滤参数未生效** — `OperationLogModel.FindList` 增加 action/username 过滤参数。
- [x] **#9 OperationLogItem 映射字段丢失** — `operationloglistlogic.go:40-50`，已修复 `UserId` 和 `TargetType`（映射 `Module` 字段）。
- [x] **#10 LoginReq 缺少验证标签** — 添加 `validate:"required,min=1,max=100"`。
- [x] **#11 SQL 注释解析逻辑脆弱** — 重写为 `removeLineComment()`，跟踪单引号字符串状态，不再误匹配。
- [x] **#12 数据库连接池未配置** — 配置 `SetMaxOpenConns(1)` / `SetMaxIdleConns(1)` + `defer sqliteDB.Close()`。
- [x] **#13 DeleteUserLogic 未阻止删除其他管理员** — `deleteuserlogic.go:28-43`，已添加后端保护，禁止删除管理员账户。
- [x] **#14 静态文件处理器缺少安全头** — 添加 `X-Content-Type-Options` / `X-Frame-Options` / `X-XSS-Protection` / `Referrer-Policy`。
- [ ] **#15 前端 Token 存 localStorage 有 XSS 风险** — `stores/user.ts:14-16`。
- [ ] **#16 前端 isAdmin 可被绕过** — `stores/user.ts:8,15,22`，localStorage 可修改。
- [ ] **#17 大量 any 类型** — `home.vue:198,210` 等，TypeScript 形同虚设。
- [ ] **#18 空 catch 块吞掉所有错误** — 所有 `.vue` 文件，页面加载失败无提示。
- [ ] **#19 命令行 flag 与配置文件优先级不清** — `todo.go:56-71`。
- [ ] **#20 API 响应拦截器未统一解包** — `request.ts:19`，调用方需 `.data` 取业务数据。

### LOW — 建议改进

- [ ] **#21** Register TOCTOU 竞态 — 并发可能报不友好错误，直接 Insert 捕获 UNIQUE 更优
- [x] **#22** checkAdmin 方法重复 7 次 — 已通过 AdminMiddleware 统一拦截，移除所有重复 checkAdmin()
- [ ] **#23** CategoryModel.FindById 不过滤 is_deleted — 与其他 model 策略不一致
- [ ] **#24** UserModel.Update 不更新 password — ToggleUserStatusLogic 应使用专用 UpdateStatus
- [x] **#25** 数据库连接未在关闭时清理 — 已添加 `defer sqliteDB.Close()`
- [x] **#26** init.sql 中 PRAGMA 重复 — init.sql 中的 PRAGMA 已移除，仅 InitDB 执行
- [ ] **#27** etc/todo-api.yaml 中 AccessSecret 明文 — 嵌入二进制的默认配置含敏感值
- [ ] **#28** 静态文件路径遍历（低风险） — Go embed.FS 安全，但切换 FS 需注意
- [ ] **#29** 登录/注册成功后硬编码延迟跳转 — setTimeout 延迟不必要
- [ ] **#30** 管理员登录后强制跳转管理页面 — 管理员可能想先使用 todo 功能
- [ ] **#31** 筛选条件改变时未重置页码 — 第 3 页切换筛选可能返回空
- [ ] **#32** v-model:current-page 与 @current-change 冲突 — 可能有时序问题
- [ ] **#33** config.vue el-switch 状态管理 — switchValue 无法处理多个布尔配置项
- [ ] **#34** 密码确认验证器类型不安全 — `callback: any` 应为 `(error?: Error) => void`
- [ ] **#35** 无障碍性缺失 — 缺少 ARIA 属性、语义化标签、键盘导航
- [ ] **#36** 默认管理员密码 admin123 — 明文写在注释中，建议首次启动后强制修改
- [ ] **#37** Model 层冗余别名方法 — FindOne→FindById 等多余包装
- [ ] **#38** TaskModel.CountStats 未被使用 — StatLogic 用 FindList 全量计算，应改用 CountStats
- [ ] **#39** TaskModel 多个方法未被使用 — FindByUserId、FindByCategoryId 等已定义未使用

---

## 技术债备注

| 问题 | 位置 | 说明 |
|------|------|------|
| 数据库连接池未配置 | `db/init.go` | `sql.Open` 后未设 `SetMaxOpenConns` 等 |
| 配置热更新缺失 | admin config | 修改系统配置后需实时读取数据库 |
| 环境变量管理 | `request.ts` | API 地址硬编码 |
| App.vue 几乎为空 | `App.vue` (41B) | 缺全局错误处理和布局 |
| BatchDelete SQL 拼接 | model 层 | IN 子句用循环拼接占位符，需注意安全性 |
