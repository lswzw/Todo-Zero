import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUserInfo } from '@/api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userId = ref(Number(localStorage.getItem('userId')) || 0)
  const username = ref(localStorage.getItem('username') || '')
  const isAdmin = ref(false)

  function setLogin(data: { token: string; isAdmin: number }, name: string) {
    token.value = data.token
    isAdmin.value = data.isAdmin === 1
    username.value = name
    localStorage.setItem('token', data.token)
    // 不再存储 isAdmin 到 localStorage，防止篡改
  }

  // 从服务端获取真实用户信息，覆盖本地缓存
  async function fetchUserInfo() {
    try {
      const res = await getUserInfo() as any
      userId.value = res.id
      username.value = res.username
      isAdmin.value = res.isAdmin === 1
    } catch {
      // token 无效时清空
      logout()
    }
  }

  function setUserInfo(data: { id: number; username: string; isAdmin: number; status: number }) {
    userId.value = data.id
    username.value = data.username
    isAdmin.value = data.isAdmin === 1
    localStorage.setItem('username', data.username)
  }

  function logout() {
    token.value = ''
    userId.value = 0
    username.value = ''
    isAdmin.value = false
    localStorage.removeItem('token')
    localStorage.removeItem('userId')
    localStorage.removeItem('username')
  }

  return { token, userId, username, isAdmin, setLogin, setUserInfo, fetchUserInfo, logout }
})
