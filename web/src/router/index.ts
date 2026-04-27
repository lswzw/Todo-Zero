import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login.vue'),
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/register.vue'),
    },
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/home.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/task/:id',
      name: 'TaskDetail',
      component: () => import('@/views/task-detail.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/trash',
      name: 'Trash',
      component: () => import('@/views/trash.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin',
      name: 'Admin',
      component: () => import('@/views/admin/layout.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        { path: '', redirect: '/admin/user' },
        { path: 'user', name: 'AdminUser', component: () => import('@/views/admin/user.vue') },
        { path: 'config', name: 'AdminConfig', component: () => import('@/views/admin/config.vue') },
        { path: 'log', name: 'AdminLog', component: () => import('@/views/admin/log.vue') },
        { path: 'login-log', name: 'AdminLoginLog', component: () => import('@/views/admin/login-log.vue') },
        { path: 'backup', name: 'AdminBackup', component: () => import('@/views/admin/backup.vue') },
      ],
    },
  ],
})

// 是否已从服务端验证过身份
let authVerified = false

router.beforeEach(async (to) => {
  const userStore = useUserStore()

  // 未登录直接跳转登录页
  if (to.meta.requiresAuth && !userStore.token) {
    return { name: 'Login' }
  }

  // 已登录但未验证身份：从服务端获取真实 isAdmin，防止 localStorage 篡改
  if (userStore.token && !authVerified) {
    await userStore.fetchUserInfo()
    authVerified = true
  }

  // 非管理员访问管理页 → 跳转首页
  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    return { name: 'Home' }
  }
})

// 登出时重置验证状态
export function resetAuthVerified() {
  authVerified = false
}

export default router
