# Todo App API 接口文档

> Base URL: `http://{host}:8888`
> 版本: v1.3.0

---

## 通用说明

### 认证方式

需要认证的接口，在请求头中携带 JWT Token：

```
Authorization: <token>
```

登录接口返回的 `token` 字段即为认证凭证。

### 响应格式

**成功响应** — 直接返回业务数据，各接口不同。

**错误响应** — 统一格式：

```json
{
  "code": 20003,
  "msg": "密码错误"
}
```

### 错误码一览

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10001 | 服务器内部错误 |
| 10002 | 请求参数错误 |
| 20001 | 用户名已存在 |
| 20002 | 用户不存在 |
| 20003 | 密码错误 |
| 20004 | 用户已被禁用 |
| 20005 | 注册已关闭 |
| 20006 | 原密码错误 |
| 30001 | 任务不存在 |
| 40001 | 无权限操作 |
| 40002 | 需要管理员权限 |

---

## 一、公开接口（无需登录）

### 1.1 检查注册开关

判断系统是否允许新用户注册。

```
GET /api/v1/user/check-register
```

**响应示例：**

```json
{
  "allowRegister": true
}
```

---

### 1.2 用户注册

```
POST /api/v1/user/register
Content-Type: application/json
```

**请求参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名，3-20字符 |
| password | string | 是 | 密码，6-20字符 |

**请求示例：**

```json
{
  "username": "testuser",
  "password": "test1234"
}
```

**响应示例：**

```json
{
  "id": 2,
  "username": "testuser"
}
```

---

### 1.3 用户登录

```
POST /api/v1/user/login
Content-Type: application/json
```

**请求参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**请求示例：**

```json
{
  "username": "admin",
  "password": "admin123"
}
```

**响应示例：**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "isAdmin": 1
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| token | string | JWT 令牌，后续请求需携带 |
| isAdmin | int64 | 是否管理员：0=否，1=是 |

---

## 二、用户接口（需登录）

### 2.1 获取用户信息

```
GET /api/v1/user/info
Authorization: <token>
```

**响应示例：**

