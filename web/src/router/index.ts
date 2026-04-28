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

let authVerified = false

router.beforeEach(async (to) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.token) {
    return { name: 'Login' }
  }

  if (userStore.token && !authVerified) {
    await userStore.fetchUserInfo()
    authVerified = true
  }

  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    return { name: 'Home' }
  }
})

export function resetAuthVerified() {
  authVerified = false
}

export default router
