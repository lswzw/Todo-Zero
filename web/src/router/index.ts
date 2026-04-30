import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import(/* webpackChunkName: "login" */ '@/views/login.vue'),
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import(/* webpackChunkName: "register" */ '@/views/register.vue'),
    },
    {
      path: '/',
      name: 'Home',
      component: () => import(/* webpackChunkName: "home" */ '@/views/home.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/task/:id',
      name: 'TaskDetail',
      component: () => import(/* webpackChunkName: "task-detail" */ '@/views/task-detail.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/trash',
      name: 'Trash',
      component: () => import(/* webpackChunkName: "trash" */ '@/views/trash.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin',
      name: 'Admin',
      component: () => import(/* webpackChunkName: "admin-layout" */ '@/views/admin/layout.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        { path: '', redirect: '/admin/user' },
        {
          path: 'user',
          name: 'AdminUser',
          component: () => import(/* webpackChunkName: "admin-user" */ '@/views/admin/user.vue'),
        },
        {
          path: 'config',
          name: 'AdminConfig',
          component: () => import(/* webpackChunkName: "admin-config" */ '@/views/admin/config.vue'),
        },
        {
          path: 'log',
          name: 'AdminLog',
          component: () => import(/* webpackChunkName: "admin-log" */ '@/views/admin/log.vue'),
        },
        {
          path: 'login-log',
          name: 'AdminLoginLog',
          component: () => import(/* webpackChunkName: "admin-login-log" */ '@/views/admin/login-log.vue'),
        },
        {
          path: 'backup',
          name: 'AdminBackup',
          component: () => import(/* webpackChunkName: "admin-backup" */ '@/views/admin/backup.vue'),
        },
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