```json
{
  "id": 1,
  "username": "admin",
  "isAdmin": 1,
  "status": 1
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 用户ID |
| username | string | 用户名 |
| isAdmin | int64 | 是否管理员：0=否，1=是 |
| status | int64 | 状态：0=禁用，1=启用 |

---

### 2.2 修改密码

```
PUT /api/v1/user/password
Authorization: <token>
Content-Type: application/json
```

**请求参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| oldPassword | string | 是 | 原密码 |
| newPassword | string | 是 | 新密码，6-20字符 |

**请求示例：**

```json
{
  "oldPassword": "old123",
  "newPassword": "new456"
}
```

**响应示例：**

```json
{}
```

---

## 三、任务接口（需登录）

### 3.1 创建任务

```
POST /api/v1/task
Authorization: <token>
Content-Type: application/json
```

**请求参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 任务标题，最长100字符 |
| content | string | 否 | 任务内容，最长1000字符 |
| priority | int64 | 否 | 优先级：1=高，2=中(默认)，3=低 |
| categoryId | int64 | 否 | 分类ID |

**请求示例：**

```json
{
  "title": "完成项目文档",
  "content": "需要完成API文档和用户手册",
  "priority": 1,
  "categoryId": 1
}
```

**响应示例：**

```json
{
  "id": 1
}
```

---

### 3.2 任务列表

```
GET /api/v1/task?page=1&pageSize=10
Authorization: <token>
```

**查询参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int64 | 否 | 页码，默认1 |
| pageSize | int64 | 否 | 每页数量，默认10 |
| status | int64 | 否 | 筛选状态：0=待办，1=已完成 |
| categoryId | int64 | 否 | 筛选分类ID |
| priority | int64 | 否 | 筛选优先级：1/2/3 |
| keyword | string | 否 | 搜索关键词（标题模糊匹配） |

**响应示例：**

```json
{
  "total": 15,
  "list": [
    {
      "id": 1,
      "title": "完成项目文档",
      "content": "需要完成API文档和用户手册",
      "status": 0,
      "priority": 1,
      "categoryId": 1,
      "categoryName": "工作",
      "createTime": "2026-04-23 10:06"
    }
  ]
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| total | int64 | 总记录数 |
| list[].id | int64 | 任务ID |
| list[].title | string | 任务标题 |
| list[].content | string | 任务内容 |
| list[].status | int64 | 状态：0=待办，1=已完成 |
| list[].priority | int64 | 优先级：1=高，2=中，3=低 |
| list[].categoryId | int64 | 分类ID，0表示未分类 |
| list[].categoryName | string | 分类名称 |
| list[].createTime | string | 创建时间 |

---

### 3.3 任务详情

```
GET /api/v1/task/:id
Authorization: <token>
```

**路径参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 任务ID |

**响应示例：**

```json
{
  "id": 1,
  "title": "完成项目文档",
  "content": "需要完成API文档和用户手册",
  "status": 0,
  "priority": 1,
  "categoryId": 1,
  "categoryName": "工作",
  "createTime": "2026-04-23 10:06",
  "updateTime": "2026-04-23 10:06"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| updateTime | string | 最后更新时间 |

---

### 3.4 更新任务

```
PUT /api/v1/task/:id
Authorization: <token>
Content-Type: application/json
```

**路径参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 任务ID |

**请求参数：** 所有字段均可选，仅传需要更新的字段。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 否 | 任务标题 |
| content | string | 否 | 任务内容 |
| priority | int64 | 否 | 优先级：1/2/3 |
| categoryId | int64 | 否 | 分类ID |

**请求示例：**

```json
{
  "title": "更新后的标题",
  "priority": 2
}
```

**响应示例：**

```json
{}
```

---

### 3.5 切换任务状态

将待办切换为已完成，或将已完成切换为待办。

```
PATCH /api/v1/task/:id/toggle
Authorization: <token>
```

**路径参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 任务ID |

**响应示例：**

```json
{}
```

---

### 3.6 删除任务

```
DELETE /api/v1/task/:id
Authorization: <token>
```

**路径参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 任务ID |

**响应示例：**

```json
{}
```

---

### 3.7 批量操作任务

```
POST /api/v1/task/batch
Authorization: <token>
Content-Type: application/json
```

**请求参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ids | []int64 | 是 | 任务ID列表 |
| action | string | 是 | 操作类型：`complete`=完成，`undo`=取消完成，`delete`=删除 |

**请求示例：**

```json
{
  "ids": [1, 2, 3],
  "action": "complete"
}
```

**响应示例：**

```json
{}
```

---

## 四、分类接口（需登录）

### 4.1 分类列表

```
GET /api/v1/category
Authorization: <token>
```

**响应示例：**

```json
{
  "list": [
    { "id": 1, "name": "工作" },
    { "id": 2, "name": "生活" },
    { "id": 3, "name": "学习" }
  ]
}
```

> 系统预置分类（无 user_id）对所有用户可见，用户自建分类仅自己可见。

---

### 4.2 新增分类

```
POST /api/v1/category
Authorization: <token>
Content-Type: application/json
```

**请求参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 分类名称，最长20字符 |

**请求示例：**

```json
{
  "name": "健身"
}
```

**响应示例：**

```json
{
  "id": 4
}
```

---

## 五、统计接口（需登录）

### 5.1 统计概览

```
GET /api/v1/stat
Authorization: <token>
```

**响应示例：**

```json
{
  "total": 15,
  "done": 8,
  "todo": 7,
  "doneRate": 53
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| total | int64 | 任务总数 |
| done | int64 | 已完成数 |
| todo | int64 | 待办数 |
| doneRate | int64 | 完成率（百分比整数） |

---

## 六、管理员接口

> 以下接口需要管理员权限（isAdmin=1），否则返回 `40002 需要管理员权限`。

### 6.1 用户列表

```
GET /api/v1/admin/user?page=1&pageSize=10
Authorization: <admin-token>
```

**查询参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int64 | 否 | 页码，默认1 |
| pageSize | int64 | 否 | 每页数量，默认10 |
| keyword | string | 否 | 搜索关键词（用户名模糊匹配） |

**响应示例：**

```json
{
  "total": 3,
  "list": [
    {
      "id": 1,
      "username": "admin",
      "isAdmin": 1,
      "status": 1,
      "createTime": "2026-04-23 09:50"
    },
    {
      "id": 2,
      "username": "testuser",
      "isAdmin": 0,
      "status": 1,
      "createTime": "2026-04-23 09:59"
    }
  ]
}
```

---

### 6.2 重置用户密码

```
PUT /api/v1/admin/user/:id/password
Authorization: <admin-token>
Content-Type: application/json
```

**路径参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 用户ID |

**请求参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| newPassword | string | 是 | 新密码，6-20字符 |

**请求示例：**

```json
{
  "newPassword": "newpass123"
}
```

**响应示例：**

```json
{}
```

---

### 6.3 切换用户状态（禁用/启用）

```
PATCH /api/v1/admin/user/:id/toggle
Authorization: <admin-token>
```

**路径参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 用户ID |

**响应示例：**

```json
{}
```

---

### 6.4 删除用户

```
DELETE /api/v1/admin/user/:id
Authorization: <admin-token>
```

**路径参数：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 用户ID |

> 不可删除自己，否则返回错误。

**响应示例：**

```json
{}
```

---

### 6.5 系统配置列表

```
GET /api/v1/admin/config
Authorization: <admin-token>
```

**响应示例：**

```json
{
  "list": [
    {
      "key": "allow_register",
      "value": "true",
      "remark": "是否允许公开注册: true/false"
    },
    {
      "key": "site_name",
      "value": "Todo App",
      "remark": "站点名称"
    }
  ]
}
```

---

### 6.6 更新系统配置

```
PUT /api/v1/admin/config
Authorization: <admin-token>
Content-Type: application/json
```

**请求参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| key | string | 是 | 配置键 |
| value | string | 是 | 配置值 |

**请求示例：**

```json
{
  "key": "allow_register",
  "value": "false"
}
```

**响应示例：**

```json
{}
```

---

### 6.7 操作日志列表

```
GET /api/v1/admin/log/operation?page=1&pageSize=10
Authorization: <admin-token>
```

**查询参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int64 | 否 | 页码，默认1 |
| pageSize | int64 | 否 | 每页数量，默认10 |
| action | string | 否 | 筛选操作类型 |
| username | string | 否 | 筛选用户名 |

**响应示例：**

```json
{
  "total": 2,
  "list": [
    {
      "id": 1,
      "userId": 2,
      "username": "testuser",
      "action": "register",
      "targetType": "user",
      "targetId": 2,
      "detail": "用户注册",
      "ip": "",
      "createTime": "2026-04-23 09:59"
    }
  ]
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| action | string | 操作类型：register / login / create_task / update_task / delete_task 等 |
| targetType | string | 操作对象类型：user / task / category / config |
| targetId | int64 | 操作对象ID |
| detail | string | 操作详情 |
| ip | string | 操作者IP |

---

### 6.8 登录日志列表

```
GET /api/v1/admin/log/login?page=1&pageSize=10
Authorization: <admin-token>
```

**查询参数：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int64 | 否 | 页码，默认1 |
| pageSize | int64 | 否 | 每页数量，默认10 |
| username | string | 否 | 筛选用户名 |

**响应示例：**

```json
{
  "total": 3,
  "list": [
    {
      "id": 5,
      "userId": 1,
      "username": "admin",
      "ip": "",
      "status": 1,
      "remark": "登录成功",
      "createTime": "2026-04-23 10:11"
    },
    {
      "id": 4,
      "userId": 1,
      "username": "admin",
      "ip": "",
      "status": 0,
      "remark": "密码错误",
      "createTime": "2026-04-23 10:10"
    }
  ]
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| status | int64 | 状态：0=失败，1=成功 |
| remark | string | 备注信息（失败原因等） |
