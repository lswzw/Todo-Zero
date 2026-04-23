// API 响应类型定义

export interface TaskItem {
  id: number
  title: string
  content: string
  status: number
  priority: number
  categoryId: number
  categoryName: string
  createTime: string
}

export interface TaskDetail extends TaskItem {
  updateTime: string
}

export interface TaskListResp {
  list: TaskItem[]
  total: number
}

export interface StatResp {
  total: number
  done: number
  todo: number
  doneRate: number
}

export interface CategoryItem {
  id: number
  name: string
}

export interface CategoryListResp {
  list: CategoryItem[]
}

export interface UserListItem {
  id: number
  username: string
  isAdmin: number
  status: number
  createTime: string
}

export interface UserListResp {
  list: UserListItem[]
  total: number
}

export interface UserInfo {
  id: number
  username: string
  isAdmin: number
  status: number
}

export interface LoginResp {
  token: string
  isAdmin: number
}

export interface CheckRegisterResp {
  allowRegister: boolean
}

export interface ConfigItem {
  key: string
  value: string
  remark: string
  _value?: string
}

export interface ConfigListResp {
  list: ConfigItem[]
}

export interface OperationLogItem {
  id: number
  userId: number
  username: string
  action: string
  targetType: string
  targetId: number
  detail: string
  ip: string
  createTime: string
}

export interface OperationLogResp {
  list: OperationLogItem[]
  total: number
}

export interface LoginLogItem {
  id: number
  userId: number
  username: string
  ip: string
  status: number
  remark: string
  createTime: string
}

export interface LoginLogResp {
  list: LoginLogItem[]
  total: number
}
