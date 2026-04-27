#!/bin/bash
# ============================================
# Todo App 集成测试脚本 (v1.4.0)
# 用法: bash test.sh
# 依赖: curl, jq
# 说明: 自动删除数据库、重启服务、执行测试、清理还原
# ============================================

set -uo pipefail

# ========== 配置 ==========
BASE_URL="http://localhost:8888"
SERVER_DIR="$(cd "$(dirname "$0")/server" && pwd)"
DB_FILE="${SERVER_DIR}/data/todo.db"
SERVER_BIN="${SERVER_DIR}/todo-api"
SERVER_CONFIG="${SERVER_DIR}/etc/todo-api.yaml"
SERVER_PID=""
PORT=8888

ADMIN_USER="admin"
ADMIN_PASS="admin123"        # 初始管理员密码（init.sql 预置，不受复杂度限制）
ADMIN_NEW_PASS="admin456"    # 修改后密码（满足复杂度：字母+数字）
TEST_USER="zhangsan"
TEST_USER_PASS="test123"     # 注册密码（满足复杂度：字母+数字）
TEST_USER_NEW_PASS="test654" # 重置后密码（满足复杂度：字母+数字）
WEAK_PASS="123456"           # 不满足复杂度的密码（纯数字）

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

# 计数器
PASS=0
FAIL=0

# ========== 工具函数 ==========
pass() { ((PASS++)); printf "  ${GREEN}✓ PASS${NC} %s\n" "$1"; }
fail() { ((FAIL++)); printf "  ${RED}✗ FAIL${NC} %s\n" "$1"; }
section() { printf "\n${YELLOW}━━━ %s ━━━${NC}\n" "$1"; }

get() {
    local path="$1" token="${2:-}"
    if [ -n "$token" ]; then
        curl -s "${BASE_URL}${path}" -H "Authorization: Bearer ${token}"
    else
        curl -s "${BASE_URL}${path}"
    fi
}

post() {
    local path="$1" body="$2" token="${3:-}"
    if [ -n "$token" ]; then
        curl -s -X POST "${BASE_URL}${path}" \
            -H "Authorization: Bearer ${token}" \
            -H "Content-Type: application/json" -d "$body"
    else
        curl -s -X POST "${BASE_URL}${path}" \
            -H "Content-Type: application/json" -d "$body"
    fi
}

put() {
    local path="$1" body="$2" token="${3:-}"
    if [ -n "$token" ]; then
        curl -s -X PUT "${BASE_URL}${path}" \
            -H "Authorization: Bearer ${token}" \
            -H "Content-Type: application/json" -d "$body"
    else
        curl -s -X PUT "${BASE_URL}${path}" \
            -H "Content-Type: application/json" -d "$body"
    fi
}

patch_req() {
    local path="$1" token="$2"
    curl -s -X PATCH "${BASE_URL}${path}" \
        -H "Authorization: Bearer ${token}"
}

del() {
    local path="$1" token="$2"
    curl -s -X DELETE "${BASE_URL}${path}" \
        -H "Authorization: Bearer ${token}"
}

jval() { echo "$1" | jq -r "$2" 2>/dev/null; }

assert_eq() {
    local resp="$1" key="$2" expected="$3" desc="$4"
    local actual
    actual=$(jval "$resp" ".$key")
    if [ "$actual" = "$expected" ]; then
        pass "${desc}: ${key}=${actual}"
    else
        fail "${desc}: ${key}=${actual} (expected ${expected})"
    fi
}

assert_not_empty() {
    local resp="$1" key="$2" desc="$3"
    local actual
    actual=$(jval "$resp" ".$key")
    if [ -n "$actual" ] && [ "$actual" != "null" ]; then
        pass "${desc}: ${key}=${actual}"
    else
        fail "${desc}: ${key} 为空"
    fi
}

assert_ok() {
    local resp="$1" desc="$2"
    local code
    code=$(jval "$resp" ".code")
    if [ "$code" = "0" ]; then
        pass "${desc}: 成功"
    else
        fail "${desc}: code=${code} msg=$(jval "$resp" ".msg")"
    fi
}

