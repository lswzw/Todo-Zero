# Todo-Zero 已完成开发记录

> 从 `todo-next.md` 迁出的已完成项，按版本归档。后续完成新项时同步更新本文档和 todo-next。

---

## v1.1.0 — P0 必须修复（影响正确性/性能）

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

## v1.2.0 — 代码审查 HIGH 修复

> 详见 [code-review.md](./code-review.md)，共 39 个发现。

- [x] **#1 管理员路由无 RBAC 中间件** — 已创建 `AdminMiddleware`，在路由级别统一拦截，移除各 logic 中重复的 `checkAdmin()`（7 处）。
- [x] **#2 JWT Secret 硬编码默认值** — 首次启动自动生成 64 位随机密钥持久化到 `system_configs`，后续启动从数据库加载。
- [x] **#3 操作日志系统形同虚设** — 已创建 `OperationLogMiddleware`，自动记录 admin 路由的所有写操作（POST/PUT/PATCH/DELETE）。

---

## v1.3.0 — 代码审查 MEDIUM + LOW 修复

### MEDIUM

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

### LOW

- [x] **#22** checkAdmin 方法重复 7 次 — 已通过 AdminMiddleware 统一拦截，移除所有重复 checkAdmin()
- [x] **#25** 数据库连接未在关闭时清理 — 已添加 `defer sqliteDB.Close()`
- [x] **#26** init.sql 中 PRAGMA 重复 — init.sql 中的 PRAGMA 已移除，仅 InitDB 执行

### P1 同期完成

- [x] **操作日志自动记录** — 通过 `OperationLogMiddleware` 实现，与 #3 同步完成。
- [x] **JWT Secret 安全提醒** — 通过自动生成+持久化实现，与 #2 同步完成，不再需要手动修改配置。

---

## 已解决的技术债

| 问题 | 位置 | 解决版本 | 说明 |
|------|------|---------|------|
| 数据库连接池未配置 | `db/init.go` | v1.3.0 | `SetMaxOpenConns(1)` / `SetMaxIdleConns(1)` |
| 数据库连接未关闭 | `todo.go` | v1.3.0 | `defer sqliteDB.Close()` |
| init.sql PRAGMA 重复 | `db/init.sql` | v1.3.0 | 移除重复 PRAGMA，仅 InitDB 执行 |

---

## v1.4.0 — 代码审查 MEDIUM 修复（前端安全 + 类型化 + 后端配置优先级）

- [x] **#15 前端 Token 存 localStorage 有 XSS 风险** — `stores/user.ts` 不再将 `isAdmin` 存入 localStorage，改为每次路由跳转从服务端 `getUserInfo` 接口获取真实身份，防止本地篡改。
- [x] **#16 前端 isAdmin 可被绕过** — 路由守卫增加 `fetchUserInfo()` 服务端验证机制，首次进入需认证页面时从 API 获取真实 `isAdmin`，不再信任 localStorage 中的值。登出时调用 `resetAuthVerified()` 重置验证状态。
- [x] **#17 大量 any 类型** — 新增 `web/src/types/index.ts` 定义全部 API 响应类型（`TaskItem`、`UserListItem`、`StatResp` 等 15 个接口）。所有 `.vue` 和 `.ts` 文件中的 `any` 替换为具体类型；`FormInstance` 替代无类型 `ref()`；validator 回调类型 `any` → `(error?: Error) => void`。
- [x] **#18 空 catch 块吞掉所有错误** — 关键数据加载（统计、分类、任务列表、用户列表、配置等）添加 `ElMessage.error` 提示；操作类 catch 保留注释说明错误已由拦截器处理。
- [x] **#19 命令行 flag 与配置文件优先级不清** — `todo.go` 重构配置加载逻辑：先加载配置文件（或嵌入式默认），再统一由 flag 覆盖非默认值，无论是否使用 `-f` 指定配置文件，flag 始终可以覆盖。
- [x] **#20 API 响应拦截器未统一解包** — `request.ts` 响应拦截器改为：成功时（`code === 0`）直接返回 `data` 字段，调用方无需再 `.data`；业务错误在拦截器统一 `ElMessage.error` 并 reject，调用方不再需要手动处理。

---

## v1.5.0 — 输入验证修复 + LOW 审查项修复

### HIGH — 输入验证 & 注入安全审计

- [x] **#40 `validate` 标签未在运行时执行** — 删除无效 `validate` 标签，为每个结构体实现 `Validate()` 方法（go-zero `validation.Validator` 接口），`httpx.Parse` 自动调用。
- [x] **#41 `options` 标签完全无效** — 将 `options:"xxx"` 独立标签改为 go-zero 原生的 `json/form:"xxx,options=xxx"` 内联语法，框架自动校验枚举值。
- [x] **#42 `BatchTaskReq.Action` 无枚举校验** — `Validate()` 方法校验枚举 + logic 层双重校验，非法 Action 返回错误而非静默忽略。

