import request from './request'

// 用户
export const checkRegister = () => request.get('/user/check-register')
export const login = (data: { username: string; password: string }) => request.post('/user/login', data)
export const register = (data: { username: string; password: string }) => request.post('/user/register', data)
export const getUserInfo = () => request.get('/user/info')
export const changePassword = (data: { oldPassword: string; newPassword: string }) => request.put('/user/password', data)

// 任务
export const getTaskList = (params: Record<string, unknown>) => request.get('/task', { params })
export const getTaskDetail = (id: number) => request.get(`/task/${id}`)
export const createTask = (data: Record<string, unknown>) => request.post('/task', data)
export const updateTask = (id: number, data: Record<string, unknown>) => request.put(`/task/${id}`, data)
export const toggleTask = (id: number) => request.patch(`/task/${id}/toggle`)
export const deleteTask = (id: number) => request.delete(`/task/${id}`)
export const batchTask = (data: { ids: number[]; action: string }) => request.post('/task/batch', data)

// 分类
export const getCategoryList = () => request.get('/category')
export const createCategory = (data: { name: string }) => request.post('/category', data)

// 统计
export const getStat = () => request.get('/stat')

// 管理员 - 用户管理
export const getUserList = (params: Record<string, unknown>) => request.get('/admin/user', { params })
export const resetPassword = (id: number, data: { newPassword: string }) => request.put(`/admin/user/${id}/password`, data)
export const toggleUserStatus = (id: number) => request.patch(`/admin/user/${id}/toggle`)
export const deleteUser = (id: number) => request.delete(`/admin/user/${id}`)

// 管理员 - 系统配置
export const getConfigList = () => request.get('/admin/config')
export const updateConfig = (data: { key: string; value: string }) => request.put('/admin/config', data)

// 管理员 - 日志
export const getOperationLogList = (params: Record<string, unknown>) => request.get('/admin/log/operation', { params })
export const getLoginLogList = (params: Record<string, unknown>) => request.get('/admin/log/login', { params })
