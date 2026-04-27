<template>
  <div class="admin-card">
    <div class="card-header">
      <h2>用户管理</h2>
      <el-input v-model="keyword" placeholder="搜索用户名" clearable style="width: 240px" @clear="loadUsers" @keyup.enter="loadUsers">
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
    </div>

    <el-table :data="users" stripe style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column label="角色" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.isAdmin === 1" type="danger" size="small">管理员</el-tag>
          <span v-else style="color: #909399; font-size: 13px">普通用户</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.status === 1" type="success" size="small">启用</el-tag>
          <el-tag v-else type="info" size="small">禁用</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="注册时间" width="180" />
      <el-table-column label="操作" width="260">
        <template #default="{ row }">
          <el-button size="small" @click="openResetDialog(row)">重置密码</el-button>
          <el-button size="small" :type="row.status === 1 ? 'warning' : 'success'" :disabled="row.isAdmin === 1" @click="handleToggleStatus(row)">
            {{ row.status === 1 ? '禁用' : '启用' }}
          </el-button>
          <el-popconfirm title="确定删除该用户？" @confirm="handleDelete(row)">
            <template #reference>
              <el-button size="small" type="danger" :disabled="row.isAdmin === 1">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination v-model:current-page="page" :page-size="10" :total="total" layout="total, prev, pager, next" />
    </div>

    <!-- 重置密码弹窗 -->
    <el-dialog v-model="resetDialogVisible" title="重置密码" width="420px" destroy-on-close>
      <el-form ref="resetFormRef" :model="resetForm" :rules="resetRules" label-width="80px">
        <el-form-item label="用户名">
          <el-input :model-value="resetForm.username" disabled />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="resetForm.newPassword" type="password" show-password placeholder="请输入新密码(6-20位，需含字母和数字)" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="resetLoading" @click="handleResetPassword">确定重置</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getUserList, resetPassword, toggleUserStatus, deleteUser } from '@/api'
import type { UserListItem } from '@/types'

const users = ref<UserListItem[]>([])
const total = ref(0)
const page = ref(1)
const keyword = ref('')

const resetDialogVisible = ref(false)
const resetLoading = ref(false)
const resetFormRef = ref<FormInstance>()
const resetForm = ref({ id: 0, username: '', newPassword: '' })
const resetRules = {
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度6-20位', trigger: 'blur' },
    { pattern: /^(?=.*[a-zA-Z])(?=.*\d)/, message: '密码必须包含字母和数字', trigger: 'blur' },
  ],
}

onMounted(() => loadUsers())

watch(page, () => loadUsers())

async function loadUsers() {
  try {
    const params: Record<string, unknown> = { page: page.value, pageSize: 10 }
    if (keyword.value) params.keyword = keyword.value
    const res = await getUserList(params)
    users.value = res.list || []
    total.value = res.total || 0
  } catch {
    ElMessage.error('加载用户列表失败')
  }
}

function openResetDialog(row: UserListItem) {
  resetForm.value = { id: row.id, username: row.username, newPassword: '' }
  resetDialogVisible.value = true
}

async function handleResetPassword() {
  await resetFormRef.value?.validate()
  resetLoading.value = true
  try {
    await resetPassword(resetForm.value.id, { newPassword: resetForm.value.newPassword })
    ElMessage.success('密码重置成功')
    resetDialogVisible.value = false
    loadUsers()
  } catch {
    ElMessage.error('重置密码失败')
  } finally {
    resetLoading.value = false
  }
}

async function handleToggleStatus(row: UserListItem) {
  try {
    await toggleUserStatus(row.id)
    ElMessage.success(row.status === 1 ? '已禁用' : '已启用')
    loadUsers()
  } catch {
    ElMessage.error('切换用户状态失败')
  }
}

async function handleDelete(row: UserListItem) {
  try {
    await deleteUser(row.id)
    ElMessage.success('已删除')
    loadUsers()
  } catch {
    ElMessage.error('删除用户失败')
  }
}
</script>

<style scoped>
.admin-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.card-header h2 {
  font-size: 18px;
  color: #303133;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