### LOW — 代码审查（v1.2.0 审查）

- [x] **#21** Register TOCTOU 竞态 — 去掉 FindOneByUsername 预检查，直接 Insert 捕获 UNIQUE 约束错误映射为 UserAlreadyExist
- [x] **#23** CategoryModel.FindById 不过滤 is_deleted — 已验证：categories 表无 is_deleted 字段，Delete 为硬删除，无需过滤
- [x] **#24** UserModel.Update 不更新 password — 已验证：ToggleUserStatusLogic 已使用 UpdateStatus，Update 不含 password 是正确设计
- [x] **#27** etc/todo-api.yaml 中 AccessSecret 明文 — 添加注释说明首次启动自动生成，yaml 仅为 fallback
- [x] **#28** 静态文件路径遍历（低风险） — 已验证：Go embed.FS 天然防路径遍历，无需修改
- [x] **#29** 登录/注册成功后硬编码延迟跳转 — 已验证：当前代码无 setTimeout，立即 router.push
- [x] **#30** 管理员登录后强制跳转管理页面 — 管理员登录后也跳转首页，通过导航栏进入管理
- [x] **#31** 筛选条件改变时未重置页码 — 筛选 @change 触发 onFilterChange 重置 page=1 再 loadTasks
- [x] **#32** v-model:current-page 与 @current-change 冲突 — 移除 @current-change，改用 watch(page) 加载数据
- [x] **#33** config.vue el-switch 状态管理 — 移除共享 switchValue，改为 `:model-value="item._value === 'true'"` 绑定
- [x] **#34** 密码确认验证器类型不安全 — 已验证：当前签名 `(error?: Error) => void` 比 `any` 更安全，无需修改
- [x] **#36** 默认管理员密码 admin123 — 移除 init.sql 中的明文密码注释
- [x] **#37** Model 层冗余别名方法 — 删除 CategoryModel.FindById、UserModel.FindById、SystemConfigModel.FindOneByKey 及实现
- [x] **#38** TaskModel.CountStats 未被使用 — StatLogic 改用 CountStats SQL 聚合，不再 FindList 全量加载
- [x] **#39** TaskModel 多个方法未被使用 — 删除 FindByUserId、FindByCategoryId、CountByStatus、BatchDelete 方法及实现

---

## v1.6.0 — 输入验证 MEDIUM + LOW 修复 + 状态筛选 bug 修复

### MEDIUM — 输入验证 & 注入安全审计

- [x] **#43 多个查询/过滤参数缺少长度验证** — TaskListReq/UserListReq Keyword max=50，LoginLogReq/OperationLogReq Username max=20，OperationLogReq Action max=20，分页参数 1-100 限制，新增 TaskListReq/LoginLogReq/OperationLogReq/UserListReq 的 Validate() 方法
- [x] **#44 `UpdateConfigReq.Value` 无长度/格式限制** — Key max=50 长度校验，Value max=500（原有）
- [x] **#45 `CreateTaskReq.CategoryId` 无验证** — Validate() 校验 CategoryId >= 0，负数拒绝；UpdateTaskReq 同步校验
- [x] **#46 `LoginReq.Username` max=100 过大** — LoginReq 用户名 max 降为 50（注册限制 3-20，登录放宽至 50 即可）
- [x] **#47 密码无复杂度要求** — 后端 validatePassword 增加字母+数字校验，前端注册/改密/重置密码均添加 pattern 正则提示
- [x] **#48 后端缺少输出编码** — 已验证：Go `encoding/json` 默认对 `<>&` 做 Unicode 转义（`\u003c`），JSON API 天然防 XSS，无需额外编码

### LOW — 输入验证 & 注入安全审计

- [x] **#49 API 响应路径缺少安全头** — 新增 SecurityHeadersMiddleware，所有 API 路由添加 X-Content-Type-Options/X-Frame-Options/X-XSS-Protection/Referrer-Policy

### Bug 修复

- [x] **状态筛选"待办"无效** — 根因：go-zero `form:"status,optional"` 对 int64 零值无法区分"未传"和"传了0"，Status/Priority/CategoryId 改为 `default=-1`，tasklistlogic 正确识别 `status==0` 为"待办"筛选

### 其他

