<template>
  <div class="admin-card">
    <div class="card-header">
      <h2>{{ t('admin.userManagement') }}</h2>
      <el-input
        v-model="keyword"
        :placeholder="t('admin.searchUsername')"
        clearable
        style="width: 240px"
        @clear="loadUsers"
        @keyup.enter="loadUsers"
      >
        <template #prefix
          ><el-icon><Search /></el-icon
        ></template>
      </el-input>
    </div>

    <el-table :data="users" stripe style="width: 100%">
      <el-table-column prop="id" :label="t('admin.id')" width="80" />
      <el-table-column prop="username" :label="t('admin.username')" />
      <el-table-column :label="t('admin.role')" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.isAdmin === 1" type="danger" size="small">{{ t('admin.adminRole') }}</el-tag>
          <span v-else style="color: #909399; font-size: 13px">{{ t('admin.normalUser') }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('admin.status')" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.status === 1" type="success" size="small">{{ t('admin.enabled') }}</el-tag>
          <el-tag v-else type="info" size="small">{{ t('admin.disabled') }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" :label="t('admin.registerTime')" width="180" />
      <el-table-column :label="t('admin.actions')" width="260">
        <template #default="{ row }">
          <el-button size="small" @click="openResetDialog(row)">{{ t('admin.resetPassword') }}</el-button>
          <el-button
            size="small"
            :type="row.status === 1 ? 'warning' : 'success'"
            :disabled="row.isAdmin === 1"
            @click="handleToggleStatus(row)"
          >
            {{ row.status === 1 ? t('admin.disabled') : t('admin.enabled') }}
          </el-button>
          <el-popconfirm :title="t('admin.deleteConfirm')" @confirm="handleDelete(row)">
            <template #reference>
              <el-button size="small" type="danger" :disabled="row.isAdmin === 1">{{ t('common.delete') }}</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination v-model:current-page="page" :page-size="10" :total="total" layout="total, prev, pager, next" />
    </div>

    <!-- 重置密码弹窗 -->
    <el-dialog v-model="resetDialogVisible" :title="t('admin.resetPasswordTitle')" width="420px" destroy-on-close>
      <el-form ref="resetFormRef" :model="resetForm" :rules="resetRules" label-width="80px">
        <el-form-item :label="t('admin.username')">
          <el-input :model-value="resetForm.username" disabled />
        </el-form-item>
        <el-form-item :label="t('auth.newPassword')" prop="newPassword">
          <el-input
            v-model="resetForm.newPassword"
            type="password"
            show-password
            :placeholder="t('admin.enterNewPassword')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="resetLoading" @click="handleResetPassword">{{
          t('admin.confirmReset')
        }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { getUserList, resetPassword, toggleUserStatus, deleteUser } from '@/api'
import type { UserListItem } from '@/types'

const { t } = useI18n()

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
    { required: true, message: () => t('auth.enterNewPassword'), trigger: 'blur' },
    { min: 6, max: 20, message: () => t('auth.passwordLength'), trigger: 'blur' },
    { pattern: /^(?=.*[a-zA-Z])(?=.*\d)/, message: () => t('auth.passwordComplexity'), trigger: 'blur' },
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
    ElMessage.error(t('admin.loadUsersFailed'))
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
    ElMessage.success(t('admin.passwordResetSuccess'))
    resetDialogVisible.value = false
    loadUsers()
  } catch {
    ElMessage.error(t('admin.resetPasswordFailed'))
  } finally {
    resetLoading.value = false
  }
}

async function handleToggleStatus(row: UserListItem) {
  try {
    await toggleUserStatus(row.id)
    ElMessage.success(row.status === 1 ? t('admin.disabledUser') : t('admin.enabledUser'))
    loadUsers()
  } catch {
    ElMessage.error(t('admin.toggleUserStatusFailed'))
  }
}

async function handleDelete(row: UserListItem) {
  try {
    await deleteUser(row.id)
    ElMessage.success(t('home.deleted'))
    loadUsers()
  } catch {
    ElMessage.error(t('admin.deleteUserFailed'))
  }
}
</script>

<style scoped>
.admin-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.card-header h2 {
  font-size: 18px;
  background: linear-gradient(135deg, #303133, #667eea);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
