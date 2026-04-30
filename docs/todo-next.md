# Todo-Zero 开发计划

> 安全评审（2026-04-29）全部 54 项已修复并归档至 [changelog.md](./changelog.md) v2.5.0

---

## v3.0.0 - 界面现代化视觉重构

> 里程碑：毛玻璃效果 + 现代渐变 + 精致动画，全面提升视觉体验
> 详细规划参考：[todo-next-web.md](./todo-next-web.md)

---

### Phase 1: 基础架构准备 v3.0.1

- [x] **T1.1** 🔴 高 | 创建全局样式变量文件 `styles/variables.css` ✅ 2026-04-30
  - 定义统一色彩系统：`--primary`, `--primary-dark`, `--primary-light`
  - 定义背景变量：`--bg-primary`, `--bg-glass`, `--bg-glass-dark`
  - 定义阴影变量：`--shadow-glow`, `--shadow-card`
  - 定义间距、圆角、过渡时间等基础变量

- [x] **T1.2** 🔴 高 | 创建毛玻璃组件 `components/GlassCard.vue` ✅ 2026-04-30
  - 封装可复用的玻璃态卡片组件
  - 支持 `dark` / `hoverable` props
  - 实现 `backdrop-filter: blur(20px)` + 半透明背景
  - 悬停效果：`translateY(-2px)` + 阴影增强
  - 包含 `-webkit-backdrop-filter` 兼容前缀

- [x] **T1.3** 🟡 中 | 创建渐变背景组件 `components/GradientBg.vue` ✅ 2026-04-30
  - 支持 `primary` / `secondary` / `dark` 三种变体
  - 主渐变：`linear-gradient(135deg, #667eea, #764ba2)`
  - 次渐变：`linear-gradient(135deg, #f5f7fa, #c3cfe2)`
  - 暗色渐变：`linear-gradient(135deg, #1a1a2e, #16213e)`
  - 添加 `radial-gradient` 装饰层 + 浮动动画

- [x] **T1.4** 🔴 高 | 更新主样式文件 `styles/main.css` ✅ 2026-04-30
  - 引入全局变量文件
  - 重置基础样式（body 背景、字体等）
  - 整合全局过渡/动画基础类

---

### Phase 2: 公共组件重构 v3.0.2

- [x] **T2.1** 🔴 高 | 导航栏毛玻璃化 ✅ 2026-04-30
  - 顶部导航栏添加 `backdrop-filter: blur(20px)`
  - 半透明白色背景 `rgba(255, 255, 255, 0.7)`
  - 底部边框使用 `rgba(255, 255, 255, 0.3)`
  - 添加 `position: sticky` 固定效果

- [x] **T2.2** 🔴 高 | 按钮样式升级 ✅ 2026-04-30
  - 主按钮：渐变背景 `linear-gradient(135deg, #667eea, #764ba2)`
  - 添加 `hover` 缩放 + 亮度增强动画
  - 添加 `active` 按下反馈
  - 次要按钮：描边 + 渐变文字

- [x] **T2.3** 🟡 中 | 输入框样式升级 ✅ 2026-04-30
  - 聚焦时添加光晕效果 `box-shadow` 渐变辉光
  - 输入框背景微透明 `rgba(255, 255, 255, 0.5)`
  - 聚焦边框渐变色过渡
  - placeholder 颜色统一

- [x] **T2.4** 🔴 高 | 卡片组件升级 ✅ 2026-04-30
  - 统一使用 `GlassCard` 组件
  - 添加阴影层级系统（sm / md / lg）
  - 圆角统一 `16px`
  - 悬停阴影增强 + 微上移

---

### Phase 3: 页面重构 v3.0.3

- [x] **T3.1** 🔴 高 | 登录页重构 `views/Login.vue` ✅ 2026-04-30
  - 使用 `GradientBg` 动态渐变背景（primary 变体）
  - 登录卡片使用 `GlassCard` 毛玻璃效果
  - 添加背景浮动装饰元素（气泡/圆形光斑）
  - 输入框聚焦光晕动画
  - 登录按钮渐变样式
  - 表单标题添加渐变文字效果

