# Todo App 前端开发文档

> 技术栈：Vue 3 + TypeScript + Vite + Element Plus
> 版本：v1.3.0

---

## 一、项目概述

Todo App 前端是个人待办事项管理平台的 Web 端，提供用户注册登录、任务管理、管理员后台等功能。开发期前后端分离，最终目标是将构建产物嵌入 Go 二进制实现单体部署。

---

## 二、技术选型

| 维度 | 选择 | 版本 | 说明 |
|------|------|------|------|
| 框架 | Vue 3 | 3.5.x | Composition API + `<script setup>` |
| 语言 | TypeScript | 5.6.x | 类型安全 |
| 构建工具 | Vite | 5.4.x | 快速开发 + HMR |
| UI 组件库 | Element Plus | 2.13.x | 企业级组件库，中文支持 |
| 图标 | @element-plus/icons-vue | 2.3.x | Element Plus 官方图标 |
| 路由 | Vue Router | 4.6.x | SPA 路由管理 |
| 状态管理 | Pinia | 3.0.x | 轻量级，TS 友好 |
| HTTP 请求 | Axios | 1.15.x | 拦截器 + 请求封装 |
| CSS 预处理 | 原生 CSS | — | scoped 样式，全局变量 |

---

## 三、项目结构

```
web/
├── index.html                 # 入口 HTML
├── package.json               # 依赖管理
├── vite.config.ts             # Vite 配置（别名 + 代理）
├── tsconfig.json              # TypeScript 配置
├── tsconfig.app.json          # 应用 TS 配置（含路径别名）
├── tsconfig.node.json         # Node 端 TS 配置
├── public/                    # 静态资源
└── src/
    ├── main.ts                # 应用入口（注册插件）
    ├── App.vue                # 根组件
    ├── style.css              # 全局样式（CSS 变量）
    ├── vite-env.d.ts          # Vite 类型声明
    ├── api/                   # API 接口封装
    │   ├── request.ts         # Axios 实例（拦截器 + 错误处理）
    │   └── index.ts           # 全部 23 个接口函数
    ├── router/                # 路由
    │   └── index.ts           # 路由定义 + 导航守卫
    ├── stores/                # 状态管理
    │   └── user.ts            # 用户状态（token/登录态持久化）
    └── views/                 # 页面组件
        ├── login.vue          # 登录页
        ├── register.vue       # 注册页
        ├── home.vue           # 任务主页
        └── admin/             # 管理后台
            ├── layout.vue     # 后台布局（顶栏 + 侧边栏）
            ├── user.vue       # 用户管理
            ├── config.vue     # 系统设置
            ├── log.vue        # 操作日志
            └── login-log.vue  # 登录日志
```

---

## 四、路由设计

| 路径 | 页面 | 认证 | 说明 |
|------|------|------|------|
| `/login` | 登录页 | 无 | 紫色渐变背景，居中卡片 |
| `/register` | 注册页 | 无 | 同登录页风格，支持注册关闭状态 |
| `/` | 任务主页 | 需登录 | 统计卡片 + 任务列表 + 筛选/搜索/批量操作 |
| `/admin/user` | 用户管理 | 需登录 + 管理员 | 用户列表/重置密码/禁用/删除 |
| `/admin/config` | 系统设置 | 需登录 + 管理员 | 注册开关/站点名称 |
| `/admin/log` | 操作日志 | 需登录 + 管理员 | 操作日志列表 + 筛选 |
| `/admin/login-log` | 登录日志 | 需登录 + 管理员 | 登录日志列表 + 搜索 |

路由守卫逻辑：
- 未登录访问需认证页面 → 重定向到 `/login`
- 非管理员访问管理员页面 → 重定向到 `/`
- 登录后管理员 → 跳转 `/admin/user`，普通用户 → 跳转 `/`

---

## 五、API 封装

### 5.1 请求实例 (`api/request.ts`)

- baseURL: `/api/v1`（开发环境通过 Vite 代理转发到 `http://localhost:8888`）
- 请求拦截：自动附加 `Authorization` 头（从 localStorage 读取 token）
- 响应拦截：直接返回 `response.data`；错误时显示 ElMessage；401 自动跳转登录

### 5.2 接口列表 (`api/index.ts`)