- [x] **首页导航栏添加"管理后台"入口** — 管理员登录后首页导航栏显示"管理后台"按钮（`v-if="userStore.isAdmin"`），普通用户不可见
- [x] **#50 SQL 表名拼接编码习惯** — 已验证：`tableName()` 返回硬编码常量，当前安全；编码风格问题留待重构
- [x] **staticcheck SA5008 配置** — 添加 `.staticcheck.conf` 排除 go-zero 扩展 json tag 选项误报

---

## v1.4.0 — 登录限流 + 分类 CRUD + 健康检查

### 安全性

- [x] **登录限流** — 新增 `LoginRateLimitMiddleware`，同 IP 15 分钟内最多 5 次登录尝试，超出后锁定 15 分钟；返回 HTTP 429 + `Retry-After` 头；后台 goroutine 定期清理过期记录防止内存泄漏
- [x] **登录限流错误响应** — `xerr.ErrorResponse` 增加 429 状态码处理，限流错误返回 `42901` code

### 功能补全

- [x] **补全分类 CRUD** — 新增 `PUT /api/v1/category/:id`（更新分类）和 `DELETE /api/v1/category/:id`（删除分类）API 端点
  - 更新分类：支持修改 name/color/icon/sort 字段，系统分类不可修改
  - 删除分类：仅可删除自己的非系统分类，系统分类拒绝删除
  - 权限校验：只能操作自己的分类，系统分类（isSystem=1）受保护
  - 输入验证：`UpdateCategoryReq`/`DeleteCategoryReq` 均实现 `Validate()` 方法
- [x] **分类颜色支持** — `CreateCategoryReq` 新增 `color` 可选字段；`CategoryItem` 返回 color/icon/sort/isSystem 完整信息；`CreateCategoryLogic` 支持自定义颜色
- [x] **分类管理 UI** — 首页任务栏添加"分类管理"按钮，弹窗支持：
  - 添加分类（名称 + 颜色选择器）
  - 编辑分类名称/颜色（内联编辑，失焦保存）
  - 删除分类（确认提示，系统分类不可删除）
  - 分类标签按颜色渲染（自动亮度检测决定文字颜色）

### 基础设施

- [x] **健康检查端点** — 新增 `GET /health`，无 JWT/中间件要求，ping DB 验证连通性，返回 `{"status":"ok"}`
- [x] **ServiceContext 暴露 DB** — 新增 `DB *sql.DB` 字段，供健康检查等场景使用

---

## v1.7.0 — 补全任务时间/标签字段

### 功能补全

- [x] **补全任务时间/标签字段** — 数据库已有 `start_time`、`end_time`、`reminder`、`tags` 字段，从 API 到前端全链路暴露
  - API 层：`CreateTaskReq`、`UpdateTaskReq`、`TaskDetailResp`、`TaskItem` 新增 4 个字段
  - Go types：对应结构体同步新增，`UpdateTaskReq` 时间/标签字段使用 `*string` 指针类型（区分"未提供"和"清空"）
  - 输入验证：时间字段校验 `2006-01-02 15:04` 格式，标签 max=200 字符
  - Create Logic：`parseNullTime` 辅助函数将字符串转为 `sql.NullTime`，写入 DB
  - Update Logic：`*string` 指针映射到 `sql.NullTime`/`string`，支持清空操作
  - Detail/List Logic：`formatNullTime` 辅助函数将 `sql.NullTime` 转为字符串返回
  - 前端 TS 类型：`TaskItem` 新增 4 字段，新增 `TaskFormData` 接口
  - 前端表单：新增 3 个 `el-date-picker`（开始时间/截止时间/提醒时间）+ 标签输入框
  - 前端列表：标签逗号拆分为独立 `el-tag` 展示；截止时间逾期红色高亮

---

## v1.8.0 — 任务详情页 + 单元测试 + 空 catch 修复

### P1 — 功能补全

- [x] **补全任务详情页** — 新增 `task-detail.vue` 详情页，路由 `/task/:id`，首页任务卡片点击跳转
  - 顶部导航栏（返回列表 + 编辑/切换状态/删除操作）
  - 加载中/错误/正常三种状态展示
  - 标题区（带状态圆圈，可点击切换）
  - 元信息行（优先级、分类、标签、状态标签）
  - 信息卡片网格（创建/开始/截止/提醒/更新时间）
  - 内容区（完整展示，支持换行）
  - 编辑弹窗（复用表单逻辑）
  - 截止时间过期高亮

