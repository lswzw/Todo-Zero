<template>
  <div class="admin-card">
    <div class="card-header">
      <h2>{{ t('log.loginLog') }}</h2>
      <el-input
        v-model="keyword"
        :placeholder="t('log.searchUsername')"
        clearable
        style="width: 240px"
        @clear="loadLogs"
        @keyup.enter="loadLogs"
      >
        <template #prefix
          ><el-icon><Search /></el-icon
        ></template>
      </el-input>
    </div>

    <el-table :data="logs" stripe style="width: 100%">
      <el-table-column prop="id" :label="t('admin.id')" width="70" />
      <el-table-column prop="username" :label="t('admin.username')" width="130" />
      <el-table-column :label="t('admin.status')" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? t('log.success') : t('log.failed') }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="ip" :label="t('log.ipAddress')" width="150" />
      <el-table-column prop="remark" :label="t('log.remark')" width="150" show-overflow-tooltip />
      <el-table-column prop="createTime" :label="t('log.loginTime')" width="180" />
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
import { getLoginLogList } from '@/api'
import type { LoginLogItem } from '@/types'

const { t } = useI18n()

const logs = ref<LoginLogItem[]>([])
const total = ref(0)
const page = ref(1)
const keyword = ref('')

onMounted(() => loadLogs())

watch(page, () => loadLogs())

async function loadLogs() {
  try {
    const params: Record<string, unknown> = { page: page.value, pageSize: 10 }
    if (keyword.value) params.username = keyword.value
    const res = await getLoginLogList(params)
    logs.value = res.list || []
    total.value = res.total || 0
  } catch {
    ElMessage.error(t('log.loadLoginLogFailed'))
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
