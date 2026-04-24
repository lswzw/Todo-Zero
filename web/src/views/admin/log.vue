<template>
  <div class="admin-card">
    <div class="card-header">
      <h2>操作日志</h2>
      <div class="header-actions">
        <el-input v-model="keyword" placeholder="搜索用户名" clearable style="width: 180px" @clear="loadLogs" @keyup.enter="loadLogs">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="action" placeholder="操作类型" clearable style="width: 140px" @change="loadLogs">
          <el-option label="创建" value="create" />
          <el-option label="更新" value="update" />
          <el-option label="删除" value="delete" />
          <el-option label="配置变更" value="config" />
        </el-select>
      </div>
    </div>

    <el-table :data="logs" stripe style="width: 100%">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="username" label="操作用户" width="120" />
      <el-table-column label="操作类型" width="110">
        <template #default="{ row }">
          <el-tag :type="actionTagType(row.action)" size="small">{{ row.action }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="targetType" label="对象类型" width="100" />
      <el-table-column prop="detail" label="操作详情" show-overflow-tooltip />
      <el-table-column prop="ip" label="IP地址" width="140" />
      <el-table-column prop="createTime" label="操作时间" width="180" />
    </el-table>

    <div class="pagination">
      <el-pagination v-model:current-page="page" :page-size="10" :total="total" layout="total, prev, pager, next" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { getOperationLogList } from '@/api'
import type { OperationLogItem } from '@/types'

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
    // 错误已由拦截器处理
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