- [x] **添加 Go 单元测试** — model 层 25 个测试 + logic 层 5 个测试，覆盖 CRUD、软删除、权限校验、密码验证等
  - `model/test_helper_test.go` — 测试辅助：SQLite 内存数据库初始化、测试数据插入
  - `model/taskmodel_test.go` (7) — Insert/FindOne/SoftDelete/Update/UpdateStatus/FindList(过滤+分页)/CountStats
  - `model/usermodel_test.go` (8) — Insert/FindOne/FindByUsername/SoftDelete/Update/UpdateStatus/UpdatePassword/FindList/UNIQUE 约束
  - `model/categorymodel_test.go` (8) — Insert/FindOne/系统分类不可删/用户分类可删/FindAll/FindSystem/Update/CountByUser
  - `logic/task/taskdetaillogic_test.go` (5) — 成功获取/无权限/任务不存在/缺少 userId/无分类
  - `logic/user/loginlogic_test.go` (5) — 登录成功/用户不存在/用户禁用/密码错误/管理员登录

### P2 — 空 catch 修复

- [x] **空 catch 修复** — 所有空 catch 添加 `ElMessage.error` 错误提示，ElMessageBox 取消与 API 错误分离
  - `views/home.vue` (9 处) — 切换状态/删除/批量/提交任务/改密码/添加分类/更新分类/删除分类
  - `views/task-detail.vue` (3 处) — 切换状态/删除/编辑提交
  - `views/admin/user.vue` (3 处) — 重置密码/切换状态/删除用户
  - `views/admin/config.vue` (1 处) — 保存配置
  - `views/admin/log.vue` (1 处) — 加载操作日志
  - `views/admin/login-log.vue` (1 处) — 加载登录日志
  - 保留静默的 4 处（合理设计）：checkRegister 非关键、登录错误由拦截器处理、token 失效执行 logout、loadCategories 降级

---

## v1.9.0 — 定时清理任务 + 日志自动清理

### P2 — 后台定时清理

- [x] **实现定时清理任务** — 后台 goroutine 每小时检查，根据系统配置自动清理
  - `task_auto_delete_days > 0` 时，永久删除已完成任务（status=2, is_deleted=0）且 `update_time` 超过 N 天的记录
  - `task_trash_retention_days > 0` 时（默认 30 天），永久删除软删除任务（is_deleted=1）且 `update_time` 超过 N 天的记录
  - 新增 `task_trash_retention_days` 系统配置项（默认 30 天）
  - 配置项每次运行时从数据库热加载，修改后无需重启服务

- [x] **日志自动清理** — 操作日志和登录日志共享 `log_auto_delete_days` 配置
  - `log_auto_delete_days > 0` 时，删除 `created_at` / `create_time` 超过 N 天的日志
  - 新增 `log_auto_delete_days` 系统配置项（默认 0 = 不清理）

### Model 层新增方法

- `TaskModel.HardDeleteCompletedBefore(ctx, beforeTime)` — 永久删除指定时间前的已完成任务
- `TaskModel.HardDeleteSoftDeletedBefore(ctx, beforeTime)` — 永久删除指定时间前的软删除任务
- `OperationLogModel.DeleteOlderThan(ctx, beforeTime)` — 删除指定时间前的操作日志
- `LoginLogModel.DeleteOlderThan(ctx, beforeTime)` — 删除指定时间前的登录日志

### 数据库索引

- `idx_tasks_completed_cleanup` ON tasks (status, is_deleted, update_time)
- `idx_tasks_soft_deleted` ON tasks (is_deleted, update_time)
- `idx_operation_logs_created_at` ON operation_logs (created_at)
- `idx_login_log_create_time` ON login_log (create_time)

### 新增文件

- `scheduler/cleanup.go` — 定时清理调度器
- `scheduler/cleanup_test.go` — 5 个测试覆盖配置读取和各种清理场景

### 前端

- `admin/config.vue` — 新增 `task_auto_delete_days`、`task_trash_retention_days`、`log_auto_delete_days` 配置项中文标题和描述

### 关闭的待办项

- ~~**前端 TypeScript 类型化**~~ — 不适用：v1.4.0 已定义 15 个接口，v1.9.0 清除最后 1 个 `any`，当前 0 个 `any` 残留
- ~~**提取公共组件**~~ — 不适用：项目仅 4 处分页各 1 行 `<el-pagination>`，抽象成本 > 收益
- ~~**数据库迁移机制**~~ — 不适用：SQLite 单文件应用，无老库升级用户，init.sql + ensureIndexes 已覆盖新库初始化和索引补全，引入迁移框架属过度设计
- ~~**健康检查端点**~~ — 已于 v1.4.0 完成 `GET /health`

---

## v1.9.0 — ESLint + Prettier 代码规范

### P2 — 前端代码规范

