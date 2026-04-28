import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUserInfo } from '@/api'
import type { UserInfo } from '@/types'
import { useStorage } from '@/composables/useStorage'

export const useUserStore = defineStore('user', () => {
  const { value: token, remove: removeToken } = useStorage<string>('token', '')
  const { value: userId, remove: removeUserId } = useStorage<number>('userId', 0)
  const { value: username, remove: removeUsername } = useStorage<string>('username', '')
  const isAdmin = ref(false)

  function setLogin(data: { token: string; isAdmin: number }, name: string) {
    token.value = data.token
    isAdmin.value = data.isAdmin === 1
    username.value = name
  }

  async function fetchUserInfo() {
    try {
      const res = (await getUserInfo()) as UserInfo
      userId.value = res.id
      username.value = res.username
      isAdmin.value = res.isAdmin === 1
    } catch {
      logout()
    }
  }

  function setUserInfo(data: { id: number; username: string; isAdmin: number; status: number }) {
    userId.value = data.id
    username.value = data.username
    isAdmin.value = data.isAdmin === 1
  }

  function logout() {
    removeToken()
    removeUserId()
    removeUsername()
    isAdmin.value = false
  }

  return { token, userId, username, isAdmin, setLogin, setUserInfo, fetchUserInfo, logout }
})
