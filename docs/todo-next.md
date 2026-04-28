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

### 后端（已完成，详见 changelog.md）

- ✅ **后端错误码国际化** — `xerr` 错误消息支持 `Accept-Language` 请求头，返回对应语言错误描述（中文/英文）
- ✅ **导出 CSV 中文列名国际化** — `exporttasklogic.go` 中 CSV 表头根据语言参数动态切换

---

## P4-B — 性能优化 （v2.3.0）

> 并发请求优化、Vite 构建分包、批量删除并行化、搜索防抖预留、虚拟滚动预留已完成。

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
