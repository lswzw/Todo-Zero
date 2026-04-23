#!/bin/bash
# ============================================
# Todo App 集成测试脚本
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
ADMIN_PASS="admin123"
NEW_ADMIN_PASS="admin456"
TEST_USER="zhangsan"
TEST_USER_PASS="123456"
TEST_USER_NEW_PASS="654321"

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
    local path="$1" body="$2" token="$3"
    curl -s -X PUT "${BASE_URL}${path}" \
        -H "Authorization: Bearer ${token}" \
        -H "Content-Type: application/json" -d "$body"
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
    if [ -z "$code" ] || [ "$code" = "null" ]; then
        pass "${desc}: 成功"
    else
        fail "${desc}: code=${code} msg=$(jval "$resp" ".msg")"
    fi
}

assert_err() {
    local resp="$1" desc="$2"
    local code
    code=$(jval "$resp" ".code")
    if [ -n "$code" ] && [ "$code" != "null" ]; then
        pass "${desc}: 正确拒绝 code=${code}"
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
        if curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/api/v1/user/check-register" 2>/dev/null | grep -q "200"; then
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

# ========== 1. 静态文件 ==========
section "1. 静态文件服务"
assert_http "/" "200"
assert_http "/login" "200"

# ========== 2. 无需认证接口 ==========
section "2. 无需认证接口"

RESP=$(get "/api/v1/user/check-register")
assert_eq "$RESP" "allowRegister" "false" "检查注册开关"

RESP=$(post "/api/v1/user/register" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_PASS}\"}")
assert_err "$RESP" "注册(已关闭)"

RESP=$(post "/api/v1/user/login" "{\"username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_PASS}\"}")
ADMIN_TOKEN=$(jval "$RESP" ".token")
if [ -n "$ADMIN_TOKEN" ] && [ "$ADMIN_TOKEN" != "null" ]; then
    pass "管理员登录: 获取到 token"
else
    fail "管理员登录: 失败 → $RESP"
    printf "${RED}无法继续测试${NC}\n"
    exit 1
fi
assert_eq "$RESP" "isAdmin" "1" "管理员角色"

# ========== 3. 开启注册 + 注册用户 ==========
section "3. 开启注册 + 注册用户"

RESP=$(put "/api/v1/admin/config" "{\"key\":\"allow_register\",\"value\":\"true\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "开启注册"

RESP=$(post "/api/v1/user/register" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_PASS}\"}")
assert_not_empty "$RESP" "id" "注册用户(${TEST_USER})"
TEST_USER_ID=$(jval "$RESP" ".id")

RESP=$(post "/api/v1/user/login" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_PASS}\"}")
USER_TOKEN=$(jval "$RESP" ".token")
if [ -n "$USER_TOKEN" ] && [ "$USER_TOKEN" != "null" ]; then
    pass "普通用户登录: 获取到 token"
else
    fail "普通用户登录: 失败 → $RESP"
fi

# ========== 4. 用户信息 ==========
section "4. 用户信息"
RESP=$(get "/api/v1/user/info" "$ADMIN_TOKEN")
assert_eq "$RESP" "username" "admin" "获取用户信息"
assert_eq "$RESP" "isAdmin" "1" "用户角色"

# ========== 5. 修改密码 ==========
section "5. 修改密码"

