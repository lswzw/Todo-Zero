# Todo-Zero 下一步开发清单

> 基于 v1.0.0 代码审查，按优先级分类。

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

## 技术债备注

| 问题 | 位置 | 说明 |
|------|------|------|
| 数据库连接池未配置 | `db/init.go` | `sql.Open` 后未设 `SetMaxOpenConns` 等 |
| 配置热更新缺失 | admin config | 修改系统配置后需实时读取数据库 |
| 环境变量管理 | `request.ts` | API 地址硬编码 |
| App.vue 几乎为空 | `App.vue` (41B) | 缺全局错误处理和布局 |
| BatchDelete SQL 拼接 | model 层 | IN 子句用循环拼接占位符，需注意安全性 |
