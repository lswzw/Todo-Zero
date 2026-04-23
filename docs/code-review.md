# Todo-Zero v1.1.0 代码审查报告

> 全面代码审查，覆盖后端 Go 代码、前端 Vue/TS 代码、入口配置、架构设计。
> 排除 P0 已修复项（状态映射、N+1 查询、索引、密码哈希泄露）。

---

## 汇总

| 严重性 | 数量 |
|--------|------|
| CRITICAL | 0 |
| HIGH | 3 |
| MEDIUM | 15 |
| LOW | 18 |
| INFO | 3 |
| **总计** | **39** |

| 分类 | 数量 |
|------|------|
| Security | 7 |
| Bug | 14 |
| Performance | 2 |
| CodeStyle | 7 |
| Architecture | 4 |
| Maintainability | 5 |

---

## HIGH — 必须修复

### #1 管理员路由无 RBAC 中间件
- **分类**: Security
- **文件**: `server/internal/handler/routes.go:20-65` + 所有 `logic/admin/*.go`
- **问题**: `/api/v1/admin/*` 仅通过 JWT 保护，每个 logic 手动调用 `checkAdmin()`。新增 admin logic 若忘记调用，任何用户都能越权访问。
- **建议**: 实现路由级 admin 中间件统一拦截。

### #2 JWT Secret 硬编码默认值
- **分类**: Security
- **文件**: `server/todo.go:34`, `server/etc/todo-api.yaml:7`
- **问题**: 默认密钥 `"todo-app-jwt-secret-key-2024"` 硬编码在二进制和配置中，攻击者可伪造任意用户 token。
- **建议**: 首次启动自动生成随机密钥持久化到数据库；启动时若检测默认值打印醒目警告。

### #3 操作日志系统形同虚设
- **分类**: Architecture
- **文件**: 整个 logic 层
- **问题**: `OperationLogModel` 存在，但只有 `RegisterLogic` 写入一条记录。所有管理操作（删除用户、重置密码、修改配置等）均无操作日志。操作日志页面几乎永远为空。
- **建议**: 实现日志中间件或 AOP 切面，自动记录所有写操作。

---

## MEDIUM — 应该修复

### #4 批量操作静默吞掉错误
- **分类**: Bug
- **文件**: `server/internal/logic/task/batchtasklogic.go:38-58`
- **问题**: 循环中 `_ = l.svcCtx.TaskModel.Update/Delete` 错误完全丢弃，前端收到"成功"但可能部分失败。
- **建议**: 收集失败 ID 返回部分成功/失败信息，或使用事务保证原子性。

### #5 ToggleUserStatus 存在 TOCTOU 竞态
- **分类**: Bug/Security
- **文件**: `server/internal/logic/admin/toggleuserstatuslogic.go:33-47`
- **问题**: 先 FindOne 读 status 再内存取反后 Update，并发请求可能导致状态错误覆盖。模型层已有 `UpdateStatus` 方法但未使用。
- **建议**: 用 SQL 原子操作 `UPDATE ... SET status = CASE WHEN status = 1 THEN 0 ELSE 1 END WHERE id = ?`。

### #6 ToggleTask 存在同样的 TOCTOU 竞态
- **分类**: Bug
- **文件**: `server/internal/logic/task/toggletasklogic.go:34-51`
- **问题**: 同 #5，应使用已有的 `TaskModel.UpdateStatus` 做原子更新。

### #7 UpdateTask 无法清空可选字段
- **分类**: Bug
- **文件**: `server/internal/logic/task/updatetasklogic.go:44-58`
- **问题**: `if req.Title != ""` / `if req.CategoryId != 0` 导致无法将字段清空回零值。
- **建议**: 使用指针类型 `*string`/`*int64` 区分"未提供"和"清空"。

### #8 操作日志过滤参数未生效
- **分类**: Bug
- **文件**: `server/internal/logic/admin/operationloglistlogic.go:33`, `server/internal/model/operationlogmodel_gen.go:39`
- **问题**: `OperationLogReq` 定义了 `Action`/`Username` 过滤字段，但 `FindList` 只接受 page/pageSize，完全没有过滤功能。
- **建议**: `FindList` 增加 `action`/`username` 参数支持过滤查询。