- [x] **T3.2** 🔴 高 | 注册页重构 `views/Register.vue` ✅ 2026-04-30
  - 与登录页风格统一
  - 同样使用 `GradientBg` + `GlassCard`
  - 表单字段与登录页一致的输入框样式
  - 添加注册成功动画反馈

- [x] **T3.3** 🔴 高 | 主页重构 `views/Home.vue` ✅ 2026-04-30
  - 导航栏毛玻璃化（T2.1 成果）
  - 统计卡片使用 `GlassCard` + 渐变色图标
  - 统计卡片每个指标使用不同渐变色区分
  - 任务卡片列表使用 `GlassCard` + 悬停效果
  - 任务状态切换视觉增强
  - 整体页面背景使用 `GradientBg` secondary 变体

- [x] **T3.4** 🟡 中 | 任务详情页重构 `views/TaskDetail.vue` ✅ 2026-04-30
  - 页面背景使用 `GradientBg` secondary 变体
  - 详情区域使用 `GlassCard` 玻璃态布局
  - 子任务列表卡片化
  - 评论/操作区域毛玻璃化
  - 添加页面进入过渡动画

- [x] **T3.5** 🟡 中 | 回收站重构 `views/Trash.vue` ✅ 2026-04-30
  - 页面背景使用 `GradientBg` secondary 变体
  - 任务卡片使用 `GlassCard` 玻璃态
  - 恢复/删除按钮渐变样式
  - 空状态插图优化
  - 列表项悬停效果

- [x] **T3.6** 🟡 中 | 管理员面板重构 `views/Admin.vue` ✅ 2026-04-30
  - 侧边栏毛玻璃化 `backdrop-filter: blur(20px)`
  - 侧边栏半透明背景
  - 内容区域使用 `GlassCard` 包裹
  - 数据统计卡片渐变图标
  - 表格行悬停效果增强

---

### Phase 4: 交互动画 v3.0.4

- [x] **T4.1** 🟡 中 | 添加页面过渡动画 ✅ 2026-04-30
  - 路由切换 `<Transition>` 淡入淡出 `fade`
  - 可选 `slide-up` / `slide-right` 切换效果
  - 动画时长 `0.3s ease`
  - 避免路由切换闪烁

- [x] **T4.2** 🟡 中 | 添加卡片悬停效果 ✅ 2026-04-30
  - `hover` 时 `transform: translateY(-4px) scale(1.01)`
  - 阴影从 `--shadow-card` 过渡到 `--shadow-glow`
  - 过渡时长 `0.3s ease`
  - 使用 `will-change: transform` 性能提示

- [x] **T4.3** 🟢 低 | 添加按钮点击涟漪效果 ✅ 2026-04-30
  - 点击位置产生涟漪扩散动画
  - 涟漪颜色 `rgba(255, 255, 255, 0.3)`
  - 动画时长 `0.6s`
  - 封装为 `Ripple` 指令或组件

- [x] **T4.4** 🟡 中 | 添加任务状态切换动画 ✅ 2026-04-30
  - 勾选时 checkbox 缩放弹跳动画
  - 文字划线过渡动画
  - 完成状态卡片透明度渐变
  - 使用 `transform` + `opacity` 确保性能

---

### Phase 5: 性能优化与兼容 v3.0.5

- [x] **T5.1** 🔴 高 | backdrop-filter 降级处理 ✅ 2026-04-30
  - 检测 `backdrop-filter` 支持情况
  - 不支持时回退到 `rgba(255, 255, 255, 0.95)` 实色背景
  - Safari 旧版本添加 `-webkit-backdrop-filter`
  - iOS Safari 特殊处理

- [x] **T5.2** 🟡 中 | CSS 动画性能优化 ✅ 2026-04-30
  - 所有动画仅使用 `transform` / `opacity`
  - 添加 `transform: translateZ(0)` 触发 GPU 加速
  - 合理使用 `will-change` 属性
  - 复杂动画使用 `requestAnimationFrame`

- [x] **T5.3** 🟡 中 | 懒加载优化 ✅ 2026-04-30
  - 路由组件 `defineAsyncComponent` 懒加载
  - 图片使用 `loading="lazy"`
  - 非首屏组件延迟渲染