- [x] **ESLint + Prettier** — 完整的代码规范和格式化配置
  - ESLint 9 flat config + `@eslint/js` + `typescript-eslint` + `eslint-plugin-vue` + `eslint-config-prettier` + `eslint-plugin-prettier`
  - `.prettierrc` — 无分号、单引号、120 字符行宽、尾逗号
  - `eslint.config.js` — TypeScript warn-any、Vue 宽松规则、Prettier 集成
  - npm scripts: `lint` (eslint --fix) / `format` (prettier --write)
  - 全部源文件自动格式化，消除 87 个格式错误和 1 个 `any` 类型警告
  - `stores/user.ts` — `as any` 替换为 `as UserInfo` 类型断言

---

## v2.0.0 — 技术债清理

### 高优先级

- [x] **环境变量管理** — `request.ts` 的 `baseURL` 从硬编码 `/api/v1` 改为 `import.meta.env.VITE_API_BASE_URL`
  - 新增 `.env.production`，`VITE_API_BASE_URL=/api/v1`
  - `.env` 新增 `VITE_API_BASE_URL=/api/v1`
  - `vite.config.ts` 代理目标也改为读取环境变量

- [x] **BatchDelete SQL 拼接** — `LoginLogModel` 补齐 `DeleteBatch` 方法，与 `OperationLogModel` 实现一致
  - 两个 model 均使用参数化占位符 `?`，无 SQL 注入风险
  - 同步更新 `loginlogic_test.go` 和 `cleanup_test.go` 的 mock 对象

### 中优先级

- [x] **SQL 表名拼接** — 全部 6 个 model 文件（44+ 处）从 `` `...` + m.tableName() + `...` `` 统一改为 `fmt.Sprintf`
  - `usermodel_gen.go`、`categorymodel_gen.go`、`systemconfigmodel_gen.go`、`taskmodel_gen.go`、`operationlogmodel_gen.go`、`loginlogmodel_gen.go`

- [x] **配置热更新** — `SystemConfigModel` 增加带 TTL 的内存缓存
  - 使用 `sync.Map` 存储，30 秒 TTL
  - `FindByKey` 命中缓存直接返回，未命中查库后写入
  - `Update`/`Insert` 时自动清除对应 key 的缓存

### 低优先级

- [x] **App.vue 全局错误处理** — 增加 `onErrorCaptured` 全局错误边界

### 测试修复

- `test.sh` 系统配置期望值从 5 更新为 7（匹配 init.sql 实际配置项数）

---

## v2.0.0 — P3-A 数据安全与部署

### Docker 镜像

- [x] **Docker 镜像** — 三阶段构建（Node 编译前端 → Go 编译后端 → Alpine 运行），镜像仅 31MB
  - `Dockerfile` — 多阶段构建，前端产物输出到 `server/dist`，Go 静态编译，最终 Alpine 运行
  - `.dockerignore` — 排除无关文件加速构建
  - 运行：`docker run -d -p 8888:8888 -v todo-data:/app/data todo-zero`
  - 自定义：`docker run -d -p 9090:9090 -v todo-data:/app/data todo-zero -port 9090 -jwt-secret your-secret`

### API 文档

- [x] **API 文档** — 基于框架能力生成 OpenAPI 规范 + 交互式文档 UI
  - `goctl api doc` 从 `.api` 文件自动生成 Markdown 文档（`docs/api/doc/todo.md`）
  - `goctl-swagger` 插件生成 Swagger 2.0 JSON（`docs/api/openapi.json`）
  - `.api` 文件添加 `@doc` 注解，接口标题完整显示
  - 修正 `options` 标签格式（`options=` → `options:`），解决插件兼容问题
  - `server/docs/openapi.json` 通过 `//go:embed` 嵌入二进制
  - `GET /api-docs` — Scalar 交互式 API 文档页面
  - `GET /openapi.json` — Swagger 2.0 规范 JSON
  - **仅 `-debug` 模式开启**，生产环境不暴露
  - Config 新增 `Debug bool` 字段，命令行 `-debug` 参数

---

## v2.1.0 — P3-B 实用功能

### 数据导出

- [x] **导出任务为 CSV/JSON** — 新增 `GET /api/v1/task/export` 端点
  - 支持 `format=json`（默认）和 `format=csv` 两种导出格式
  - CSV 使用 UTF-8 BOM 编码，中文列名表头（ID/标题/内容/状态/优先级/分类/开始时间/截止时间/提醒/标签）
  - 支持与任务列表相同的筛选条件（keyword/status/priority/categoryId）
  - 仅导出当前用户的未删除任务，权限隔离与列表接口一致
  - 新增 `TaskModel.FindAllForExport` 方法（无分页，全量查询）
  - 前端首页添加"导出"按钮，支持 JSON/CSV 格式选择

### 拖拽排序

