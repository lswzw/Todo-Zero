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
