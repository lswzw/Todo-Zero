# Todo-Zero 下一步开发清单

> 已完成项已迁移至 [changelog.md](./changelog.md)，本文档仅记录未完成项。

---

## v2.4.0 — 代码质量与安全加固

### 已完成

- [x] **README 文档重构** — 重新编写 README.md，添加 emoji 图标、完整的 API 接口列表、页面路由、技术栈说明等
- [x] **清理无效目录** — 删除 web/src/assets、web/src/components、web/src/utils 三个空目录
- [x] **测试修复** — 修复 `test_helper_test.go` 中 `tasks` 表缺失 `sort_order` 字段问题，所有 34 个单元测试通过
- [x] **ESLint 配置完善** — 添加浏览器全局变量（window、document、localStorage、Blob、setTimeout），消除 9 个 no-undef 错误
- [x] **前端 localStorage XSS 风险修复** — 创建 `useStorage` 组合式函数，替换 `stores/user.ts` 中直接的 localStorage 操作
- [x] **API 速率限制中间件** — 新增 `APIRateLimitMiddleware`，限制每个 IP 每分钟最多 100 次 API 请求，添加到 category、stat、task、user 等路由组
- [x] **任务提醒功能** — 使用浏览器 Notification API 实现任务提醒，支持权限请求、定时检查、点击跳转等功能
- [x] **标签管理功能** — 实现完整的标签 CRUD 功能，包括数据库表、模型、Handler、Logic、API 和前端集成

---

## 待开发项

### P3 — 性能优化

- [ ] **前端懒加载** — 路由级别代码分割，按需加载组件
- [ ] **数据库查询优化** — 添加更多索引，优化复杂查询

### P4 — 开发者体验

- [ ] **添加 Makefile** — 统一构建、测试、部署命令
- [ ] **添加 git hooks** — pre-commit 自动 lint 和格式化

---

> 已完成项已迁移至 [changelog.md](./changelog.md)