- [x] **任务列表拖拽排序** — 新增 `PUT /api/v1/task/sort` 端点
  - 数据库新增 `sort_order` 字段（默认 0），任务列表按 `sort_order ASC, id DESC` 排序
  - `SortTaskReq` 接收 `orders` 数组（id + sortOrder），批量更新排序
  - 权限校验：排序前逐个验证任务归属，跨用户排序返回 40001
  - 空数组请求拒绝（RequestParamError）
  - 前端集成 `vuedraggable@next`，拖拽手柄 + ghost 样式，拖拽结束后自动同步排序
  - `TaskItem`、`TaskDetailResp`、`TrashItem` 响应均包含 `sortOrder` 字段

### 其他

- [x] **test.sh 集成测试更新** — 新增拖拽排序测试区块（6 项：排序成功、sortOrder 验证、空 orders 拒绝、跨用户拒绝、无 token 拒绝）
- [x] **单元测试 mock 修复** — `taskdetaillogic_test.go`、`cleanup_test.go` mock 补齐 `FindAllForExport`、`UpdateSortOrder` 方法

---

## v2.2.0 — P4-A 国际化 (i18n)

### 基础设施

- [x] **安装 vue-i18n@9** — `npm install vue-i18n`
- [x] **创建翻译文件** — `src/locales/zh-CN.ts` + `src/locales/en.ts`，按模块拆分：`common`、`auth`、`home`、`taskDetail`、`trash`、`admin`、`config`、`log`、`backup`
- [x] **注册 i18n 插件** — `src/locales/index.ts` 创建 i18n 实例，`main.ts` 中 `app.use(i18n)`，移除硬编码 `zh-cn` import
- [x] **Element Plus 国际化联动** — `App.vue` 使用 `el-config-provider` + `:locale` 动态绑定 Element Plus locale，与 vue-i18n 语言同步切换
- [x] **语言持久化** — 语言偏好存入 `localStorage('locale')`，刷新后恢复
- [x] **语言切换 composable** — `src/composables/useLocale.ts`，封装 `setLocale`/`currentLocale`/`elementLocale`/`localeOptions`

### 中文提取（前端全量）

- [x] **提取 `request.ts`** — 4 处：网络错误、登录过期、请求失败（成功/错误分支各 2 处）
- [x] **提取 `login.vue`** — ~10 处：表单标签、placeholder、验证消息、登录成功提示、注册链接、语言切换
- [x] **提取 `register.vue`** — ~20 处：表单标签、placeholder、验证消息（含密码不一致/复杂度）、注册成功提示、登录链接、注册关闭页、语言切换
- [x] **提取 `home.vue`** — ~80 处：统计标签、状态/优先级/分类筛选、任务列表、批量操作、新增/编辑弹窗、修改密码弹窗、分类管理弹窗、导出、拖拽、所有 ElMessage/ElMessageBox
- [x] **提取 `task-detail.vue`** — ~45 处：导航栏、详情页标签、信息卡片、编辑弹窗、所有提示消息
- [x] **提取 `trash.vue`** — ~25 处：导航栏、列表、批量操作、永久删除确认、所有提示消息
- [x] **提取 `admin/layout.vue`** — ~10 处：侧边栏菜单、导航栏、退出登录
- [x] **提取 `admin/config.vue`** — ~20 处：configMeta 标题/描述（8 个配置项）、开关 active/inactive 文本、保存提示
- [x] **提取 `admin/user.vue`** — ~20 处：表头、角色/状态标签、重置密码弹窗、所有提示消息
- [x] **提取 `admin/log.vue`** — ~15 处：表头、操作类型筛选、所有提示消息
- [x] **提取 `admin/login-log.vue`** — ~10 处：表头、状态标签、所有提示消息
- [x] **提取 `admin/backup.vue`** — ~25 处：备份状态、列表、恢复对话框（危险操作提示、确认输入、所有提示消息）

### 语言切换入口

- 登录页/注册页 — `el-select` 下拉切换（中文/English）
- 首页导航栏 — `el-select` 下拉切换
- 管理后台导航栏 — `el-select` 下拉切换

### 技术细节

- 表单验证消息使用 `() => t('key')` 函数式写法，确保语言切换后验证消息也更新
- `config.vue` 的 `configMeta` 使用 `computed` + `tm`/`rt` API 实现响应式翻译
- `backup.vue` 的确认输入框使用 `t('backup.confirmRestore')` 作为验证值，中英文切换时匹配对应文本
- `request.ts` 使用 `i18n.global.t()` 全局方式调用（非组件内）

---

## v2.3.0 — P4-B 性能优化

### P1 — 中优先级

- [x] **并发请求优化** — `home.vue` 操作后同时刷新任务列表和统计，改为 `Promise.all([loadTasks(), loadStat()])` 并发执行
  - `handleToggle`、`handleDelete`、`handleBatch`、`handleSubmitTask` 均适用
  - `onMounted` 中 `loadStat()` + `loadTasks()` 同理