assert_err() {
    local resp="$1" desc="$2"
    local code
    code=$(jval "$resp" ".code")
    if [ -n "$code" ] && [ "$code" != "null" ] && [ "$code" != "0" ]; then
        pass "${desc}: 正确拒绝 code=${code}"
    else
        fail "${desc}: 不应成功 → $resp"
    fi
}

assert_err_msg() {
    local resp="$1" expected_msg="$2" desc="$3"
    local code msg
    code=$(jval "$resp" ".code")
    msg=$(jval "$resp" ".msg")
    if [ -n "$code" ] && [ "$code" != "null" ] && [ "$code" != "0" ]; then
        if echo "$msg" | grep -q "$expected_msg"; then
            pass "${desc}: 正确拒绝 msg=${msg}"
        else
            fail "${desc}: 拒绝但消息不匹配 msg=${msg} (expected: ${expected_msg})"
        fi
    else
        fail "${desc}: 不应成功 → $resp"
    fi
}

assert_http() {
    local path="$1" expected="$2"
    local actual
    actual=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}${path}")
    if [ "$actual" = "$expected" ]; then
        pass "GET ${path} → ${actual}"
    else
        fail "GET ${path} → ${actual} (expected ${expected})"
    fi
}

assert_http_post() {
    local path="$1" body="$2" expected="$3"
    local actual
    actual=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}${path}" \
        -H "Content-Type: application/json" -d "$body")
    if [ "$actual" = "$expected" ]; then
        pass "POST ${path} → ${actual}"
    else
        fail "POST ${path} → ${actual} (expected ${expected})"
    fi
}

# ========== 服务管理 ==========
start_server() {
    section "环境准备"

    for cmd in curl jq; do
        if ! command -v "$cmd" &>/dev/null; then
            printf "${RED}错误: %s 未安装${NC}\n" "$cmd"
            exit 1
        fi
    done
    pass "依赖检查: curl, jq"

    if [ ! -x "$SERVER_BIN" ]; then
        printf "${RED}错误: 服务未编译 %s${NC}\n" "$SERVER_BIN"
        printf "请先编译: cd server && go build -o todo-api .\n"
        exit 1
    fi

    # 杀掉占用端口的旧进程
    local old_pid
    old_pid=$(lsof -ti :${PORT} 2>/dev/null || true)
    if [ -n "$old_pid" ]; then
        kill $old_pid 2>/dev/null
        sleep 1
        old_pid=$(lsof -ti :${PORT} 2>/dev/null || true)
        if [ -n "$old_pid" ]; then
            kill -9 $old_pid 2>/dev/null
            sleep 1
        fi
        pass "停止旧服务 (pid: $old_pid)"
    fi

    # 删除数据库
    if [ -f "$DB_FILE" ]; then
        rm -f "$DB_FILE"
        pass "删除旧数据库"
    else
        pass "无旧数据库"
    fi
    rm -f "${DB_FILE}-wal" "${DB_FILE}-shm" 2>/dev/null

    # 启动服务
    cd "$SERVER_DIR"
    "$SERVER_BIN" -f "$SERVER_CONFIG" >/dev/null 2>&1 &
    SERVER_PID=$!
    cd - >/dev/null

    # 等待服务就绪
    local retry=0
    while [ $retry -lt 20 ]; do
        if curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/health" 2>/dev/null | grep -q "200"; then
            pass "服务启动成功 (pid: ${SERVER_PID})"
            return
        fi
        sleep 0.5
        ((retry++))
    done

    printf "${RED}错误: 服务启动超时${NC}\n"
    kill $SERVER_PID 2>/dev/null
    exit 1
}

cleanup() {
    printf "\n${CYAN}━━━ 清理还原 ━━━${NC}\n"
    if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
        kill "$SERVER_PID" 2>/dev/null
        wait "$SERVER_PID" 2>/dev/null
        pass "服务已停止"
    fi
    rm -f "$DB_FILE" "${DB_FILE}-wal" "${DB_FILE}-shm" 2>/dev/null
    pass "数据库已删除"
}

