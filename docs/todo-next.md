# Todo-Zero 下一步开发清单

> 已完成项已迁移至 [changelog.md](./changelog.md)，本文档仅记录未完成项。

---

## P1 ~ P3-B — 已全部完成

详见 [changelog.md](./changelog.md) v1.1.0 ~ v2.1.0。

---

## 技术债 — 已全部清理

详见 [changelog.md](./changelog.md) v2.0.0 技术债清理。

---

## P4-A — 国际化 (i18n) ✅ 已完成 (v2.2.0)

详见 [changelog.md](./changelog.md) v2.2.0。

### 后端（可选，未实施）

- [ ] **后端错误码国际化** — `xerr` 错误消息支持 `Accept-Language` 请求头，返回对应语言错误描述（当前仅中文）
- [ ] **导出 CSV 中文列名国际化** — `exporttasklogic.go` 中 CSV 表头硬编码中文列名，需根据语言参数切换

---

## P4-B — 性能优化 （v2.3.0）

> 当前状态：路由懒加载已实现，但存在全量导入、无缓存、无分包等问题。

### P0 — 高优先级

- [ ] **API 缓存（Pinia Store）** — 对低频变化数据引入缓存，避免跨页面重复请求
  - `getCategoryList()` — `home.vue` 和 `task-detail.vue` 都独立请求，应缓存到 store
  - `getStat()` — 操作后每次重新请求，应缓存 + 操作后失效
  - `getConfigList()` — `config.vue` 和 `backup.vue` 重复请求，应缓存
  - 缓存策略：首次加载存入 store，增删改操作后 `invalidate`，TTL 5 分钟自动过期

- [ ] **Element Plus 按需导入** — 当前全量 `import ElementPlus` + 全量注册图标，改为：
  - 安装 `unplugin-vue-components` + `unplugin-auto-import`
  - `vite.config.ts` 配置 Element Plus resolver
  - 移除 `main.ts` 中 `import ElementPlus` 和 `import * as ElementPlusIconsVue` 全量注册
  - 预期减少 bundle 体积 30%+

### P1 — 中优先级

- [ ] **并发请求优化** — `home.vue` 操作后同时刷新任务列表和统计，改为 `Promise.all([loadTasks(), loadStat()])` 并发执行
  - `handleToggle`、`handleDelete`、`handleBatch`、`handleSubmitTask` 均适用
  - `onMounted` 中 `loadStat()` + `loadCategories()` + `loadTasks()` 同理

- [ ] **Vite 构建分包** — `vite.config.ts` 添加 `rollupOptions.output.manualChunks`
  - `vendor` — vue/vue-router/pinia/axios
  - `element-plus` — Element Plus 组件库
  - `icons` — Element Plus Icons
  - 分离后利于浏览器缓存，首屏仅加载必要 chunk

### P2 — 低优先级

- [ ] **批量删除并行化** — `trash.vue` 中 `handleBatchPermanentDelete` 使用 `for...of` 逐个 `await`，改为 `Promise.allSettled` 并行删除

- [ ] **搜索防抖预留** — 当前搜索仅 `@keyup.enter` + `@clear` 触发，若将来改为实时搜索（`@input`），需引入防抖（`lodash-es/debounce` 或手写）

- [ ] **虚拟滚动预留** — 当前 `pageSize=10` 性能无问题；若未来增大分页或取消分页，需引入 `vue-virtual-scroller` 或 `@tanstack/vue-virtual`
