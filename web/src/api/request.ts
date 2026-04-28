import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import i18n from '@/locales'

const { t } = i18n.global

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = token
  }
  return config
})

request.interceptors.response.use(
  (response) => {
    // Blob 响应直接返回（文件下载）
    if (response.config.responseType === 'blob') {
      return response.data
    }
    const res = response.data
    // 业务成功：解包 data 字段，调用方直接拿到业务数据
    if (res.code === 0) {
      return res.data
    }
    // 业务错误（非 code=0）
    if (res.code === 40001) {
      localStorage.removeItem('token')
      router.push('/login')
      ElMessage.error(t('auth.loginExpired'))
    } else {
      ElMessage.error(res.msg || t('auth.requestFailed'))
    }
    return Promise.reject(new Error(res.msg || t('auth.requestFailed')))
  },
  (error) => {
    if (error.response) {
      const data = error.response.data
      if (data?.code === 40001) {
        localStorage.removeItem('token')
        router.push('/login')
        ElMessage.error(t('auth.loginExpired'))
      } else {
        ElMessage.error(data?.msg || t('auth.requestFailed'))
      }
    } else {
      ElMessage.error(t('auth.networkError'))
    }
    return Promise.reject(error)
  },
)

export default request
