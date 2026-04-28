import request from './request'
import type {
  LoginResp,
  CheckRegisterResp,
  UserInfo,
  TaskListResp,
  TaskDetail,
  TrashListResp,
  CategoryListResp,
  StatResp,
  UserListResp,
  ConfigListResp,
  OperationLogResp,
  LoginLogResp,
  BackupListResp,
  TriggerBackupResp,
  RestoreBackupResp,
} from '@/types'

// 用户
export const checkRegister = () => request.get<never, CheckRegisterResp>('/user/check-register')
export const login = (data: { username: string; password: string }) =>
  request.post<never, LoginResp>('/user/login', data)
export const register = (data: { username: string; password: string }) => request.post('/user/register', data)
export const getUserInfo = () => request.get<never, UserInfo>('/user/info')
export const changePassword = (data: { oldPassword: string; newPassword: string }) =>
  request.put('/user/password', data)

// 任务
export const getTaskList = (params: Record<string, unknown>) => request.get<never, TaskListResp>('/task', { params })
export const getTaskDetail = (id: number) => request.get<never, TaskDetail>(`/task/${id}`)
export const createTask = (data: Record<string, unknown>) => request.post('/task', data)
export const updateTask = (id: number, data: Record<string, unknown>) => request.put(`/task/${id}`, data)
export const toggleTask = (id: number) => request.patch(`/task/${id}/toggle`)
export const deleteTask = (id: number) => request.delete(`/task/${id}`)
export const batchTask = (data: { ids: number[]; action: string }) => request.post('/task/batch', data)
export const sortTask = (data: { orders: { id: number; sortOrder: number }[] }) => request.put('/task/sort', data)
export const getTrashList = (params: Record<string, unknown>) => request.get<never, TrashListResp>('/task/trash', { params })
export const restoreTask = (id: number) => request.patch(`/task/${id}/restore`)
export const permanentDeleteTask = (id: number) => request.delete(`/task/${id}/permanent`)

// 导出任务（Blob 下载）
export const exportTasks = (params: Record<string, unknown>) =>
  request.get('/task/export', { params, responseType: 'blob' })

// 分类
export const getCategoryList = () => request.get<never, CategoryListResp>('/category')
export const createCategory = (data: { name: string; color?: string }) => request.post('/category', data)
export const updateCategory = (id: number, data: { name?: string; color?: string; icon?: string; sort?: number }) =>
  request.put(`/category/${id}`, data)
export const deleteCategory = (id: number) => request.delete(`/category/${id}`)

// 统计
export const getStat = () => request.get<never, StatResp>('/stat')

// 管理员 - 用户管理
export const getUserList = (params: Record<string, unknown>) =>
  request.get<never, UserListResp>('/admin/user', { params })
export const resetPassword = (id: number, data: { newPassword: string }) =>
  request.put(`/admin/user/${id}/password`, data)
export const toggleUserStatus = (id: number) => request.patch(`/admin/user/${id}/toggle`)
export const deleteUser = (id: number) => request.delete(`/admin/user/${id}`)

// 管理员 - 系统配置
export const getConfigList = () => request.get<never, ConfigListResp>('/admin/config')
export const updateConfig = (data: { key: string; value: string }) => request.put('/admin/config', data)

// 管理员 - 日志
export const getOperationLogList = (params: Record<string, unknown>) =>
  request.get<never, OperationLogResp>('/admin/log/operation', { params })
export const getLoginLogList = (params: Record<string, unknown>) =>
  request.get<never, LoginLogResp>('/admin/log/login', { params })

// 管理员 - 数据库备份
export const getBackupList = () => request.get<never, BackupListResp>('/admin/backup')
export const triggerBackup = () => request.post<never, TriggerBackupResp>('/admin/backup')
export const downloadBackup = (fileName: string) => `/api/v1/admin/backup/download/${fileName}`
export const restoreBackup = (fileName: string) => request.post<never, RestoreBackupResp>(`/admin/backup/restore/${fileName}`)
