<template>
  <div class="admin-card">
    <div class="card-header">
      <h2>数据库备份</h2>
      <el-button type="primary" :loading="backupLoading" @click="handleBackup">立即备份</el-button>
    </div>

    <!-- 备份状态 -->
    <div class="backup-status">
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="自动备份">
          {{ backupEnabled ? '已开启' : '未开启' }}
          <el-tag :type="backupEnabled ? 'success' : 'info'" size="small" style="margin-left: 8px">
            {{ backupEnabled ? 'ON' : 'OFF' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="备份间隔">
          {{ backupIntervalHours }} 小时
        </el-descriptions-item>
        <el-descriptions-item label="最大备份数">
          {{ backupMaxCount }} 份
        </el-descriptions-item>
        <el-descriptions-item label="当前备份数">
          {{ backups.length }} 份
        </el-descriptions-item>
      </el-descriptions>
      <div class="status-tip">
        提示：在「系统设置」中可开启自动备份并配置备份间隔和最大备份数
      </div>
    </div>

    <!-- 备份列表 -->
    <el-table :data="backups" stripe style="margin-top: 20px" empty-text="暂无备份">
      <el-table-column label="文件名" prop="fileName" min-width="240" />
      <el-table-column label="大小" width="120">
        <template #default="{ row }">
          {{ formatFileSize(row.fileSize) }}
        </template>
      </el-table-column>
      <el-table-column label="创建时间" prop="createTime" width="180" />
      <el-table-column label="操作" width="160" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="handleDownload(row.fileName)">
            下载
          </el-button>
          <el-button type="warning" link size="small" @click="handleRestore(row.fileName)">
            恢复
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 恢复确认对话框 -->
    <el-dialog v-model="restoreDialogVisible" title="恢复数据库" width="480px" :close-on-click-modal="false">
      <el-alert type="error" :closable="false" show-icon style="margin-bottom: 16px">
        <template #title>
          <strong>危险操作！</strong>
        </template>
        恢复操作将用备份数据<strong>替换当前所有数据</strong>，此操作不可撤销。
      </el-alert>
      <div style="margin-bottom: 12px">
        <p>将恢复备份：<strong>{{ restoreFileName }}</strong></p>
        <p style="margin-top: 8px; color: #e6a23c">系统会在恢复前自动创建安全备份，以防需要回退。</p>
        <p style="margin-top: 8px; color: #f56c6c">恢复后所有用户需重新登录。</p>
      </div>
      <el-form @submit.prevent>
        <el-form-item label='请输入 "确认恢复" 以继续'>
          <el-input v-model="restoreConfirmText" placeholder="确认恢复" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="restoreDialogVisible = false">取消</el-button>
        <el-button type="danger" :loading="restoreLoading" :disabled="restoreConfirmText !== '确认恢复'" @click="doRestore">
          确认恢复
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBackupList, triggerBackup, downloadBackup, restoreBackup, getConfigList } from '@/api'
import type { BackupItem, ConfigItem } from '@/types'

const backups = ref<BackupItem[]>([])
const backupLoading = ref(false)
const backupEnabled = ref(false)
const backupIntervalHours = ref(24)
const backupMaxCount = ref(7)

const restoreDialogVisible = ref(false)
const restoreFileName = ref('')
const restoreConfirmText = ref('')
const restoreLoading = ref(false)

onMounted(() => {
  loadBackups()
  loadConfig()
})

async function loadBackups() {
  try {
    const res = await getBackupList()
    backups.value = res.list || []
  } catch {
    ElMessage.error('加载备份列表失败')
  }
}

async function loadConfig() {
  try {
    const res = await getConfigList()
    const configs = res.list || []
    const findVal = (key: string) => configs.find((c: ConfigItem) => c.key === key)?.value
    backupEnabled.value = findVal('db_backup_enabled') === '1'
    backupIntervalHours.value = Number(findVal('db_backup_interval_hours')) || 24
    backupMaxCount.value = Number(findVal('db_backup_max_count')) || 7
  } catch {
    // ignore
  }
}

async function handleBackup() {
  try {
    await ElMessageBox.confirm('确定要立即创建数据库备份吗？', '手动备份', { type: 'info' })
  } catch {
    return
  }

  backupLoading.value = true
  try {
    const res = await triggerBackup()
    ElMessage.success(`备份成功：${res.fileName}（${formatFileSize(res.fileSize)}）`)
    loadBackups()
  } catch {
    ElMessage.error('备份失败')
  } finally {
    backupLoading.value = false
  }
}

function handleDownload(fileName: string) {
  const url = downloadBackup(fileName)
  const token = localStorage.getItem('token')
  const link = document.createElement('a')
  link.href = `${url}?token=${token}`
  link.download = fileName
  link.click()
}

function handleRestore(fileName: string) {
  restoreFileName.value = fileName
  restoreConfirmText.value = ''
  restoreDialogVisible.value = true
}

async function doRestore() {
  restoreLoading.value = true
  try {
    const res = await restoreBackup(restoreFileName.value)
    restoreDialogVisible.value = false
    ElMessage.success({
      message: `数据库已恢复！安全备份：${res.preRestoreBackup}。请重新登录。`,
      duration: 8000,
    })
    // Clear token and redirect to login after a short delay
    setTimeout(() => {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }, 2000)
  } catch {
    ElMessage.error('恢复失败')
  } finally {
    restoreLoading.value = false
  }
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(i > 0 ? 1 : 0) + ' ' + units[i]
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

.backup-status {
  margin-top: 12px;
}

.status-tip {
  margin-top: 12px;
  font-size: 13px;
  color: #909399;
}
</style>