### #9 OperationLogItem 映射字段丢失
- **分类**: Bug
- **文件**: `server/internal/logic/admin/operationloglistlogic.go:40-50`
- **问题**: `UserId` 硬编码为 0，`TargetType` 为空字符串。前端"对象类型"列始终为空。
- **建议**: 正确映射数据库字段，或从 API 响应类型中移除无用字段。

### #10 LoginReq 缺少验证标签
- **分类**: Security/Bug
- **文件**: `server/internal/types/types.go:99-102`
- **问题**: `LoginReq` 的 Username/Password 无 `validate` 标签，空请求和超长字符串可导致 bcrypt 消耗大量 CPU。
- **建议**: 添加 `validate:"required,min=1,max=100"`。

### #11 SQL 注释解析逻辑脆弱
- **分类**: Bug
- **文件**: `server/internal/db/init.go:96-120`
- **问题**: `strings.Contains(line, "--")` 会误匹配字符串值中的注释符（如 `INSERT INTO t VALUES('a--b')`），当前 init.sql 碰巧避开但逻辑本身脆弱。
- **建议**: 使用支持多语句的驱动方法，或对初始化 SQL 采用更健壮的解析策略。

### #12 数据库连接池未配置
- **分类**: Performance
- **文件**: `server/internal/db/init.go:26`
- **问题**: `sql.Open` 后未设 `SetMaxOpenConns` 等，长时间运行可能连接泄漏。
- **建议**: 为 SQLite 配置 `SetMaxOpenConns(1)` 或小量连接。

### #13 DeleteUserLogic 未阻止删除其他管理员
- **分类**: Security
- **文件**: `server/internal/logic/admin/deleteuserlogic.go:28-43`
- **问题**: 只有"不能删除自己"的检查，一个管理员可删除所有其他管理员，导致系统无法管理。前端虽禁用按钮但 API 无保护。
- **建议**: 后端增加保护，不允许删除其他管理员账户。

### #14 静态文件处理器缺少安全头
- **分类**: Security
- **文件**: `server/todo.go:108-148`
- **问题**: 缺少 `X-Content-Type-Options: nosniff`、`X-Frame-Options: DENY` 等安全响应头，`getContentType` 默认返回 `application/octet-stream` 配合无 nosniff 存在 MIME 嗅探风险。
- **建议**: 添加安全响应头。

### #15 前端 Token 存 localStorage 有 XSS 风险
- **分类**: Security
- **文件**: `web/src/stores/user.ts:14-16`
- **问题**: JWT token 存 localStorage，任何 XSS 漏洞都能窃取。
- **建议**: 优先使用 HttpOnly cookie，或确保无 XSS 漏洞。

### #16 前端 isAdmin 可被绕过
- **分类**: Security
- **文件**: `web/src/stores/user.ts:8,15,22`, `web/src/router/index.ts:44`
- **问题**: `isAdmin` 存 localStorage，用户可通过 DevTools 修改绕过路由守卫访问管理页面。虽后端有 checkAdmin 但 UI 仍渲染。
- **建议**: 路由守卫通过 API 实时验证用户角色。

### #17 大量 any 类型
- **分类**: Maintainability
- **文件**: `web/src/views/home.vue:198,210` 等多处
- **问题**: 几乎所有响应和表单数据用 `any`，TypeScript 形同虚设。
- **建议**: 定义 Task、User、Category、Stat 等接口类型。

### #18 空 catch 块吞掉所有错误
- **分类**: Bug/Maintainability
- **文件**: 所有 `.vue` 文件中的 `catch {}`
- **问题**: 页面加载失败时无任何提示，用户看到空白。
- **建议**: catch 中至少添加日志或"加载失败"提示。

### #19 命令行 flag 与配置文件优先级不清
- **分类**: Bug
- **文件**: `server/todo.go:56-71`
- **问题**: 指定配置文件时命令行 flags 完全不生效，但 flag 始终被 Parse，用户可能误以为命令行参数会覆盖配置文件。
- **建议**: 统一配置优先级策略并明确文档说明。

### #20 API 响应拦截器未统一解包
- **分类**: Bug/CodeStyle
- **文件**: `web/src/api/request.ts:19`
- **问题**: 拦截器返回 `response.data`（整个 body），调用方需 `.data` 取业务数据，但当前用 `as any` 绕过。成功/错误处理路径不一致。
- **建议**: 拦截器统一解包，成功时返回 `response.data.data`。

