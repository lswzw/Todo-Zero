# Todo-Zero 开发清单

> 基于 v1.1.0 全面代码审查，详见 `docs/code-review.md`。

---

## HIGH — 必须修复

- [ ] **#1 管理员路由无 RBAC 中间件** — 新增 admin logic 忘记 checkAdmin 即可越权
- [ ] **#2 JWT Secret 硬编码默认值** — 可伪造任意用户 token，首次启动应自动生成随机密钥
- [ ] **#3 操作日志系统形同虚设** — 只有注册写了一条记录，所有管理操作均无日志

---

## MEDIUM — 应该修复

- [ ] **#4 批量操作静默吞错** — BatchTask 错误被 `_ =` 丢弃，前端收到"成功"但可能部分失败
- [ ] **#5 ToggleUserStatus TOCTOU 竞态** — 先查再改，并发可覆盖状态；应使用 SQL 原子操作
- [ ] **#6 ToggleTask TOCTOU 竞态** — 同 #5，应使用已有的 UpdateStatus
- [ ] **#7 UpdateTask 无法清空可选字段** — `if req.X != ""` 导致无法清空回零值
- [ ] **#8 操作日志过滤参数未生效** — FindList 只接受 page/pageSize，action/username 被忽略
- [ ] **#9 OperationLogItem 映射字段丢失** — UserId 硬编码 0，TargetType 为空
- [ ] **#10 LoginReq 缺少验证标签** — 空/超长请求可导致 bcrypt DoS
- [ ] **#11 SQL 注释解析逻辑脆弱** — `--` 在字符串值中会被误剥离
- [ ] **#12 数据库连接池未配置** — 未设 SetMaxOpenConns，长时间运行可能连接泄漏
- [ ] **#13 DeleteUser 未阻止删除其他管理员** — API 层面无保护
- [ ] **#14 静态文件缺少安全响应头** — 无 nosniff/X-Frame-Options，存在 MIME 嗅探风险
- [ ] **#15 Token 存 localStorage 有 XSS 风险** — 任何 XSS 可窃取 token
- [ ] **#16 前端 isAdmin 可被绕过** — localStorage 可 DevTools 篡改，路由守卫应 API 验证
- [ ] **#17 大量 any 类型** — TypeScript 形同虚设，应定义接口类型
- [ ] **#18 空 catch 块吞掉所有错误** — 页面加载失败无提示
- [ ] **#19 flag 与配置文件优先级不清** — 指定配置文件时 flags 完全不生效
- [ ] **#20 API 响应拦截器未统一解包** — 调用方需 `.data` 但用 `as any` 绕过

---

## LOW — 建议改进

- [ ] **#21** Register TOCTOU 竞态 — 先查再插，直接 Insert 捕获 UNIQUE 更优
- [ ] **#22** checkAdmin 重复 7 次 — 违反 DRY，应提取公共函数
- [ ] **#23** 分类删除策略不一致 — categories 物理删除 vs tasks/users 软删除
- [ ] **#24** ToggleUserStatus 应用 UpdateStatus 而非通用 Update
- [ ] **#25** 缺少 defer sqliteDB.Close() — WAL 文件可能未清理
- [ ] **#26** init.sql 中 PRAGMA 与 InitDB 重复
- [ ] **#27** 嵌入配置中 AccessSecret 明文
- [ ] **#29** 登录/注册成功后硬编码延迟跳转 — 不必要
- [ ] **#30** 管理员登录后强制跳转管理页面 — 应可配置
- [ ] **#31** 筛选条件改变时未重置页码
- [ ] **#33** config.vue switchValue 单一 ref 无法多配置
- [ ] **#34** 密码确认验证器类型不安全
- [ ] **#35** 无障碍性缺失 — 缺 ARIA/语义化标签/键盘导航
- [ ] **#36** 默认管理员密码 admin123 — 注释明文，应首次启动强制修改
- [ ] **#38** StatLogic 应使用 CountStats 而非 FindList 全量加载
- [ ] **#39** Model 层多个方法已定义但未使用（BatchDelete、UpdateStatus 等）

---

## P0 已完成 (v1.1.0)

- [x] 前后端状态/优先级映射不一致
- [x] 任务列表 N+1 查询
- [x] 数据库缺少索引
- [x] 用户列表返回密码哈希
