<template>
  <div class="admin-card">
    <div class="card-header">
      <h2>{{ t('log.operationLog') }}</h2>
      <div class="header-actions">
        <el-input
          v-model="keyword"
          :placeholder="t('log.searchUsername')"
          clearable
          style="width: 180px"
          @clear="loadLogs"
          @keyup.enter="loadLogs"
        >
          <template #prefix
            ><el-icon><Search /></el-icon
          ></template>
        </el-input>
        <el-select
          v-model="action"
          :placeholder="t('log.actionType')"
          clearable
          style="width: 140px"
          @change="loadLogs"
        >
          <el-option :label="t('log.create')" value="create" />
          <el-option :label="t('log.update')" value="update" />
          <el-option :label="t('log.delete')" value="delete" />
          <el-option :label="t('log.configChange')" value="config" />
        </el-select>
      </div>
    </div>

    <el-table :data="logs" stripe style="width: 100%">
      <el-table-column prop="id" :label="t('admin.id')" width="70" />
      <el-table-column prop="username" :label="t('log.opUser')" width="120" />
      <el-table-column :label="t('log.actionType')" width="110">
        <template #default="{ row }">
          <el-tag :type="actionTagType(row.action)" size="small">{{ row.action }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="targetType" :label="t('log.targetType')" width="100" />
      <el-table-column prop="detail" :label="t('log.detail')" show-overflow-tooltip />
      <el-table-column prop="ip" :label="t('log.ipAddress')" width="140" />
      <el-table-column prop="createTime" :label="t('log.opTime')" width="180" />
    </el-table>

    <div class="pagination">
      <el-pagination v-model:current-page="page" :page-size="10" :total="total" layout="total, prev, pager, next" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { getOperationLogList } from '@/api'
import type { OperationLogItem } from '@/types'

const { t } = useI18n()

const logs = ref<OperationLogItem[]>([])
const total = ref(0)
const page = ref(1)
const keyword = ref('')
const action = ref('')

onMounted(() => loadLogs())

watch(page, () => loadLogs())

async function loadLogs() {
  try {
    const params: Record<string, unknown> = { page: page.value, pageSize: 10 }
    if (keyword.value) params.username = keyword.value
    if (action.value) params.action = action.value
    const res = await getOperationLogList(params)
    logs.value = res.list || []
    total.value = res.total || 0
  } catch {
    ElMessage.error(t('log.loadOperationLogFailed'))
  }
}

function actionTagType(a: string) {
  if (a.includes('create') || a === 'register') return 'success'
  if (a.includes('update')) return 'warning'
  if (a.includes('delete')) return 'danger'
  return 'info'
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

.header-actions {
  display: flex;
  gap: 8px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