---

## LOW — 建议改进

### #21 Register TOCTOU 竞态
- `server/internal/logic/user/registerlogic.go:38-44`
- 先查再插，并发可能报不友好错误。直接 Insert 捕获 UNIQUE 约束更优。

### #22 checkAdmin 方法重复 7 次
- 所有 `logic/admin/*.go` 各自实现完全相同的 `checkAdmin()`，违反 DRY。

### #23 CategoryModel.FindById 不过滤 is_deleted
- `server/internal/model/categorymodel_gen.go:97-105`
- categories 用物理删除而 users/tasks 用软删除，策略不一致。

### #24 UserModel.Update 不更新 password
- `server/internal/model/usermodel_gen.go:63-68`
- Update SQL 不含 password，ToggleUserStatusLogic 应使用专用 UpdateStatus。

### #25 数据库连接未在关闭时清理
- `server/todo.go:74-78`
- 缺少 `defer sqliteDB.Close()`，WAL 文件可能未正确清理。

### #26 init.sql 中 PRAGMA 重复
- `server/internal/db/init.sql:6` + `server/internal/db/init.go:31`
- InitDB 已执行 PRAGMA，init.sql 中多余。

### #27 etc/todo-api.yaml 中 AccessSecret 明文
- 嵌入二进制的默认配置包含敏感值，应运行时从环境变量读取。

### #28 静态文件路径遍历（低风险）
- `server/todo.go:112-113`
- Go embed.FS 本身安全，但切换到非 embed 文件系统时需额外注意。

### #29 登录/注册成功后硬编码延迟跳转
- `web/src/views/login.vue:63-65`, `web/src/views/register.vue:92`
- setTimeout 延迟不必要，loading 状态已防重复提交。

### #30 管理员登录后强制跳转管理页面
- `web/src/views/login.vue:64`
- 管理员可能想先使用 todo 功能。

### #31 筛选条件改变时未重置页码
- `web/src/views/home.vue:45-56`
- 第 3 页切换筛选会请求 page=3，可能返回空。

### #32 v-model:current-page 与 @current-change 冲突
- `web/src/views/home.vue:118-123`
- 可能有时序问题。

### #33 config.vue el-switch 状态管理
- `web/src/views/admin/config.vue:13-25`
- switchValue 单一 ref 无法处理多个布尔配置项，`_value` 不是响应式属性。

### #34 密码确认验证器类型不安全
- `web/src/views/home.vue:226-229`
- `callback: any` 应为 `(error?: Error) => void`。

### #35 无障碍性缺失
- 所有 `.vue` 文件
- 缺少 ARIA 属性、语义化标签、键盘导航支持。

### #36 默认管理员密码 admin123
- `server/internal/db/init.sql:101-104`
- 明文写在注释中，建议首次启动后强制修改。

### #37 Model 层冗余别名方法
- FindOne→FindById、FindOneByKey→FindByKey 等多余包装。

### #38 TaskModel.CountStats 未被使用
- `server/internal/model/taskmodel_gen.go:184-198`
- StatLogic 用 FindList 全量加载计算，应改用更高效的 CountStats。

### #39 TaskModel 多个方法未被使用
- FindByUserId、FindByCategoryId、CountByStatus、BatchDelete、UpdateStatus 已定义但未使用。

---

## 优先修复建议（Top 10）

| 优先级 | 编号 | 问题 | 影响 |
|--------|------|------|------|
| 1 | #2 | JWT Secret 硬编码 | 可伪造任意用户 token |
| 2 | #1 | 管理员路由无中间件 | 新增 logic 遗忘 checkAdmin 即可越权 |
| 3 | #3 | 操作日志系统形同虚设 | 审计功能缺失 |
| 4 | #4 | 批量操作静默吞错 | 数据不一致 |
| 5 | #8 | 操作日志过滤参数未生效 | 前端功能失效 |
| 6 | #5/#6 | TOCTOU 竞态条件 | 并发场景数据错误 |
| 7 | #7 | UpdateTask 无法清空字段 | 功能缺陷 |
| 8 | #14 | 缺少安全响应头 | Web 安全基础防护 |
| 9 | #16 | 前端权限可被绕过 | 需要 API 级别防护 |
| 10 | #17 | 大量 any 类型 | TypeScript 形同虚设 |
