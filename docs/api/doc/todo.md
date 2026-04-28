### 1. "获取备份列表"

1. route definition

- Url: /api/v1/admin/backup
- Method: GET
- Request: `-`
- Response: `BackupListResp`

2. request definition



3. response definition



```golang
type BackupListResp struct {
	List []BackupItem `json:"list"`
}
```

### 2. "触发数据库备份"

1. route definition

- Url: /api/v1/admin/backup
- Method: POST
- Request: `TriggerBackupReq`
- Response: `TriggerBackupResp`

2. request definition



```golang
type TriggerBackupReq struct {
}
```


3. response definition



```golang
type TriggerBackupResp struct {
	FileName string `json:"fileName"`
	FileSize int64 `json:"fileSize"`
}
```

### 3. "下载备份文件"

1. route definition

- Url: /api/v1/admin/backup/download/:fileName
- Method: GET
- Request: `DownloadBackupReq`
- Response: `-`

2. request definition



```golang
type DownloadBackupReq struct {
	FileName string `path:"fileName"`
}
```


3. response definition


### 4. "从备份恢复数据库"

1. route definition

- Url: /api/v1/admin/backup/restore/:fileName
- Method: POST
- Request: `RestoreBackupReq`
- Response: `RestoreBackupResp`

2. request definition



```golang
type RestoreBackupReq struct {
	FileName string `path:"fileName"`
}
```


3. response definition



```golang
type RestoreBackupResp struct {
	PreRestoreBackup string `json:"preRestoreBackup"`
}
```

### 5. "获取系统配置"

1. route definition

- Url: /api/v1/admin/config
- Method: GET
- Request: `-`
- Response: `ConfigListResp`

2. request definition



3. response definition



```golang
type ConfigListResp struct {
	List []ConfigItem `json:"list"`
}
```

### 6. "更新系统配置"

1. route definition

- Url: /api/v1/admin/config
- Method: PUT
- Request: `UpdateConfigReq`
- Response: `UpdateConfigResp`

2. request definition



```golang
type UpdateConfigReq struct {
	Key string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}
```


3. response definition



```golang
type UpdateConfigResp struct {
}
```

### 7. "获取登录日志"

1. route definition

- Url: /api/v1/admin/log/login
- Method: GET
- Request: `LoginLogReq`
- Response: `LoginLogResp`

2. request definition



```golang
type LoginLogReq struct {
	Page int64 `form:"page,default=1"`
	PageSize int64 `form:"pageSize,default=10"`
	Username string `form:"username,optional"`
}
```


3. response definition



```golang
type LoginLogResp struct {
	Total int64 `json:"total"`
	List []LoginLogItem `json:"list"`
}
```

### 8. "获取操作日志"

1. route definition

- Url: /api/v1/admin/log/operation
- Method: GET
- Request: `OperationLogReq`
- Response: `OperationLogResp`

2. request definition



```golang
type OperationLogReq struct {
	Page int64 `form:"page,default=1"`
	PageSize int64 `form:"pageSize,default=10"`
	Action string `form:"action,optional"`
	Username string `form:"username,optional"`
}
```


3. response definition



```golang
type OperationLogResp struct {
	Total int64 `json:"total"`
	List []OperationLogItem `json:"list"`
}
```

### 9. "获取用户列表"

1. route definition

- Url: /api/v1/admin/user
- Method: GET
- Request: `UserListReq`
- Response: `UserListResp`

2. request definition



```golang
type UserListReq struct {
	Page int64 `form:"page,default=1"`
	PageSize int64 `form:"pageSize,default=10"`
	Keyword string `form:"keyword,optional"`
}
```


3. response definition



```golang
type UserListResp struct {
	Total int64 `json:"total"`
	List []UserListItem `json:"list"`
}
```

### 10. "删除用户"

1. route definition

- Url: /api/v1/admin/user/:id
- Method: DELETE
- Request: `DeleteUserReq`
- Response: `DeleteUserResp`

2. request definition



```golang
type DeleteUserReq struct {
	Id int64 `path:"id"`
}
```


3. response definition



```golang
type DeleteUserResp struct {
}
```

### 11. "重置用户密码"

1. route definition

- Url: /api/v1/admin/user/:id/password
- Method: PUT
- Request: `ResetPasswordReq`
- Response: `ResetPasswordResp`

2. request definition