- [x] **T5.4** 🔴 高 | 首屏加载优化 ✅ 2026-04-30
  - 关键 CSS 内联到 HTML
  - 预加载关键字体
  - 骨架屏替换 loading
  - LCP < 2.5s, FID < 100ms 目标

---

## 进度总览

| 阶段 | 任务数 | 🔴高 | 🟡中 | 🟢低 | 进度 |
|------|--------|------|------|------|------|
| Phase 1 - 基础架构 | 4 | 3 | 1 | 0 | 4/4 ✅ |
| Phase 2 - 公共组件 | 4 | 3 | 1 | 0 | 4/4 ✅ |
| Phase 3 - 页面重构 | 6 | 3 | 3 | 0 | 6/6 ✅ |
| Phase 4 - 交互动画 | 4 | 0 | 3 | 1 | 4/4 ✅ |
| Phase 5 - 性能优化 | 4 | 2 | 2 | 0 | 4/4 ✅ |
| **合计** | **22** | **11** | **10** | **1** | **22/22 ✅** |

---

## 执行顺序（按优先级排序）

1. **T1.1** → **T1.4** → **T1.2** → **T1.3** — 基础架构先就位
2. **T2.1** → **T2.2** → **T2.4** → **T2.3** — 公共组件先行
3. **T3.1** → **T3.2** → **T3.3** → **T3.4** → **T3.5** → **T3.6** — 页面逐个重构
4. **T5.1** → **T5.4** — 关键性能保障（与页面重构穿插进行）
5. **T4.1** → **T4.2** → **T4.4** → **T4.3** — 交互动画收尾
6. **T5.2** → **T5.3** — 剩余性能优化

---

## 完成标准 v3.0.0

- [x] 所有页面采用毛玻璃设计风格
- [x] 保持原有功能完整性（无功能回归）
- [x] 性能指标达标：LCP < 2.5s, FID < 100ms
- [x] 响应式适配：移动端 / 平板 / 桌面端
- [x] 浏览器兼容：Chrome / Firefox / Safari / Edge
- [x] 代码通过 ESLint 检查
- [x] 更新 [changelog.md](./changelog.md) v3.0.0 版本记录 ✅ 2026-04-30

---

## v3.1.0 — 国际化增强 + Bug 修复 + 代码清理

### 国际化增强

- **useLocale 全面升级** — `composables/useLocale.ts` 新增全局 singleton `syncedLocale`、跨 tab `storage` 事件同步、`setLocale`/`toggleLocale`/`getLocaleLabel`/`getLocaleAlias` 方法、`currentLocaleOption`/`currentLocaleLabel`/`currentLocaleAlias` 计算属性
- **语言切换统一为悬浮按钮** — 所有页面统一使用右下角圆形悬浮按钮 `LocaleSwitch.vue`，移除各页面独立的 `el-select` 下拉切换
- **新增翻译键** — 添加 `common.switchLanguage` 等翻译键
- **Element Plus 语言联动** — 主应用使用 `ElConfigProvider` + `:locale` 动态绑定，Element Plus 组件语言随 vue-i18n 同步切换

### Bug 修复

- **登录页错误提示不显示** — 根因：`request.ts`（`.ts` 文件）中 `ElMessage` 的 CSS 未被 `unplugin-vue-components` 处理；修复：在 `main.ts` 显式导入 `element-plus/es/components/message/style/css` 和 `message-box/style/css`
- **管理员退出登录弹窗不一致** — 统一为 home.vue 样式（无图标、自定义按钮文字、禁止点击遮罩关闭、禁止 ESC 关闭）
- **管理员退出登录按钮多余背景色** — 从 `text type="danger"` 改为 `link` 类型

### 代码清理

- **清理无效文件** — 删除 `web/src/assets`、`web/src/components`（旧 `GlassCard.vue`）、`web/src/utils` 等空目录/废弃文件
- **清理未使用的 API 和类型** — 删除 `api/index.ts` 中未使用的 Tag 相关 API 函数，删除 `types/index.ts` 中未使用的 `TagItem`/`TagListResp` 类型
- **修复类型警告** — 移除 `LocaleSwitch.vue` 中未使用的 `LocaleValue` 导入，移除未使用的 `props` 赋值
