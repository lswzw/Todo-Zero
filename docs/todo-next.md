# Todo-Zero 下一步开发清单

> 已完成项已迁移至 [changelog.md](./changelog.md)，本文档仅记录未完成项。

---

## 待开发项

暂无待开发项。所有 v2.4.0 及之前的任务已完成，归档于 [changelog.md](./changelog.md)。

---

## 安全评审报告

> 本次安全评审完成于 2026-04-29，审计范围：server/internal 下所有后端 Go 代码（120+ 文件）
>
> **汇总：HIGH 级问题 13 个，MEDIUM 级问题 35 个，LOW 级问题 18 个**

---

## 🔴 HIGH 级安全问题（需立即处理）

### 1. 密码明文存储（严重）
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/model/usermodel_gen.go` |
| **问题** | `Password` 字段以明文存储，注释声称 bcrypt 加密但实际未实现 |
| **风险** | 数据库泄露后所有用户密码完全暴露 |
| **修复** | 在 `Insert`/`UpdatePassword` 中使用 `bcrypt.GenerateFromPassword`，登录时用 `bcrypt.CompareHashAndPassword` 验证 |

### 2. JWT Secret 硬编码
| 项目 | 内容 |
|------|------|
| **文件** | `server/etc/todo-api.yaml` (L8), `server/internal/config/config.go` (L10) |
| **问题** | JWT Secret 硬编码为 `"todo-app-jwt-secret-key-2024"` |
| **风险** | 攻击者获取密钥后可伪造任意用户 JWT Token |
| **修复** | 使用环境变量或密钥管理服务存储 JWT Secret |

### 3. Token 缺少过期时间验证
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/pkg/jwtx/jwt.go` |
| **问题** | `GetUserIdFromCtx` 和 `GetIsAdminFromCtx` 未验证 Token 是否过期 |
| **风险** | 过期 Token 可能仍被使用 |
| **修复** | 添加 Token 过期验证逻辑 |

### 4. 管理员操作缺少权限验证 - ResetPassword
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/logic/admin/resetpasswordlogic.go` (L28-L38) |
| **问题** | 任何已认证用户可重置任意用户密码，未验证管理员身份 |
| **风险** | 普通用户可修改管理员密码 |
| **修复** | 添加 `jwtx.GetIsAdminFromCtx` 管理员验证，不能重置自己/其他管理员密码 |

### 5. 管理员操作缺少权限验证 - ToggleUserStatus
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/logic/admin/toggleuserstatuslogic.go` (L27-L42) |
| **问题** | 任何已认证用户可切换任意用户状态，未验证管理员身份 |
| **风险** | 普通用户可禁用管理员账户 |
| **修复** | 添加管理员验证，不能操作自己/其他管理员 |

### 6. SortTaskReq 完全缺失验证
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/types/types.go`, `server/internal/types/types_validate.go` |
| **问题** | `SortTaskReq` 及其 `SortOrderItem.Id`/`SortOrder` 无任何验证 |
| **风险** | 负数索引可能导致数据越界或异常行为 |
| **修复** | 实现完整 Validate() 方法，验证 Id > 0，SortOrder >= 0 |

### 7. DownloadBackupReq 文件名无路径遍历防护
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/types/types.go`, `server/internal/types/types_validate.go` |
| **问题** | `FileName` 字段无路径遍历过滤，攻击者可用 `../../../etc/passwd` 下载任意文件 |
| **风险** | 敏感文件泄露 |
| **修复** | 验证 FileName 只含合法字符合，不含 `..` 或绝对路径 |

### 8. RestoreBackupReq 文件名无路径遍历防护
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/types/types.go`, `server/internal/types/types_validate.go` |
| **问题** | 同 DownloadBackupReq，可能导致任意文件覆盖 |
| **风险** | 系统文件被恶意覆盖 |
| **修复** | 同上，验证 FileName 在合法备份目录范围内 |

### 9. SQL 注入风险 - 动态表名拼接
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/scheduler/backup.go` (L235, L243) |
| **问题** | 使用 `fmt.Sprintf("DELETE FROM %s", table)` 拼接 SQL 表名 |
| **风险** | 虽然当前 table 来源硬编码，但存在潜在注入风险 |
| **修复** | 使用白名单验证 table 值 |

### 10. 敏感配置直接返回客户端
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/logic/admin/configlistlogic.go` (L27-L43) |
| **问题** | 返回所有系统配置项未过滤敏感配置 |
| **风险** | 备份配置、安全配置等敏感信息泄露 |
| **修复** | 过滤或脱敏敏感配置项后再返回 |

### 11. 路径遍历 TOCTOU 漏洞
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/scheduler/backup.go` (L179-L184) |
| **问题** | 先计算绝对路径再检查前缀，os.Stat 和后续操作间存在时间窗口 |
| **风险** | TOCTOU 攻击可能在检查后替换文件 |
| **修复** | 打开文件后使用 O_NOFOLLOW 标志或使用 os.OpenFile 而非后续操作 |