```golang
type ResetPasswordReq struct {
	Id int64 `path:"id"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=20"`
}
```


3. response definition



```golang
type ResetPasswordResp struct {
}
```

### 12. "切换用户状态(禁用/启用)"

1. route definition

- Url: /api/v1/admin/user/:id/toggle
- Method: PATCH
- Request: `ToggleUserStatusReq`
- Response: `ToggleUserStatusResp`

2. request definition



```golang
type ToggleUserStatusReq struct {
	Id int64 `path:"id"`
}
```


3. response definition



```golang
type ToggleUserStatusResp struct {
}
```

### 13. "获取分类列表"

1. route definition

- Url: /api/v1/category
- Method: GET
- Request: `-`
- Response: `CategoryListResp`

2. request definition



3. response definition



```golang
type CategoryListResp struct {
	List []CategoryItem `json:"list"`
}
```

### 14. "创建分类"

1. route definition

- Url: /api/v1/category
- Method: POST
- Request: `CreateCategoryReq`
- Response: `CreateCategoryResp`

2. request definition



```golang
type CreateCategoryReq struct {
	Name string `json:"name" validate:"required,max=20"`
	Color string `json:"color,optional"`
}
```


3. response definition



```golang
type CreateCategoryResp struct {
	Id int64 `json:"id"`
}
```

### 15. "更新分类"

1. route definition

- Url: /api/v1/category/:id
- Method: PUT
- Request: `UpdateCategoryReq`
- Response: `UpdateCategoryResp`

2. request definition



```golang
type UpdateCategoryReq struct {
	Id int64 `path:"id"`
	Name *string `json:"name,optional"`
	Color *string `json:"color,optional"`
	Icon *string `json:"icon,optional"`
	Sort *int64 `json:"sort,optional"`
}
```


3. response definition



```golang
type UpdateCategoryResp struct {
}
```

### 16. "删除分类"

1. route definition

- Url: /api/v1/category/:id
- Method: DELETE
- Request: `DeleteCategoryReq`
- Response: `DeleteCategoryResp`

2. request definition



```golang
type DeleteCategoryReq struct {
	Id int64 `path:"id"`
}
```


3. response definition



```golang
type DeleteCategoryResp struct {
}
```

### 17. "获取任务统计概览"

1. route definition

- Url: /api/v1/stat
- Method: GET
- Request: `-`
- Response: `StatResp`

2. request definition



3. response definition



```golang
type StatResp struct {
	Total int64 `json:"total"`
	Done int64 `json:"done"`
	Todo int64 `json:"todo"`
	DoneRate int64 `json:"doneRate"`
}
```

### 18. "创建任务"

1. route definition

- Url: /api/v1/task
- Method: POST
- Request: `CreateTaskReq`
- Response: `CreateTaskResp`

2. request definition



```golang
type CreateTaskReq struct {
	Title string `json:"title" validate:"required,max=100"`
	Content string `json:"content,optional" validate:"max=1000"`
	Priority int64 `json:"priority,optional,default=2" options:"1|2|3"`
	CategoryId int64 `json:"categoryId,optional"`
	StartTime string `json:"startTime,optional"`
	EndTime string `json:"endTime,optional"`
	Reminder string `json:"reminder,optional"`
	Tags string `json:"tags,optional" validate:"max=200"`
}
```


3. response definition



```golang
type CreateTaskResp struct {
	Id int64 `json:"id"`
}
```

### 19. "获取任务列表"

1. route definition

- Url: /api/v1/task
- Method: GET
- Request: `TaskListReq`
- Response: `TaskListResp`

2. request definition



```golang
type TaskListReq struct {
	Page int64 `form:"page,default=1"`
	PageSize int64 `form:"pageSize,default=10"`
	Status int64 `form:"status,optional" options:"0|1"`
	CategoryId int64 `form:"categoryId,optional"`
	Priority int64 `form:"priority,optional" options:"1|2|3"`
	Keyword string `form:"keyword,optional"`
}
```


3. response definition



```golang
type TaskListResp struct {
	Total int64 `json:"total"`
	List []TaskItem `json:"list"`
}
```

### 20. "更新任务"

1. route definition

- Url: /api/v1/task/:id
- Method: PUT
- Request: `UpdateTaskReq`
- Response: `UpdateTaskResp`

2. request definition



```golang
type UpdateTaskReq struct {
	Id int64 `path:"id"`
	Title string `json:"title,optional" validate:"max=100"`
	Content string `json:"content,optional" validate:"max=1000"`
	Priority int64 `json:"priority,optional" options:"1|2|3"`
	CategoryId int64 `json:"categoryId,optional"`
	StartTime string `json:"startTime,optional"`
	EndTime string `json:"endTime,optional"`
	Reminder string `json:"reminder,optional"`
	Tags string `json:"tags,optional" validate:"max=200"`
}
```


3. response definition



```golang
type UpdateTaskResp struct {
}
```

### 21. "删除任务(移至回收站)"

1. route definition

- Url: /api/v1/task/:id
- Method: DELETE
- Request: `DeleteTaskReq`
- Response: `DeleteTaskResp`

2. request definition



```golang
type DeleteTaskReq struct {
	Id int64 `path:"id"`
}
```


3. response definition



```golang
type DeleteTaskResp struct {
}
```

### 22. "获取任务详情"

1. route definition

- Url: /api/v1/task/:id
- Method: GET
- Request: `TaskDetailReq`
- Response: `TaskDetailResp`

2. request definition



```golang
type TaskDetailReq struct {
	Id int64 `path:"id"`
}
```


3. response definition



```golang
type TaskDetailResp struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Status int64 `json:"status"`
	Priority int64 `json:"priority"`
	CategoryId int64 `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	Reminder string `json:"reminder"`
	Tags string `json:"tags"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}
