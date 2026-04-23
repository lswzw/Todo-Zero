import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userId = ref(Number(localStorage.getItem('userId')) || 0)
  const username = ref(localStorage.getItem('username') || '')
  const isAdmin = ref(localStorage.getItem('isAdmin') === '1')

  function setLogin(data: { token: string; isAdmin: number }, name: string) {
    token.value = data.token
    isAdmin.value = data.isAdmin === 1
    username.value = name
    localStorage.setItem('token', data.token)
    localStorage.setItem('isAdmin', String(data.isAdmin))
    localStorage.setItem('username', name)
  }

  function setUserInfo(data: { id: number; username: string; isAdmin: number; status: number }) {
    userId.value = data.id
    username.value = data.username
    isAdmin.value = data.isAdmin === 1
    localStorage.setItem('userId', String(data.id))
    localStorage.setItem('username', data.username)
    localStorage.setItem('isAdmin', String(data.isAdmin))
  }

  function logout() {
    token.value = ''
    userId.value = 0
    username.value = ''
    isAdmin.value = false
    localStorage.removeItem('token')
    localStorage.removeItem('userId')
    localStorage.removeItem('username')
    localStorage.removeItem('isAdmin')
  }

  return { token, userId, username, isAdmin, setLogin, setUserInfo, logout }
})
