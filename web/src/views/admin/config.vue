<template>
  <div class="admin-card">
    <div class="card-header">
      <h2>系统设置</h2>
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
            active-text="开启"
            inactive-text="关闭"
            @change="(val: boolean) => handleUpdate(item.key, String(val))"
          />
        </template>
        <template v-else-if="item.key === 'db_backup_enabled'">
          <el-switch
            :model-value="item._value === '1'"
            active-text="开启"
            inactive-text="关闭"
            @change="(val: boolean) => handleUpdate(item.key, val ? '1' : '0')"
          />
        </template>
        <template v-else>
          <el-input v-model="item._value" style="width: 200px" />
          <el-button type="primary" size="small" @click="handleUpdate(item.key, item._value)">保存</el-button>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getConfigList, updateConfig } from '@/api'
import type { ConfigItem } from '@/types'

const configs = ref<ConfigItem[]>([])

const configMeta: Record<string, { title: string; desc: string }> = {
  allow_register: {
    title: '开放注册',
    desc: '开启后允许新用户自行注册账号，关闭后仅管理员可创建账号',
  },
  site_name: {
    title: '站点名称',
    desc: '显示在浏览器标题栏和登录页面',
  },
  task_auto_delete_days: {
    title: '自动清理已完成任务',
    desc: '已完成任务超过指定天数后自动永久删除（0=不清理）',
  },
  task_trash_retention_days: {
    title: '回收站保留天数',
    desc: '手动删除的任务在回收站中保留的天数，超过后永久删除（0=不清理，默认30天）',
  },
  log_auto_delete_days: {
    title: '自动清理日志',
    desc: '操作日志和登录日志超过指定天数后自动删除（0=不清理）',
  },
  db_backup_enabled: {
    title: '数据库自动备份',
    desc: '开启后系统将按设定间隔自动备份SQLite数据库文件（0=关闭 1=开启）',
  },
  db_backup_interval_hours: {
    title: '备份间隔（小时）',
    desc: '自动备份的时间间隔，单位为小时（默认24小时）',
  },
  db_backup_max_count: {
    title: '最大备份数量',
    desc: '保留的最大备份数量，超过后自动删除最旧的备份（默认7份）',
  },
}

onMounted(() => loadConfigs())

async function loadConfigs() {
  try {
    const res = await getConfigList()
    configs.value = (res.list || []).map((item: ConfigItem) => ({ ...item, _value: item.value }))
  } catch {
    ElMessage.error('加载配置失败')
  }
}

async function handleUpdate(key: string, value: string) {
  try {
    await updateConfig({ key, value })
    ElMessage.success(`配置已保存：${key}`)
    loadConfigs()
  } catch {
    ElMessage.error('保存配置失败')
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
  margin-bottom: 20px;
}

.card-header h2 {
  font-size: 18px;
  color: #303133;
}

.config-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 0;
  border-bottom: 1px solid #f0f0f0;
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