- [x] **Vite 构建分包** — `vite.config.ts` 添加 `rollupOptions.output.manualChunks`
  - `vendor` — vue/vue-router/pinia/axios
  - `element-plus` — Element Plus 组件库
  - `icons` — Element Plus Icons
  - 分离后利于浏览器缓存，首屏仅加载必要 chunk

### P2 — 低优先级

- [x] **批量删除并行化** — `trash.vue` 中 `handleBatchPermanentDelete` 使用 `for...of` 逐个 `await`，改为 `Promise.allSettled` 并行删除

- [x] **搜索防抖预留** — 搜索框添加注释：若将来改为 `@input` 实时搜索，需引入 debounce（`lodash-es/debounce` 或手写）

- [x] **虚拟滚动预留** — 分页区域添加注释：若未来增大分页或取消分页，需引入 `vue-virtual-scroller` 或 `@tanstack/vue-virtual`

---

## v2.4.0 — 代码质量与安全加固

- [x] **README 文档重构** — 重新编写 README.md，添加 emoji 图标、完整的 API 接口列表、页面路由、技术栈说明等
- [x] **清理无效目录** — 删除 web/src/assets、web/src/components、web/src/utils 三个空目录
- [x] **测试修复** — 修复 `test_helper_test.go` 中 `tasks` 表缺失 `sort_order` 字段问题，所有 34 个单元测试通过
- [x] **ESLint 配置完善** — 添加浏览器全局变量（window、document、localStorage、Blob、setTimeout），消除 9 个 no-undef 错误
- [x] **前端 localStorage XSS 风险修复** — 创建 `useStorage` 组合式函数，替换 `stores/user.ts` 中直接的 localStorage 操作
- [x] **API 速率限制中间件** — 新增 `APIRateLimitMiddleware`，限制每个 IP 每分钟最多 100 次 API 请求，添加到 category、stat、task、user 等路由组
- [x] **任务提醒功能** — 使用浏览器 Notification API 实现任务提醒，支持权限请求、定时检查、点击跳转等功能
- [x] **标签管理功能** — 实现完整的标签 CRUD 功能，包括数据库表、模型、Handler、Logic、API 和前端集成
- [x] **前端懒加载** — 路由级别代码分割，按需加载组件
- [x] **数据库查询优化** — 添加更多索引，优化复杂查询
  - 新增 tasks 表索引: priority, user_category, end_time
  - 新增 users 表索引: status
  - 新增 login_log 表索引: user_time
  - 新增 operation_logs 表索引: user_created_at
  - 新增 categories 表索引: user_system
  - 新增 tags 表索引: user_id
  - 新增 task_tags 表索引: task_id, tag_id
  - 新增 system_configs 唯一索引: config_key

---

## v2.5.0 — 安全评审修复

> 安全评审完成于 2026-04-29，审计范围：server/internal 下所有后端 Go 代码（120+ 文件）
>
> 汇总：HIGH 13 个、MEDIUM 27 个、LOW 14 个，共 54 项，全部已修复

### HIGH — 高危问题（13 项）

- [x] **H1 密码明文存储** — `Insert`/`UpdatePassword` 改用 bcrypt 加密 | `usermodel_gen.go`
- [x] **H2 JWT Secret 硬编码** — 首次启动自动生成随机密钥并持久化到数据库 | `config.go`, `todo-api.yaml`
- [x] **H3 Token 过期验证缺失** — 添加注释说明 JWT 验证由中间件层处理 | `jwtx/jwt.go`
- [x] **H4 ResetPassword 管理员权限缺失** — 添加管理员身份验证 | `resetpasswordlogic.go`
- [x] **H5 ToggleUserStatus 管理员权限缺失** — 添加管理员身份验证 | `toggleuserstatuslogic.go`
- [x] **H6 SortTaskReq 验证缺失** — 实现 `Validate()` 方法 | `types_validate.go`
- [x] **H7 DownloadBackupReq 路径遍历** — 验证 FileName 防止目录穿越 | `types_validate.go`
- [x] **H8 RestoreBackupReq 路径遍历** — 验证 FileName 防止目录穿越 | `types_validate.go`
- [x] **H9 SQL 注入风险** — 动态表名使用白名单验证 | `scheduler/backup.go`
- [x] **H10 敏感配置泄露** — 过滤/脱敏敏感配置项返回 | `configlistlogic.go`
- [x] **H11 路径遍历 TOCTOU 漏洞** — 使用 `os.Open` 直接打开文件 | `downloadbackuphandler.go`
- [x] **H12 X-Forwarded-For IP 欺骗** — 只信任配置的可信代理 IP | `middleware/ip.go`
- [x] **H13 错误响应泄露内部详情** — 未知错误返回通用消息 | `xerr/response.go`