| 模块 | 函数 | 方法 | 路径 |
|------|------|------|------|
| 用户 | `checkRegister()` | GET | `/user/check-register` |
| | `login(data)` | POST | `/user/login` |
| | `register(data)` | POST | `/user/register` |
| | `getUserInfo()` | GET | `/user/info` |
| | `changePassword(data)` | PUT | `/user/password` |
| 任务 | `getTaskList(params)` | GET | `/task` |
| | `getTaskDetail(id)` | GET | `/task/:id` |
| | `createTask(data)` | POST | `/task` |
| | `updateTask(id, data)` | PUT | `/task/:id` |
| | `toggleTask(id)` | PATCH | `/task/:id/toggle` |
| | `deleteTask(id)` | DELETE | `/task/:id` |
| | `batchTask(data)` | POST | `/task/batch` |
| 分类 | `getCategoryList()` | GET | `/category` |
| | `createCategory(data)` | POST | `/category` |
| 统计 | `getStat()` | GET | `/stat` |
| 管理员 | `getUserList(params)` | GET | `/admin/user` |
| | `resetPassword(id, data)` | PUT | `/admin/user/:id/password` |
| | `toggleUserStatus(id)` | PATCH | `/admin/user/:id/toggle` |
| | `deleteUser(id)` | DELETE | `/admin/user/:id` |
| | `getConfigList()` | GET | `/admin/config` |
| | `updateConfig(data)` | PUT | `/admin/config` |
| | `getOperationLogList(params)` | GET | `/admin/log/operation` |
| | `getLoginLogList(params)` | GET | `/admin/log/login` |

---

## 六、状态管理

### User Store (`stores/user.ts`)

| 属性 | 类型 | 说明 | 持久化 |
|------|------|------|--------|
| `token` | string | JWT 令牌 | localStorage |
| `userId` | number | 用户 ID | localStorage |
| `username` | string | 用户名 | localStorage |
| `isAdmin` | boolean | 是否管理员 | localStorage |

| 方法 | 说明 |
|------|------|
| `setLogin(data, name)` | 登录成功后保存 token + 用户信息 |
| `setUserInfo(data)` | 更新用户信息（来自 /user/info 接口） |
| `logout()` | 清除所有状态 + localStorage |

---

## 七、页面说明

### 7.1 登录页 (`login.vue`)

- 全屏紫色渐变背景，居中白色卡片
- 表单：用户名 + 密码 → 调用 `login` 接口
- 登录成功：管理员 → `/admin/user`，普通用户 → `/`
- 底部注册链接：受 `checkRegister` 接口控制是否显示

### 7.2 注册页 (`register.vue`)

- 同登录页风格，增加确认密码字段
- 表单验证：用户名 3-20 位，密码 6-20 位，两次密码一致
- 注册关闭时显示锁定提示

### 7.3 任务主页 (`home.vue`)

- 顶部导航栏：Logo + 用户名 + 修改密码 + 退出登录
- 统计卡片行：总任务 / 待办 / 已完成 / 完成率
- 任务列表区：
  - 筛选：状态 / 分类 / 优先级 下拉框
  - 搜索：关键词输入框
  - 多选模式：批量完成 / 取消 / 删除
  - 任务项：状态圆圈(点击切换) + 标题/内容 + 优先级/分类标签 + 编辑/删除
  - 分页器
- 弹窗：新增/编辑任务、修改密码

### 7.4 管理后台 (`admin/`)

统一布局：60px 顶栏 + 200px 侧边栏 + 弹性内容区

| 页面 | 功能 |
|------|------|
| 用户管理 | 表格 + 搜索 + 重置密码/禁用/启用/删除 |
| 系统设置 | 注册开关(Switch) + 站点名称(输入框+保存) |
| 操作日志 | 表格 + 用户名搜索 + 操作类型筛选 + 彩色标签 |
| 登录日志 | 表格 + 用户名搜索 + 成功/失败标签 |

---

## 八、Vite 配置说明

```typescript
// vite.config.ts
{
  resolve: {
    alias: { '@': 'src/' }     // @ 指向 src 目录
  },
  server: {
    host: '0.0.0.0',           // 允许外部访问
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8888',  // 代理到后端
        changeOrigin: true,
      }
    }
  }
}
```

---

## 九、开发与构建

```bash
# 安装依赖
cd web && npm install

# 开发模式（HMR + API 代理）
npm run dev

# 类型检查
npx vue-tsc --noEmit

# 生产构建（输出到 web/dist/）
npm run build

# 预览构建产物
npm run preview
```

---

## 十、设计规范

### 配色

| 用途 | 颜色 | 色值 |
|------|------|------|
| 主色 | 紫蓝 | `#667eea` → `#764ba2`（渐变） |
| 成功 | 绿色 | `#67c23a` |
| 警告 | 橙色 | `#e6a23c` |
| 危险 | 红色 | `#f56c6c` |
| 信息 | 蓝色 | `#409eff` |
| 文字主色 | 深灰 | `#303133` |
| 文字次色 | 灰色 | `#909399` |
| 页面背景 | 浅灰 | `#f5f7fa` |

### 组件风格

- 卡片：白色背景、12px 圆角、轻阴影 `0 2px 12px rgba(0,0,0,0.04)`
- 登录/注册卡片：16px 圆角、大阴影
- 按钮：primary 使用紫蓝渐变
- 表格：斑马纹

### 交互模式

- Dialog 弹窗：表单编辑（新增任务、修改密码、重置密码）
- Popconfirm：删除确认
- ElMessage：操作成功/失败提示
- ElMessageBox：退出登录确认