```

### 23. "永久删除任务"

1. route definition

- Url: /api/v1/task/:id/permanent
- Method: DELETE
- Request: `PermanentDeleteTaskReq`
- Response: `PermanentDeleteTaskResp`

2. request definition



```golang
type PermanentDeleteTaskReq struct {
	Id int64 `path:"id"`
}
```


3. response definition



```golang
type PermanentDeleteTaskResp struct {
}
```

### 24. "恢复已删除的任务"

1. route definition

- Url: /api/v1/task/:id/restore
- Method: PATCH
- Request: `RestoreTaskReq`
- Response: `RestoreTaskResp`

2. request definition



```golang
type RestoreTaskReq struct {
	Id int64 `path:"id"`
}
```


3. response definition



```golang
type RestoreTaskResp struct {
}
```

### 25. "切换任务完成状态"

1. route definition

- Url: /api/v1/task/:id/toggle
- Method: PATCH
- Request: `ToggleTaskReq`
- Response: `ToggleTaskResp`

2. request definition



```golang
type ToggleTaskReq struct {
	Id int64 `path:"id"`
}
```


3. response definition



```golang
type ToggleTaskResp struct {
}
```

### 26. "批量操作任务"

1. route definition

- Url: /api/v1/task/batch
- Method: POST
- Request: `BatchTaskReq`
- Response: `BatchTaskResp`

2. request definition



```golang
type BatchTaskReq struct {
	Ids []int64 `json:"ids" validate:"required"`
	Action string `json:"action" validate:"required" options:"complete|undo|delete|restore"`
}
```


3. response definition



```golang
type BatchTaskResp struct {
}
```

### 27. "获取回收站列表"

1. route definition

- Url: /api/v1/task/trash
- Method: GET
- Request: `TrashListReq`
- Response: `TrashListResp`

2. request definition



```golang
type TrashListReq struct {
	Page int64 `form:"page,default=1"`
	PageSize int64 `form:"pageSize,default=10"`
}
```


3. response definition



```golang
type TrashListResp struct {
	Total int64 `json:"total"`
	List []TrashItem `json:"list"`
}
```

### 28. "检查是否允许注册"

1. route definition

- Url: /api/v1/user/check-register
- Method: GET
- Request: `-`
- Response: `CheckRegisterResp`

2. request definition



3. response definition



```golang
type CheckRegisterResp struct {
	AllowRegister bool `json:"allowRegister"`
}
```

### 29. "用户登录"

1. route definition

- Url: /api/v1/user/login
- Method: POST
- Request: `LoginReq`
- Response: `LoginResp`

2. request definition



```golang
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
```


3. response definition



```golang
type LoginResp struct {
	Token string `json:"token"`
	IsAdmin int64 `json:"isAdmin"`
}
```

### 30. "用户注册"

1. route definition

- Url: /api/v1/user/register
- Method: POST
- Request: `RegisterReq`
- Response: `RegisterResp`

2. request definition



```golang
type RegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}
```


3. response definition



```golang
type RegisterResp struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
}
```

### 31. "获取当前用户信息"

1. route definition

- Url: /api/v1/user/info
- Method: GET
- Request: `-`
- Response: `UserInfoResp`

2. request definition



3. response definition



```golang
type UserInfoResp struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	IsAdmin int64 `json:"isAdmin"`
	Status int64 `json:"status"`
}
```

### 32. "修改密码"

1. route definition

- Url: /api/v1/user/password
- Method: PUT
- Request: `ChangePasswordReq`
- Response: `ChangePasswordResp`

2. request definition



```golang
type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=20"`
}
```


3. response definition



```golang
type ChangePasswordResp struct {
}
```