### MEDIUM — 中危问题（27 项）

**密码安全**

- [x] **M1 密码强度验证** — 注册时验证密码复杂度（8位以上、含大小写/数字/特殊字符） | `types_validate.go`
- [x] **M2 统一登录错误提示** — 对外返回"用户名或密码错误"，不区分具体原因 | `loginlogic.go`
- [x] **M3 备份失败错误信息泄露** — 返回通用错误消息 | `triggerbackuplogic.go`
- [x] **M4 恢复失败错误信息泄露** — 返回通用错误消息 | `restorebackuplogic.go`
- [x] **M5 ChangePasswordReq 密码长度验证** | `types_validate.go`
- [x] **M6 LoginReq Password 长度验证** | `types_validate.go`

**输入验证缺失**

- [x] **M7-M12** 补全 DeleteTaskReq/RestoreTaskReq/PermanentDeleteTaskReq/DeleteUserReq/TaskDetailReq/ToggleTaskReq 的 `Validate()` 方法 | `types_validate.go`
- [x] **M13 ToggleUserStatusReq Validate()** | `types_validate.go`
- [x] **M14 UpdateTagReq Validate()** | `types_validate.go`
- [x] **M15 DeleteTagReq Validate()** | `types_validate.go`
- [x] **M16 CreateTagReq Name/Color 长度验证** | `types_validate.go`
- [x] **M17 ExportTaskReq Keyword 长度限制** | `types_validate.go`

**权限/认证**

- [x] **M18 CategoryList 忽略 JWT 错误** — 修复 JWT 解析错误未正确处理 | `categorylistlogic.go`
- [x] **M19 DeleteUser 忽略 JWT 错误** — 修复 JWT 解析错误未正确处理 | `deleteuserlogic.go`

**其他中危**

- [x] **M20 操作日志 Params 脱敏** | `loginlogmodel_gen.go`
- [x] **M21 审计日志 IP 脱敏** | `operationlogmiddleware.go`
- [x] **M22 X-RateLimit-Remaining 头部乱码** | `apiratelimitmiddleware.go`
- [x] **M23 备份文件权限收紧为 0600** | `scheduler/backup.go`
- [x] **M24 float64 转 int64 类型混淆** | `jwtx/jwt.go`
- [x] **M25 登录失败次数锁定机制** — 新增 `LoginRateLimitMiddleware`，可配置尝试次数/窗口/锁定时长 | `routes.go`, `loginlogic.go`
- [x] **M26 /user/check-register 限流** | `routes.go`
- [x] **M27 Debug 模式敏感信息泄露风险** | `config.go`

### LOW — 低危问题（14 项）

- [x] **L1** 修正 `vars.go` 密码加密注释与实现不符 | `model/vars.go`
- [x] **L2** `FindList` 返回密码字段 — 查询排除 password | `usermodel_gen.go`
- [x] **L3** 系统配置缓存访问控制 | `systemconfigmodel_gen.go`
- [x] **L4** 简化错误消息，避免泄露敏感信息 | `xerr/code.go`
- [x] **L5** 配置化中间件安全参数 — RateLimit 配置项 | `servicecontext.go`, `todo-api.yaml`
- [x] **L6** 清理日志中的敏感操作信息 | `scheduler/cleanup.go`
- [x] **L7** 修复测试文件接口信息暴露 | `cleanup_test.go`
- [x] **L8** /user/check-register 和 /user/register 限流 | `routes.go`
- [x] **L9** 改密额外验证（当前密码验证） | `changepasswordlogic.go`
- [x] **L10** 文件扩展名大小写绕过问题 | `downloadbackuphandler.go`
- [x] **L11** 软删除恢复权限控制 — `Restore`/`PermanentDelete` 增加 userId 参数 | `taskmodel_gen.go`, `usermodel_gen.go`
- [x] **L12** 统一 LoginLogReq Username 长度验证 | `types.go`
- [x] **L13** Keyword 查询字段长度限制 | `types_validate.go`
- [x] **L14** CSP、HSTS 等安全响应头 | `securityheadsmiddleware.go`

### 配置变更

- `etc/todo-api.yaml` 新增 `RateLimit` 配置块：
  - `APIRequestsPerMinute` — API 限流每分钟最大请求数（默认 60）
  - `LoginAttempts` — 登录限流最大尝试次数（默认 10）
  - `LoginWindowMinutes` — 登录统计窗口分钟数（默认 15）
  - `LoginLockoutMinutes` — 登录锁定时长分钟数（默认 15）