RESP=$(put "/api/v1/user/password" "{\"oldPassword\":\"${ADMIN_PASS}\",\"newPassword\":\"${NEW_ADMIN_PASS}\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "修改密码"

RESP=$(post "/api/v1/user/login" "{\"username\":\"${ADMIN_USER}\",\"password\":\"${NEW_ADMIN_PASS}\"}")
NEW_TOKEN=$(jval "$RESP" ".token")
if [ -n "$NEW_TOKEN" ] && [ "$NEW_TOKEN" != "null" ]; then
    pass "新密码登录: 成功"
    ADMIN_TOKEN="$NEW_TOKEN"
else
    fail "新密码登录: 失败 → $RESP"
fi

RESP=$(post "/api/v1/user/login" "{\"username\":\"${ADMIN_USER}\",\"password\":\"${ADMIN_PASS}\"}")
assert_err "$RESP" "旧密码登录"

# ========== 6. 分类 ==========
section "6. 分类"

RESP=$(get "/api/v1/category" "$ADMIN_TOKEN")
assert_eq "$RESP" "list | length" "4" "分类列表(4个初始分类)"

RESP=$(post "/api/v1/category" "{\"name\":\"娱乐\"}" "$ADMIN_TOKEN")
assert_not_empty "$RESP" "id" "创建分类"

RESP=$(get "/api/v1/category" "$ADMIN_TOKEN")
assert_eq "$RESP" "list | length" "5" "分类列表(新增后5个)"

# ========== 7. 任务 CRUD ==========
section "7. 任务 CRUD"

RESP=$(post "/api/v1/task" "{\"title\":\"完成项目文档\",\"content\":\"编写README\",\"priority\":1,\"status\":1,\"categoryId\":1}" "$ADMIN_TOKEN")
TASK1_ID=$(jval "$RESP" ".id")
assert_not_empty "$RESP" "id" "创建任务1"

RESP=$(post "/api/v1/task" "{\"title\":\"修复BUG\",\"content\":\"修复登录问题\",\"priority\":2,\"status\":1,\"categoryId\":2}" "$ADMIN_TOKEN")
TASK2_ID=$(jval "$RESP" ".id")
assert_not_empty "$RESP" "id" "创建任务2"

RESP=$(get "/api/v1/task?page=1&pageSize=10" "$ADMIN_TOKEN")
assert_eq "$RESP" "total" "2" "任务列表"

RESP=$(get "/api/v1/task/${TASK1_ID}" "$ADMIN_TOKEN")
assert_eq "$RESP" "title" "完成项目文档" "任务详情"

RESP=$(put "/api/v1/task/${TASK1_ID}" "{\"title\":\"完成项目文档V2\",\"content\":\"更新README\",\"priority\":2,\"status\":2,\"categoryId\":1}" "$ADMIN_TOKEN")
assert_ok "$RESP" "更新任务"

RESP=$(patch_req "/api/v1/task/${TASK1_ID}/toggle" "$ADMIN_TOKEN")
assert_ok "$RESP" "切换任务状态"

RESP=$(post "/api/v1/task/batch" "{\"ids\":[${TASK1_ID},${TASK2_ID}],\"action\":\"complete\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "批量操作"

RESP=$(del "/api/v1/task/${TASK2_ID}" "$ADMIN_TOKEN")
assert_ok "$RESP" "删除任务2"

RESP=$(get "/api/v1/task?page=1&pageSize=10" "$ADMIN_TOKEN")
assert_eq "$RESP" "total" "1" "删除后任务列表"

# ========== 8. 统计 ==========
section "8. 统计"
RESP=$(get "/api/v1/stat" "$ADMIN_TOKEN")
TOTAL=$(jval "$RESP" ".total")
if [ -n "$TOTAL" ] && [ "$TOTAL" != "null" ]; then
    pass "统计: total=${TOTAL}, done=$(jval "$RESP" ".done"), todo=$(jval "$RESP" ".todo")"
else
    fail "统计: 异常 → $RESP"
fi

# ========== 9. 管理员功能 ==========
section "9. 管理员功能"

RESP=$(get "/api/v1/admin/user?page=1&pageSize=10" "$ADMIN_TOKEN")
assert_eq "$RESP" "total" "2" "用户列表"

RESP=$(put "/api/v1/admin/user/${TEST_USER_ID}/password" "{\"newPassword\":\"${TEST_USER_NEW_PASS}\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "重置用户密码"

RESP=$(post "/api/v1/user/login" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_NEW_PASS}\"}")
if [ -n "$(jval "$RESP" ".token")" ] && [ "$(jval "$RESP" ".token")" != "null" ]; then
    pass "用户新密码登录: 成功"
else
    fail "用户新密码登录: 失败"
fi

RESP=$(patch_req "/api/v1/admin/user/${TEST_USER_ID}/toggle" "$ADMIN_TOKEN")
assert_ok "$RESP" "禁用用户"

RESP=$(post "/api/v1/user/login" "{\"username\":\"${TEST_USER}\",\"password\":\"${TEST_USER_NEW_PASS}\"}")
assert_err "$RESP" "被禁用用户登录"

RESP=$(get "/api/v1/admin/config" "$ADMIN_TOKEN")
assert_eq "$RESP" "list | length" "5" "系统配置(5项)"

RESP=$(put "/api/v1/admin/config" "{\"key\":\"allow_register\",\"value\":\"false\"}" "$ADMIN_TOKEN")
assert_ok "$RESP" "关闭注册"

RESP=$(get "/api/v1/user/check-register")
assert_eq "$RESP" "allowRegister" "false" "验证注册已关闭"

# ========== 10. 日志 ==========
section "10. 日志"

RESP=$(get "/api/v1/admin/log/operation?page=1&pageSize=5" "$ADMIN_TOKEN")
TOTAL=$(jval "$RESP" ".total")
if [ -n "$TOTAL" ] && [ "$TOTAL" != "null" ]; then
    pass "操作日志: total=${TOTAL}"
else
    fail "操作日志: 异常"
fi

RESP=$(get "/api/v1/admin/log/login?page=1&pageSize=5" "$ADMIN_TOKEN")
TOTAL=$(jval "$RESP" ".total")
if [ -n "$TOTAL" ] && [ "$TOTAL" != "null" ] && [ "$TOTAL" -gt "0" ] 2>/dev/null; then
    pass "登录日志: total=${TOTAL}"
else
    fail "登录日志: 异常"
fi

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