trap cleanup EXIT

# ========== 启动服务 ==========
start_server

# ========== 1. 健康检查端点 ==========
section "1. 健康检查端点"

RESP=$(get "/health")
if echo "$RESP" | grep -q '"status":"ok"' 2>/dev/null; then
    pass "GET /health → status=ok"
else
    fail "GET /health → 异常: $RESP"
fi

# ========== 2. 静态文件服务 ==========
section "2. 静态文件服务"
assert_http "/" "200"
assert_http "/login" "200"

# ========== 3. 无需认证接口 ==========
section "3. 无需认证接口"

RESP=$(get "/api/v1/user/check-register")
assert_eq "$RESP" "data.allowRegister" "false" "检查注册开关"

RESP=$(post "/api/v1/user/register" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_PASS}\"}")
assert_err "$RESP" "注册(已关闭)"

# 登录 — 获取管理员 token
RESP=$(post "/api/v1/user/login" "{\"username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_PASS}\"}")
ADMIN_TOKEN=$(jval "$RESP" ".data.token")
if [ -n "$ADMIN_TOKEN" ] && [ "$ADMIN_TOKEN" != "null" ]; then
    pass "管理员登录: 获取到 token"
else
    fail "管理员登录: 失败 → $RESP"
    printf "${RED}无法继续测试${NC}\n"
    exit 1
fi
assert_eq "$RESP" "data.isAdmin" "1" "管理员角色"

# 错误密码登录
RESP=$(post "/api/v1/user/login" "{\"username\":\"${ADMIN_USER}\",\"password\":\"wrongpwd\"}")
assert_err "$RESP" "错误密码登录"

# ========== 4. 开启注册 + 注册用户 ==========
section "4. 开启注册 + 注册用户"

RESP=$(put "/api/v1/admin/config" "{\"key\":\"allow_register\",\"value\":\"true\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "开启注册"

# 注册 — 正常密码
RESP=$(post "/api/v1/user/register" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_PASS}\"}")
assert_not_empty "$RESP" "data.id" "注册用户(${TEST_USER})"
TEST_USER_ID=$(jval "$RESP" ".data.id")

# 注册 — 弱密码(纯数字)
RESP=$(post "/api/v1/user/register" "{\"username\":\"weakuser\",\"password\":\"${WEAK_PASS}\"}")
assert_err_msg "$RESP" "字母和数字" "注册弱密码拒绝"

# 注册 — 空密码
RESP=$(post "/api/v1/user/register" "{\"username\":\"emptypwd\",\"password\":\"\"}")
assert_err "$RESP" "注册空密码拒绝"

# 注册 — 重复用户名
RESP=$(post "/api/v1/user/register" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_PASS}\"}")
assert_err "$RESP" "重复用户名注册拒绝"

# 普通用户登录
RESP=$(post "/api/v1/user/login" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_PASS}\"}")
USER_TOKEN=$(jval "$RESP" ".data.token")
if [ -n "$USER_TOKEN" ] && [ "$USER_TOKEN" != "null" ]; then
    pass "普通用户登录: 获取到 token"
else
    fail "普通用户登录: 失败 → $RESP"
fi

# ========== 5. 用户信息 ==========
section "5. 用户信息"

RESP=$(get "/api/v1/user/info" "$ADMIN_TOKEN")
assert_eq "$RESP" "data.username" "admin" "获取管理员信息"
assert_eq "$RESP" "data.isAdmin" "1" "管理员角色"

RESP=$(get "/api/v1/user/info" "$USER_TOKEN")
assert_eq "$RESP" "data.username" "zhangsan" "获取普通用户信息"
assert_eq "$RESP" "data.isAdmin" "0" "普通用户角色"

# 无 token 访问受保护接口
HTTP_CODE=$(curl -s -o /tmp/notoken_resp.json -w "%{http_code}" "${BASE_URL}/api/v1/user/info")
if [ "$HTTP_CODE" = "401" ]; then
    pass "无token访问用户信息拒绝: HTTP 401"
