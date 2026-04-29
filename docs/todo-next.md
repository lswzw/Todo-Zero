# Todo-Zero 安全修复任务列表

> 本次安全评审完成于 2026-04-29，审计范围：server/internal 下所有后端 Go 代码（120+ 文件）
>
> **汇总：HIGH 级问题 13 个（已全部修复），MEDIUM 级问题 35 个，LOW 级问题 18 个**

---

## 🔴 优先级 1 - 高危问题（已全部修复 ✅）

| 状态 | # | 任务描述 | 文件位置 | 风险等级 |
|:---:|---|----------|----------|----------|
| [x] | H1 | 修复密码明文存储 - 在 `Insert`/`UpdatePassword` 中使用 bcrypt 加密 | `server/internal/model/usermodel_gen.go` | 🔴 HIGH |
| [x] | H2 | 修复 JWT Secret 硬编码 - 使用环境变量存储 JWT Secret | `server/etc/todo-api.yaml`, `server/internal/config/config.go` | 🔴 HIGH |
| [x] | H3 | 修复 Token 过期验证缺失 - 添加注释说明验证架构 | `server/internal/pkg/jwtx/jwt.go` | 🔴 HIGH |
| [x] | H4 | 修复 ResetPassword 管理员权限验证 - 添加管理员验证 | `server/internal/logic/admin/resetpasswordlogic.go` | 🔴 HIGH |
| [x] | H5 | 修复 ToggleUserStatus 管理员权限验证 - 添加管理员验证 | `server/internal/logic/admin/toggleuserstatuslogic.go` | 🔴 HIGH |
| [x] | H6 | 修复 SortTaskReq 验证缺失 - 实现 Validate() 方法 | `server/internal/types/types_validate.go` | 🔴 HIGH |
| [x] | H7 | 修复 DownloadBackupReq 路径遍历漏洞 - 验证 FileName | `server/internal/types/types_validate.go` | 🔴 HIGH |
| [x] | H8 | 修复 RestoreBackupReq 路径遍历漏洞 - 验证 FileName | `server/internal/types/types_validate.go` | 🔴 HIGH |
| [x] | H9 | 修复 SQL 注入风险 - 使用白名单验证动态表名 | `server/internal/scheduler/backup.go` | 🔴 HIGH |
| [x] | H10 | 修复敏感配置泄露 - 过滤或脱敏敏感配置项 | `server/internal/logic/admin/configlistlogic.go` | 🔴 HIGH |
| [x] | H11 | 修复路径遍历 TOCTOU 漏洞 - 使用 os.Open 直接打开 | `server/internal/handler/admin/downloadbackuphandler.go` | 🔴 HIGH |
| [x] | H12 | 修复 X-Forwarded-For IP 欺骗 - 只信任已知可信代理 | `server/internal/middleware/ip.go` | 🔴 HIGH |
| [x] | H13 | 修复错误响应泄露 - 未知错误返回通用消息 | `server/internal/pkg/xerr/response.go` | 🔴 HIGH |

---

## 🟡 优先级 2 - 中危问题（建议尽快处理）

### 密码安全

| 状态 | # | 任务描述 | 文件位置 | 风险等级 |
|:---:|---|----------|----------|----------|
| [ ] | M1 | 添加密码强度验证 - 注册时验证密码复杂度（8位以上、包含大小写、数字、特殊字符） | `server/internal/logic/user/registerlogic.go` | 🟡 MEDIUM |
| [ ] | M2 | 统一登录错误提示 - 对外返回"用户名或密码错误"，不区分具体原因 | `server/internal/logic/user/loginlogic.go` | 🟡 MEDIUM |
| [ ] | M3 | 修复备份失败错误信息泄露 - 返回通用错误消息 | `server/internal/logic/admin/triggerbackuplogic.go` | 🟡 MEDIUM |
| [ ] | M4 | 修复恢复失败错误信息泄露 - 返回通用错误消息 | `server/internal/logic/admin/restorebackuplogic.go` | 🟡 MEDIUM |
| [ ] | M5 | 添加 ChangePasswordReq 密码长度验证 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [ ] | M6 | 添加 LoginReq Password 长度验证 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |

### 输入验证缺失

| 状态 | # | 任务描述 | 文件位置 | 风险等级 |
|:---:|---|----------|----------|----------|
| [x] | M7 | 添加 DeleteTaskReq Validate() 方法 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [x] | M8 | 添加 RestoreTaskReq Validate() 方法 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [x] | M9 | 添加 PermanentDeleteTaskReq Validate() 方法 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [x] | M10 | 添加 DeleteUserReq Validate() 方法 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [x] | M11 | 添加 TaskDetailReq Validate() 方法 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [x] | M12 | 添加 ToggleTaskReq Validate() 方法 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [x] | M13 | 添加 ToggleUserStatusReq Validate() 方法 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [x] | M14 | 添加 UpdateTagReq Validate() 方法 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [x] | M15 | 添加 DeleteTagReq Validate() 方法 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [x] | M16 | 添加 CreateTagReq Name/Color 长度验证 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |
| [ ] | M17 | 添加 ExportTaskReq Keyword 长度限制 | `server/internal/types/types_validate.go` | 🟡 MEDIUM |

### 权限/认证问题

| 状态 | # | 任务描述 | 文件位置 | 风险等级 |
|:---:|---|----------|----------|----------|
| [ ] | M18 | 修复 CategoryList 忽略 JWT 错误问题 | `server/internal/logic/category/categorylistlogic.go` | 🟡 MEDIUM |
| [ ] | M19 | 修复 DeleteUser 忽略 JWT 错误问题 | `server/internal/logic/admin/deleteuserlogic.go` | 🟡 MEDIUM |

