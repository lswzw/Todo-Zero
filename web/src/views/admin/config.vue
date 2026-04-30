<template>
  <div class="admin-card">
    <div class="card-header">
      <h2>{{ t('config.systemSettings') }}</h2>
    </div>

    <div v-for="item in configs" :key="item.key" class="config-item">
      <div class="config-info">
        <div class="config-title">{{ configMeta[item.key]?.title || item.key }}</div>
        <div class="config-desc">{{ configMeta[item.key]?.desc || item.remark }}</div>
      </div>
      <div class="config-control">
        <template v-if="item.key === 'allow_register'">
          <el-switch
            :model-value="item._value === 'true'"
            :active-text="t('config.open')"
            :inactive-text="t('config.close')"
            @change="(val: boolean) => handleUpdate(item.key, String(val))"
          />
        </template>
        <template v-else-if="item.key === 'db_backup_enabled'">
          <el-switch
            :model-value="item._value === '1'"
            :active-text="t('config.open')"
            :inactive-text="t('config.close')"
            @change="(val: boolean) => handleUpdate(item.key, val ? '1' : '0')"
          />
        </template>
        <template v-else>
          <el-input v-model="item._value" style="width: 200px" />
          <el-button type="primary" size="small" @click="handleUpdate(item.key, item._value)">{{
            t('common.save')
          }}</el-button>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getConfigList, updateConfig } from '@/api'
import type { ConfigItem } from '@/types'

const { t, tm, rt } = useI18n()

const configs = ref<ConfigItem[]>([])

const configMeta = computed(() => ({
  allow_register: {
    title: rt(tm('config.allowRegister.title') as string),
    desc: rt(tm('config.allowRegister.desc') as string),
  },
  site_name: {
    title: rt(tm('config.siteName.title') as string),
    desc: rt(tm('config.siteName.desc') as string),
  },
  task_auto_delete_days: {
    title: rt(tm('config.taskAutoDeleteDays.title') as string),
    desc: rt(tm('config.taskAutoDeleteDays.desc') as string),
  },
  task_trash_retention_days: {
    title: rt(tm('config.taskTrashRetentionDays.title') as string),
    desc: rt(tm('config.taskTrashRetentionDays.desc') as string),
  },
  log_auto_delete_days: {
    title: rt(tm('config.logAutoDeleteDays.title') as string),
    desc: rt(tm('config.logAutoDeleteDays.desc') as string),
  },
  db_backup_enabled: {
    title: rt(tm('config.dbBackupEnabled.title') as string),
    desc: rt(tm('config.dbBackupEnabled.desc') as string),
  },
  db_backup_interval_hours: {
    title: rt(tm('config.dbBackupIntervalHours.title') as string),
    desc: rt(tm('config.dbBackupIntervalHours.desc') as string),
  },
  db_backup_max_count: {
    title: rt(tm('config.dbBackupMaxCount.title') as string),
    desc: rt(tm('config.dbBackupMaxCount.desc') as string),
  },
}))

onMounted(() => loadConfigs())

async function loadConfigs() {
  try {
    const res = await getConfigList()
    configs.value = (res.list || []).map((item: ConfigItem) => ({ ...item, _value: item.value }))
  } catch {
    ElMessage.error(t('config.loadConfigFailed'))
  }
}

async function handleUpdate(key: string, value: string) {
  try {
    await updateConfig({ key, value })
    ElMessage.success(t('config.configSaved', { key }))
    loadConfigs()
  } catch {
    ElMessage.error(t('config.saveConfigFailed'))
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
  margin-bottom: 20px;
}

.card-header h2 {
  font-size: 18px;
  background: linear-gradient(135deg, #303133, #667eea);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.config-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.config-item:last-child {
  border-bottom: none;
}

.config-title {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.config-desc {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.config-control {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