else
    RESP=$(cat /tmp/notoken_resp.json)
    assert_err "$RESP" "无token访问用户信息拒绝"
fi

# ========== 6. 修改密码 ==========
section "6. 修改密码"

# 新密码不满足复杂度
RESP=$(put "/api/v1/user/password" "{\"oldPassword\":\"${ADMIN_PASS}\",\"newPassword\":\"${WEAK_PASS}\"}" "$ADMIN_TOKEN")
assert_err_msg "$RESP" "字母和数字" "修改密码弱密码拒绝"

# 正常修改密码
RESP=$(put "/api/v1/user/password" "{\"oldPassword\":\"${ADMIN_PASS}\",\"newPassword\":\"${ADMIN_NEW_PASS}\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "修改密码"

# 用新密码登录
RESP=$(post "/api/v1/user/login" "{\"username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_NEW_PASS}\"}")
NEW_TOKEN=$(jval "$RESP" ".data.token")
if [ -n "$NEW_TOKEN" ] && [ "$NEW_TOKEN" != "null" ]; then
    pass "新密码登录: 成功"
    ADMIN_TOKEN="$NEW_TOKEN"
else
    fail "新密码登录: 失败 → $RESP"
fi

# 旧密码不应再登录
RESP=$(post "/api/v1/user/login" "{\"username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_PASS}\"}")
assert_err "$RESP" "旧密码登录"

# ========== 7. 分类 CRUD ==========
section "7. 分类 CRUD"

# 7.1 列表 — 初始4个系统分类
RESP=$(get "/api/v1/category" "$ADMIN_TOKEN")
assert_eq "$RESP" "data.list | length" "4" "分类列表(4个初始分类)"

# 验证 CategoryItem 新增字段
FIRST_CAT_NAME=$(jval "$RESP" ".data.list[0].name")
FIRST_CAT_COLOR=$(jval "$RESP" ".data.list[0].color")
FIRST_CAT_IS_SYSTEM=$(jval "$RESP" ".data.list[0].isSystem")
if [ -n "$FIRST_CAT_NAME" ] && [ "$FIRST_CAT_NAME" != "null" ]; then
    pass "分类字段验证: name=${FIRST_CAT_NAME}, color=${FIRST_CAT_COLOR}, isSystem=${FIRST_CAT_IS_SYSTEM}"
else
    fail "分类字段验证: 缺少扩展字段"
fi

# 7.2 创建分类（带颜色）
RESP=$(post "/api/v1/category" "{\"name\":\"娱乐\",\"color\":\"#ff6b6b\"}" "$ADMIN_TOKEN")
assert_not_empty "$RESP" "data.id" "创建分类(带颜色)"
ENTERTAIN_ID=$(jval "$RESP" ".data.id")

# 创建分类（不带颜色，使用默认）
RESP=$(post "/api/v1/category" "{\"name\":\"学习\"}" "$ADMIN_TOKEN")
assert_not_empty "$RESP" "data.id" "创建分类(默认颜色)"
STUDY_ID=$(jval "$RESP" ".data.id")

# 创建分类 — 空名称
RESP=$(post "/api/v1/category" "{\"name\":\"\"}" "$ADMIN_TOKEN")
assert_err "$RESP" "创建分类空名称拒绝"

# 创建分类 — 名称超长
RESP=$(post "/api/v1/category" "{\"name\":\"超长分类名称超长分类名称超长分类名称超长分类名称\"}" "$ADMIN_TOKEN")
assert_err "$RESP" "创建分类名称超长拒绝"

# 7.3 列表 — 新增后6个
RESP=$(get "/api/v1/category" "$ADMIN_TOKEN")
assert_eq "$RESP" "data.list | length" "6" "分类列表(新增后6个)"

# 7.4 更新分类
RESP=$(put "/api/v1/category/${ENTERTAIN_ID}" "{\"name\":\"运动\",\"color\":\"#52c41a\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "更新分类"

# 验证更新后的分类名称
RESP=$(get "/api/v1/category" "$ADMIN_TOKEN")
UPDATED_NAME=$(echo "$RESP" | jq -r ".data.list[] | select(.id==${ENTERTAIN_ID}) | .name" 2>/dev/null)
if [ "$UPDATED_NAME" = "运动" ]; then
    pass "更新后分类名称: ${UPDATED_NAME}"
else
    fail "更新后分类名称: ${UPDATED_NAME} (expected 运动)"
fi

# 7.5 删除非系统分类
RESP=$(del "/api/v1/category/${STUDY_ID}" "$ADMIN_TOKEN")
assert_ok "$RESP" "删除非系统分类"

RESP=$(get "/api/v1/category" "$ADMIN_TOKEN")
assert_eq "$RESP" "data.list | length" "5" "分类列表(删除后5个)"

# 7.6 删除系统分类 — 应被拒绝
# 获取第一个系统分类ID
FIRST_SYS_ID=$(jval "$RESP" ".data.list[0].id")
RESP=$(del "/api/v1/category/${FIRST_SYS_ID}" "$ADMIN_TOKEN")
assert_err "$RESP" "删除系统分类拒绝"

# 7.7 普通用户不能操作别人的分类
# admin 创建的分类，zhangsan 尝试修改
RESP=$(put "/api/v1/category/${ENTERTAIN_ID}" "{\"name\":\"hack\"}" "$USER_TOKEN")
assert_err "$RESP" "跨用户修改分类拒绝"

RESP=$(del "/api/v1/category/${ENTERTAIN_ID}" "$USER_TOKEN")
assert_err "$RESP" "跨用户删除分类拒绝"

# 7.8 普通用户创建自己的分类
RESP=$(post "/api/v1/category" "{\"name\":\"我的分类\",\"color\":\"#722ed1\"}" "$USER_TOKEN")
assert_not_empty "$RESP" "data.id" "普通用户创建分类"
USER_CAT_ID=$(jval "$RESP" ".data.id")

# 验证用户能看到自己的+系统的分类
RESP=$(get "/api/v1/category" "$USER_TOKEN")
USER_CAT_COUNT=$(jval "$RESP" ".data.list | length")
if [ "$USER_CAT_COUNT" -ge "5" ] 2>/dev/null; then
    pass "普通用户分类列表: ${USER_CAT_COUNT}个(含系统分类)"
else
    fail "普通用户分类列表: ${USER_CAT_COUNT}个"
fi

# ========== 8. 任务 CRUD ==========
section "8. 任务 CRUD"

# 创建任务1
RESP=$(post "/api/v1/task" "{\"title\":\"完成项目文档\",\"content\":\"编写README\",\"priority\":1,\"categoryId\":1}" "$ADMIN_TOKEN")
TASK1_ID=$(jval "$RESP" ".data.id")
assert_not_empty "$RESP" "data.id" "创建任务1"

# 创建任务2
RESP=$(post "/api/v1/task" "{\"title\":\"修复BUG\",\"content\":\"修复登录问题\",\"priority\":2,\"categoryId\":2}" "$ADMIN_TOKEN")
TASK2_ID=$(jval "$RESP" ".data.id")
assert_not_empty "$RESP" "data.id" "创建任务2"

# 创建任务 — 空标题
RESP=$(post "/api/v1/task" "{\"title\":\"\",\"content\":\"test\"}" "$ADMIN_TOKEN")
assert_err "$RESP" "创建任务空标题拒绝"

# 创建任务 — 标题超长
RESP=$(post "/api/v1/task" "{\"title\":\"$(python3 -c 'print("A"*101)')\"}" "$ADMIN_TOKEN")
assert_err "$RESP" "创建任务标题超长拒绝"

# 任务列表
RESP=$(get "/api/v1/task?page=1&pageSize=10" "$ADMIN_TOKEN")
assert_eq "$RESP" "data.total" "2" "任务列表(2条)"

# 任务详情
RESP=$(get "/api/v1/task/${TASK1_ID}" "$ADMIN_TOKEN")
assert_eq "$RESP" "data.title" "完成项目文档" "任务详情"

# 更新任务
RESP=$(put "/api/v1/task/${TASK1_ID}" "{\"title\":\"完成项目文档V2\",\"content\":\"更新README\",\"priority\":2,\"categoryId\":1}" "$ADMIN_TOKEN")
assert_ok "$RESP" "更新任务"

# 切换任务状态
RESP=$(patch_req "/api/v1/task/${TASK1_ID}/toggle" "$ADMIN_TOKEN")
assert_ok "$RESP" "切换任务状态"

# 批量操作
RESP=$(post "/api/v1/task/batch" "{\"ids\":[${TASK1_ID},${TASK2_ID}],\"action\":\"complete\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "批量完成"

# 批量操作 — 非法action
RESP=$(post "/api/v1/task/batch" "{\"ids\":[${TASK1_ID}],\"action\":\"invalid\"}" "$ADMIN_TOKEN")
assert_err "$RESP" "批量操作非法action拒绝"

# 批量操作 — 空ids
RESP=$(post "/api/v1/task/batch" "{\"ids\":[],\"action\":\"complete\"}" "$ADMIN_TOKEN")
assert_err "$RESP" "批量操作空ids拒绝"

# 删除任务
RESP=$(del "/api/v1/task/${TASK2_ID}" "$ADMIN_TOKEN")
assert_ok "$RESP" "删除任务2"

RESP=$(get "/api/v1/task?page=1&pageSize=10" "$ADMIN_TOKEN")
assert_eq "$RESP" "data.total" "1" "删除后任务列表(1条)"

# 8.2 任务筛选
# 创建多条任务用于筛选
RESP=$(post "/api/v1/task" "{\"title\":\"待办任务\",\"priority\":3}" "$ADMIN_TOKEN")
TASK3_ID=$(jval "$RESP" ".data.id")

# 状态筛选 — 待办(status=0)
RESP=$(get "/api/v1/task?page=1&pageSize=10&status=0" "$ADMIN_TOKEN")
TODO_COUNT=$(jval "$RESP" ".data.total")
if [ -n "$TODO_COUNT" ] && [ "$TODO_COUNT" != "null" ]; then
    pass "状态筛选待办: total=${TODO_COUNT}"
else
    fail "状态筛选待办: 异常"
fi

# 状态筛选 — 已完成(status=2)
RESP=$(get "/api/v1/task?page=1&pageSize=10&status=2" "$ADMIN_TOKEN")
DONE_COUNT=$(jval "$RESP" ".data.total")
if [ -n "$DONE_COUNT" ] && [ "$DONE_COUNT" != "null" ]; then
    pass "状态筛选已完成: total=${DONE_COUNT}"
else
    fail "状态筛选已完成: 异常"
fi

# 优先级筛选
RESP=$(get "/api/v1/task?page=1&pageSize=10&priority=3" "$ADMIN_TOKEN")
P3_COUNT=$(jval "$RESP" ".data.total")
if [ -n "$P3_COUNT" ] && [ "$P3_COUNT" != "null" ]; then
    pass "优先级筛选普通: total=${P3_COUNT}"
else
    fail "优先级筛选普通: 异常"
fi

# 8.3 跨用户操作 — 不能操作别人的任务
RESP=$(del "/api/v1/task/${TASK1_ID}" "$USER_TOKEN")
assert_err "$RESP" "跨用户删除任务拒绝"

RESP=$(put "/api/v1/task/${TASK1_ID}" "{\"title\":\"hack\"}" "$USER_TOKEN")
assert_err "$RESP" "跨用户修改任务拒绝"

# 8.4 不存在的任务
RESP=$(get "/api/v1/task/99999" "$ADMIN_TOKEN")
assert_err "$RESP" "查询不存在的任务"

# ========== 9. 统计 ==========
section "9. 统计"
RESP=$(get "/api/v1/stat" "$ADMIN_TOKEN")
TOTAL=$(jval "$RESP" ".data.total")
if [ -n "$TOTAL" ] && [ "$TOTAL" != "null" ]; then
    pass "统计: total=${TOTAL}, done=$(jval "$RESP" ".data.done"), todo=$(jval "$RESP" ".data.todo")"
else
    fail "统计: 异常 → $RESP"
fi

# ========== 10. 管理员功能 ==========
section "10. 管理员功能"

# 10.1 用户列表
RESP=$(get "/api/v1/admin/user?page=1&pageSize=10" "$ADMIN_TOKEN")
assert_eq "$RESP" "data.total" "2" "用户列表(2个)"

# 10.2 重置用户密码
RESP=$(put "/api/v1/admin/user/${TEST_USER_ID}/password" "{\"newPassword\":\"${TEST_USER_NEW_PASS}\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "重置用户密码"

# 重置密码弱密码
RESP=$(put "/api/v1/admin/user/${TEST_USER_ID}/password" "{\"newPassword\":\"${WEAK_PASS}\"}" "$ADMIN_TOKEN")
assert_err_msg "$RESP" "字母和数字" "重置弱密码拒绝"

# 用新密码登录验证
RESP=$(post "/api/v1/user/login" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_NEW_PASS}\"}")
if [ -n "$(jval "$RESP" ".data.token")" ] && [ "$(jval "$RESP" ".data.token")" != "null" ]; then
    pass "用户新密码登录: 成功"
    USER_TOKEN=$(jval "$RESP" ".data.token")
else
    fail "用户新密码登录: 失败"
fi

# 10.3 禁用/启用用户
RESP=$(patch_req "/api/v1/admin/user/${TEST_USER_ID}/toggle" "$ADMIN_TOKEN")
assert_ok "$RESP" "禁用用户"

RESP=$(post "/api/v1/user/login" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_NEW_PASS}\"}")
assert_err "$RESP" "被禁用用户登录"

# 重新启用
RESP=$(patch_req "/api/v1/admin/user/${TEST_USER_ID}/toggle" "$ADMIN_TOKEN")
assert_ok "$RESP" "启用用户"

RESP=$(post "/api/v1/user/login" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_NEW_PASS}\"}")
if [ -n "$(jval "$RESP" ".data.token")" ] && [ "$(jval "$RESP" ".data.token")" != "null" ]; then
    pass "启用后用户登录: 成功"
    USER_TOKEN=$(jval "$RESP" ".data.token")
else
    fail "启用后用户登录: 失败"
fi

# 10.4 普通用户不能访问管理接口
RESP=$(get "/api/v1/admin/user?page=1&pageSize=10" "$USER_TOKEN")
assert_err "$RESP" "普通用户访问管理接口拒绝"

RESP=$(get "/api/v1/admin/config" "$USER_TOKEN")
assert_err "$RESP" "普通用户访问配置拒绝"

# 10.5 删除用户
RESP=$(post "/api/v1/category" "{\"name\":\"temp\"}" "$USER_TOKEN")
TEMP_USER_CAT_ID=$(jval "$RESP" ".data.id")

RESP=$(del "/api/v1/admin/user/${TEST_USER_ID}" "$ADMIN_TOKEN")
assert_ok "$RESP" "删除普通用户"

# 删除后不能登录
RESP=$(post "/api/v1/user/login" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_NEW_PASS}\"}")
assert_err "$RESP" "已删除用户登录拒绝"

