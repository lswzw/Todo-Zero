# Todo-Zero 下一步开发清单

> 已完成项已迁移至 [changelog.md](./changelog.md)，本文档仅记录未完成项。

---

## P1 ~ P4-A — 已全部完成

详见 [changelog.md](./changelog.md) v1.1.0 ~ v2.2.0。

---

## 技术债 — 已全部清理

详见 [changelog.md](./changelog.md) v2.0.0 技术债清理。

---

## P4-B — 性能优化 （v2.3.0）

### P0 — 高优先级

- [ ] **Element Plus 按需导入** — 当前全量 `import ElementPlus` + 全量注册图标，改为：
  - 安装 `unplugin-vue-components` + `unplugin-auto-import`
  - `vite.config.ts` 配置 Element Plus resolver
  - 移除 `main.ts` 中 `import ElementPlus` 和 `import * as ElementPlusIconsVue` 全量注册
  - 预期减少 bundle 体积 30%+
