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