# 不能删除管理员
RESP=$(get "/api/v1/admin/user?page=1&pageSize=10" "$ADMIN_TOKEN")
ADMIN_USER_ID=$(echo "$RESP" | jq -r '.data.list[] | select(.isAdmin==1) | .id' 2>/dev/null)
RESP=$(del "/api/v1/admin/user/${ADMIN_USER_ID}" "$ADMIN_TOKEN")
assert_err "$RESP" "删除管理员拒绝"

# 10.6 系统配置
RESP=$(get "/api/v1/admin/config" "$ADMIN_TOKEN")
assert_eq "$RESP" "data.list | length" "7" "系统配置(7项)"

RESP=$(put "/api/v1/admin/config" "{\"key\":\"allow_register\",\"value\":\"false\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "关闭注册"

RESP=$(get "/api/v1/user/check-register")
assert_eq "$RESP" "data.allowRegister" "false" "验证注册已关闭"

# 配置验证 — 空key
RESP=$(put "/api/v1/admin/config" "{\"key\":\"\",\"value\":\"test\"}" "$ADMIN_TOKEN")
assert_err "$RESP" "配置空key拒绝"

# 配置验证 — 空value
RESP=$(put "/api/v1/admin/config" "{\"key\":\"test_key\",\"value\":\"\"}" "$ADMIN_TOKEN")
assert_err "$RESP" "配置空value拒绝"

# ========== 11. 日志 ==========
section "11. 日志"

RESP=$(get "/api/v1/admin/log/operation?page=1&pageSize=5" "$ADMIN_TOKEN")
TOTAL=$(jval "$RESP" ".data.total")
if [ -n "$TOTAL" ] && [ "$TOTAL" != "null" ]; then
    pass "操作日志: total=${TOTAL}"
else
    fail "操作日志: 异常"
fi

RESP=$(get "/api/v1/admin/log/login?page=1&pageSize=5" "$ADMIN_TOKEN")
TOTAL=$(jval "$RESP" ".data.total")
if [ -n "$TOTAL" ] && [ "$TOTAL" != "null" ] && [ "$TOTAL" -gt "0" ] 2>/dev/null; then
    pass "登录日志: total=${TOTAL}"
else
    fail "登录日志: 异常"
fi

# ========== 12. 登录限流 ==========
section "12. 登录限流"

# 连续发送错误登录请求，触发限流（最多10次/15分钟）
RATE_LIMIT_TRIGGERED=false
for i in $(seq 1 15); do
    HTTP_CODE=$(curl -s -o /tmp/ratelimit_resp.json -w "%{http_code}" \
        -X POST "${BASE_URL}/api/v1/user/login" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"ratelimit_test\",\"password\":\"wrong${i}\"}")
    if [ "$HTTP_CODE" = "429" ]; then
        RATE_LIMIT_TRIGGERED=true
        RESP=$(cat /tmp/ratelimit_resp.json)
        RATE_CODE=$(jval "$RESP" ".code")
        if [ "$RATE_CODE" = "42901" ]; then
            pass "登录限流: 第${i}次触发429, code=${RATE_CODE}"
        else
            fail "登录限流: HTTP 429 但 code=${RATE_CODE}"
        fi
        # 验证 Retry-After 头
        RETRY_AFTER=$(curl -s -D - -o /dev/null -X POST "${BASE_URL}/api/v1/user/login" \
            -H "Content-Type: application/json" \
            -d "{\"username\":\"ratelimit_test\",\"password\":\"wrong\"}" 2>/dev/null | grep -i "retry-after" | tr -d '\r')
        if [ -n "$RETRY_AFTER" ]; then
            pass "登录限流: Retry-After 头存在"
        else
            fail "登录限流: Retry-After 头缺失"
        fi
        break
    fi
done

if [ "$RATE_LIMIT_TRIGGERED" = "false" ]; then
    fail "登录限流: 8次错误登录未触发限流"
fi

# ========== 13. 安全头 ==========
section "13. 安全头"

HEADERS=$(curl -s -D - -o /dev/null "${BASE_URL}/api/v1/user/check-register" 2>/dev/null)

for header in "X-Content-Type-Options" "X-Frame-Options" "X-XSS-Protection" "Referrer-Policy"; do
    if echo "$HEADERS" | grep -qi "$header"; then
        pass "安全头: ${header} 存在"
    else
        fail "安全头: ${header} 缺失"
    fi
done

# ========== 汇总 ==========
section "测试汇总"
TOTAL_TESTS=$((PASS + FAIL))
printf "  总计: %d  ${GREEN}通过: %d${NC}  ${RED}失败: %d${NC}\n" "$TOTAL_TESTS" "$PASS" "$FAIL"

if [ "$FAIL" -gt 0 ]; then
    printf "\n${RED}存在失败项，请检查！${NC}\n"
    exit 1
else
    printf "\n${GREEN}全部测试通过！${NC}\n"
    exit 0
fi