### 12. X-Forwarded-For IP 欺骗漏洞
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/middleware/apiratelimitmiddleware.go`, `server/internal/middleware/loginratelimitmiddleware.go` |
| **问题** | 直接信任 X-Forwarded-For 和 X-Real-IP 头，未验证可信代理链 |
| **风险** | 攻击者伪造 IP 绕过限流 |
| **修复** | 只信任已知可信代理链的 IP，忽略伪造头 |

### 13. 错误响应泄露内部错误详情
| 项目 | 内容 |
|------|------|
| **文件** | `server/internal/pkg/xerr/response.go` (L38-L49) |
| **问题** | 非 CodeError 类型错误直接返回 err.Error() 给客户端 |
| **风险** | SQL 错误、堆栈信息、文件路径等敏感信息泄露 |
| **修复** | 未知错误返回通用消息，不暴露 err.Error() |

---

## 🟡 MEDIUM 级安全问题（建议尽快处理）

### 密码相关
| # | 文件 | 问题 |
|---|------|------|
| M1 | `logic/user/registerlogic.go` | 注册时未验证密码强度，允许简单密码 |
| M2 | `logic/user/loginlogic.go` | 登录失败日志区分"用户不存在"和"密码错误" |
| M3 | `logic/admin/triggerbackuplogic.go` | 备份失败错误信息泄露内部详情 |
| M4 | `logic/admin/restorebackuplogic.go` | 恢复备份错误信息泄露内部详情 |
| M5 | `types/types.go` - ChangePasswordReq | OldPassword/NewPassword 无长度验证 |
| M6 | `types/types.go` - LoginReq | Password 无长度验证 |

### 输入验证缺失
| # | 文件 | 问题 |
|---|------|------|
| M7 | `types_validate.go` | DeleteTaskReq 缺 Validate()，Id 无正数验证 |
| M8 | `types_validate.go` | RestoreTaskReq 缺 Validate() |
| M9 | `types_validate.go` | PermanentDeleteTaskReq 缺 Validate() |
| M10 | `types_validate.go` | DeleteUserReq 缺 Validate() |
| M11 | `types_validate.go` | TaskDetailReq 缺 Validate() |
| M12 | `types_validate.go` | ToggleTaskReq 缺 Validate() |
| M13 | `types_validate.go` | ToggleUserStatusReq 缺 Validate() |
| M14 | `types_validate.go` | UpdateTagReq 缺 Validate() |
| M15 | `types_validate.go` | DeleteTagReq 缺 Validate() |
| M16 | `types_validate.go` | CreateTagReq 缺 Name/Color 长度验证 |
| M17 | `types/types.go` - ExportTaskReq | Keyword 无长度限制，可能导致数据库性能问题 |

### 权限/认证问题
| # | 文件 | 问题 |
|---|------|------|
| M18 | `logic/category/categorylistlogic.go` (L28) | 忽略 JWT 错误，未登录用户可能访问 |
| M19 | `logic/admin/deleteuserlogic.go` (L30) | 忽略 JWT 错误 |

### 其他
| # | 文件 | 问题 |
|---|------|------|
| M20 | `model/loginlogmodel_gen.go` (L81) | Params 字段记录请求参数可能含敏感信息 |
| M21 | `middleware/operationlogmiddleware.go` | 审计日志记录 IP 地址（含端口）等敏感信息 |
| M22 | `middleware/apiratelimitmiddleware.go` | X-RateLimit-Remaining 头部值大于9时产生乱码 |
| M23 | `scheduler/backup.go` (L49,L101) | 备份文件权限 0755 过于宽松 |

### 类型安全
| # | 文件 | 问题 |
|---|------|------|
| M24 | `pkg/jwtx/jwt.go` (L22-L23) | float64 转 int64 存在类型混淆风险 |

### 中间件/路由
| # | 文件 | 问题 |
|---|------|------|
| M25 | `handler/routes.go` (L21-L26) | /user/login 无暴力破解防护（如失败次数锁定）|
| M26 | `handler/routes.go` | /user/check-register 可能被枚举已注册用户 |
| M27 | `config/config.go` | Debug 模式可能泄露敏感信息 |

---

## 🟢 LOW 级安全问题（可计划处理）

| # | 文件 | 问题描述 |
|---|------|----------|
| L1 | `model/vars.go` (L16) | 注释声明 bcrypt 加密但实际未实现（与 H1 重复） |
| L2 | `model/usermodel_gen.go` (L96) | FindList 返回包含 password 的结果集 |
| L3 | `model/systemconfigmodel_gen.go` (L24-35) | 内存缓存无访问控制 |
| L4 | `pkg/xerr/code.go` | 错误消息过于详细（如"原密码错误"）|
| L5 | `svc/servicecontext.go` (L47-52) | 中间件配置硬编码，无安全参数配置 |
| L6 | `scheduler/cleanup.go` | 日志输出可能含敏感操作信息 |
| L7 | `scheduler/cleanup_test.go` | Mock 实现暴露内部接口 |
| L8 | `handler/routes.go` (L270-288) | /user/check-register 和 /user/register 无限流 |
| L9 | `handler/routes.go` (L308-327) | /user/info 和 /user/password 改密未要求额外验证 |
| L10 | `handler/admin/downloadbackuphandler.go` | 文件扩展名检查仅小写 .bak，可能大小写绕过 |
| L11 | `model/通用` | 软删除记录可恢复，需严格权限控制 |
| L12 | `types/types.go` - LoginLogReq | Username 长度定义与验证不一致 |
| L13 | `types/types.go` - TagListReq/UserListReq/OperationLogReq | Keyword 无长度限制 |
| L14 | `handler/routes.go` | SecurityHeadersMiddleware 缺少 CSP、HSTS 等关键头 |

---

## 安全亮点（做得好的方面）

- ✅ 所有 SQL 查询使用参数化查询，有效防止 SQL 注入
- ✅ 使用 bcrypt 加密密码（虽未在 UserModel 中应用）
- ✅ JWT 认证机制基本完善
- ✅ 任务/分类/标签等资源验证了 userId 所有權
- ✅ AdminMiddleware 对管理员路由统一验证
- ✅ DeleteUser 实现了"不能删除自己"和"不能删除管理员"的业务规则

---

> 已完成项已迁移至 [changelog.md](./changelog.md)