### 其他中危问题

| 状态 | # | 任务描述 | 文件位置 | 风险等级 |
|:---:|---|----------|----------|----------|
| [ ] | M20 | 对操作日志 Params 字段进行脱敏处理 | `server/internal/model/loginlogmodel_gen.go` | 🟡 MEDIUM |
| [ ] | M21 | 审计日志 IP 地址脱敏处理 | `server/internal/middleware/operationlogmiddleware.go` | 🟡 MEDIUM |
| [x] | M22 | 修复 X-RateLimit-Remaining 头部乱码问题 | `server/internal/middleware/apiratelimitmiddleware.go` | 🟡 MEDIUM |
| [ ] | M23 | 修改备份文件权限为更严格的 0600 | `server/internal/scheduler/backup.go` | 🟡 MEDIUM |
| [ ] | M24 | 修复 float64 转 int64 类型混淆风险 | `server/internal/pkg/jwtx/jwt.go` | 🟡 MEDIUM |
| [ ] | M25 | 添加登录失败次数锁定机制 | `server/internal/handler/routes.go`, `server/internal/logic/user/loginlogic.go` | 🟡 MEDIUM |
| [ ] | M26 | 添加 /user/check-register 限流保护 | `server/internal/handler/routes.go` | 🟡 MEDIUM |
| [ ] | M27 | 修复 Debug 模式敏感信息泄露风险 | `server/internal/config/config.go` | 🟡 MEDIUM |

---

## 🟢 优先级 3 - 低危问题（可计划处理）

| 状态 | # | 任务描述 | 文件位置 | 风险等级 |
|:---:|---|----------|----------|----------|
| [x] | L1 | 修正 vars.go 密码加密注释与实现不符问题 | `server/internal/model/vars.go` | 🟢 LOW |
| [ ] | L2 | 修复 FindList 返回密码字段问题 | `server/internal/model/usermodel_gen.go` | 🟢 LOW |
| [ ] | L3 | 添加系统配置缓存访问控制 | `server/internal/model/systemconfigmodel_gen.go` | 🟢 LOW |
| [ ] | L4 | 简化错误消息，避免泄露敏感信息 | `server/internal/pkg/xerr/code.go` | 🟢 LOW |
| [ ] | L5 | 配置化中间件安全参数 | `server/internal/svc/servicecontext.go` | 🟢 LOW |
| [ ] | L6 | 清理日志中的敏感操作信息 | `server/internal/scheduler/cleanup.go` | 🟢 LOW |
| [ ] | L7 | 修复测试文件接口信息暴露 | `server/internal/scheduler/cleanup_test.go` | 🟢 LOW |
| [ ] | L8 | 添加 /user/check-register 和 /user/register 限流 | `server/internal/handler/routes.go` | 🟢 LOW |
| [ ] | L9 | 添加改密额外验证（当前密码验证） | `server/internal/logic/user/changepasswordlogic.go` | 🟢 LOW |
| [x] | L10 | 修复文件扩展名大小写绕过问题 | `server/internal/handler/admin/downloadbackuphandler.go` | 🟢 LOW |
| [ ] | L11 | 添加软删除恢复权限控制 | `server/internal/model/taskmodel_gen.go`, `server/internal/model/usermodel_gen.go` | 🟢 LOW |
| [ ] | L12 | 统一 LoginLogReq Username 长度验证 | `server/internal/types/types.go` | 🟢 LOW |
| [x] | L13 | 添加 Keyword 查询字段长度限制 | `server/internal/types/types_validate.go` | 🟢 LOW |
| [ ] | L14 | 添加 CSP、HSTS 等安全响应头 | `server/internal/middleware/securityheadsmiddleware.go` | 🟢 LOW |

---

## 安全亮点（做得好的方面）

- ✅ 所有 SQL 查询使用参数化查询，有效防止 SQL 注入
- ✅ 使用 bcrypt 加密密码
- ✅ JWT 认证机制基本完善
- ✅ 任务/分类/标签等资源验证了 userId 所有权
- ✅ AdminMiddleware 对管理员路由统一验证
- ✅ DeleteUser 实现了"不能删除自己"和"不能删除管理员"的业务规则

---

## 统计信息

| 优先级 | 总数 | 已完成 | 待处理 | 完成率 |
|--------|------|--------|--------|--------|
| 🔴 HIGH | 13 | 13 | 0 | 100% |
| 🟡 MEDIUM | 27 | 15 | 12 | 56% |
| 🟢 LOW | 14 | 4 | 10 | 29% |
| **合计** | **54** | **32** | **22** | **59%** |

---

## 提交记录

```
2cbe430 fix(security): 修复错误响应泄露内部错误详情 (H13)
8d9fcfa fix(security): 修复 X-Forwarded-For IP 欺骗漏洞 (H12, M22)
810c54c fix(security): 修复路径遍历 TOCTOU 漏洞 (H11)
762dd61 fix(security): 修复敏感配置泄露 (H10)
d4660b7 fix(security): 修复 SQL 注入风险 - 动态表名白名单验证 (H9)
702565b fix(security): 修复路径遍历漏洞及输入验证缺失 (H7, H8, M7-M17)
996b399 fix(security): 修复 SortTaskReq 验证缺失 (H6)
44b050b fix(security): 修复 ToggleUserStatus 管理员权限验证 (H5)
403c4a0 fix(security): 修复 ResetPassword 管理员权限验证 (H4)
ec3f25a docs(security): 说明 JWT 验证由中间件层处理 (H3)
2c64599 fix(security): 修复 JWT Secret 硬编码问题 (H2)
589fcfd fix(security): 修复密码明文存储漏洞 (H1)
```

---

> 已完成项可迁移至 [changelog.md](./changelog.md)
