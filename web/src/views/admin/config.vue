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
