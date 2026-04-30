<template>
  <div class="admin-card">
    <div class="card-header">
      <h2>{{ t('backup.dbBackup') }}</h2>
      <el-button type="primary" :loading="backupLoading" @click="handleBackup">{{ t('backup.backupNow') }}</el-button>
    </div>

    <!-- 备份状态 -->
    <div class="backup-status">
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item :label="t('backup.autoBackup')">
          {{ backupEnabled ? t('backup.enabled') : t('backup.notEnabled') }}
          <el-tag :type="backupEnabled ? 'success' : 'info'" size="small" style="margin-left: 8px">
            {{ backupEnabled ? 'ON' : 'OFF' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="t('backup.backupInterval')">
          {{ backupIntervalHours }} {{ t('backup.hours') }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('backup.maxBackups')">
          {{ backupMaxCount }} {{ t('backup.copies') }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('backup.currentBackups')">
          {{ backups.length }} {{ t('backup.copies') }}
        </el-descriptions-item>
      </el-descriptions>
      <div class="status-tip">
        {{ t('backup.configTip') }}
      </div>
    </div>

    <!-- 备份列表 -->
    <el-table :data="backups" stripe style="margin-top: 20px" :empty-text="t('backup.noBackups')">
      <el-table-column :label="t('backup.fileName')" prop="fileName" min-width="240" />
      <el-table-column :label="t('backup.size')" width="120">
        <template #default="{ row }">
          {{ formatFileSize(row.fileSize) }}
        </template>
      </el-table-column>
      <el-table-column :label="t('backup.createTime')" prop="createTime" width="180" />
      <el-table-column :label="t('admin.actions')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="handleDownload(row.fileName)">
            {{ t('backup.download') }}
          </el-button>
          <el-button type="warning" link size="small" @click="handleRestore(row.fileName)">
            {{ t('backup.restore') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 恢复确认对话框 -->
    <el-dialog
      v-model="restoreDialogVisible"
      :title="t('backup.restoreTitle')"
      width="480px"
      :close-on-click-modal="false"
    >
      <el-alert type="error" :closable="false" show-icon style="margin-bottom: 16px">
        <template #title>
          <strong>{{ t('backup.dangerousAction') }}</strong>
        </template>
        {{ t('backup.restoreWarning') }}
      </el-alert>
      <div style="margin-bottom: 12px">
        <p>{{ t('backup.restoreFile', { fileName: restoreFileName }) }}</p>
        <p style="margin-top: 8px; color: #e6a23c">{{ t('backup.safetyBackupTip') }}</p>
        <p style="margin-top: 8px; color: #f56c6c">{{ t('backup.reloginTip') }}</p>
      </div>
      <el-form @submit.prevent>
        <el-form-item :label="t('backup.confirmRestorePlaceholder')">
          <el-input v-model="restoreConfirmText" :placeholder="t('backup.confirmRestorePlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="restoreDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button
          type="danger"
          :loading="restoreLoading"
          :disabled="restoreConfirmText !== t('backup.confirmRestore')"
          @click="doRestore"
        >
          {{ t('backup.confirmRestoreBtn') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getBackupList, triggerBackup, downloadBackup, restoreBackup, getConfigList } from '@/api'
import type { BackupItem, ConfigItem } from '@/types'

const { t } = useI18n()

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
    ElMessage.error(t('backup.loadBackupFailed'))
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
    await ElMessageBox.confirm(t('backup.backupConfirm'), t('backup.manualBackupTitle'), { type: 'info' })
  } catch {
    return
  }

  backupLoading.value = true
  try {
    const res = await triggerBackup()
    ElMessage.success(t('backup.backupSuccess', { fileName: res.fileName, size: formatFileSize(res.fileSize) }))
    loadBackups()
  } catch {
    ElMessage.error(t('backup.backupFailed'))
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
      message: t('backup.restoreSuccess', { fileName: res.preRestoreBackup }),
      duration: 8000,
    })
    // Clear token and redirect to login after a short delay
    setTimeout(() => {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }, 2000)
  } catch {
    ElMessage.error(t('backup.restoreFailed'))
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

.backup-status {
  margin-top: 12px;
}

.status-tip {
  margin-top: 12px;
  font-size: 13px;
  color: #909399;
}
</style>
