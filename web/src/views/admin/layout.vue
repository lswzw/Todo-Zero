<template>
  <div class="admin-layout">
    <!-- 顶部导航 -->
    <header class="admin-navbar">
      <div class="nav-left">
        <span class="logo-icon">📝</span>
        <span class="logo-text">Todo App {{ t('admin.adminPanel') }}</span>
      </div>
      <div class="nav-right">
        <el-select v-model="currentLang" size="small" style="width: 90px" @change="handleLocaleChange">
          <el-option v-for="opt in localeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
        </el-select>
        <span>{{ t('admin.admin') }}：{{ userStore.username }}</span>
        <el-button text @click="router.push('/')">{{ t('admin.homePage') }}</el-button>
        <el-button text type="danger" @click="handleLogout">{{ t('common.logout') }}</el-button>
      </div>
    </header>

    <div class="admin-body">
      <!-- 侧边栏 -->
      <aside class="admin-sidebar">
        <router-link to="/admin/user" class="sidebar-item" active-class="active">
          <span>👤</span> {{ t('admin.userManagement') }}
        </router-link>
        <router-link to="/admin/config" class="sidebar-item" active-class="active">
          <span>⚙️</span> {{ t('admin.systemSettings') }}
        </router-link>
        <router-link to="/admin/log" class="sidebar-item" active-class="active">
          <span>📋</span> {{ t('admin.operationLog') }}
        </router-link>
        <router-link to="/admin/login-log" class="sidebar-item" active-class="active">
          <span>🔑</span> {{ t('admin.loginLog') }}
        </router-link>
        <router-link to="/admin/backup" class="sidebar-item" active-class="active">
          <span>💾</span> {{ t('admin.dbBackup') }}
        </router-link>
      </aside>

      <!-- 内容区 -->
      <main class="admin-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { resetAuthVerified } from '@/router'
import { useLocale } from '@/composables/useLocale'

const { t } = useI18n()
const { currentLocale, setLocale, localeOptions } = useLocale()
const currentLang = ref(currentLocale.value)

const router = useRouter()
const userStore = useUserStore()

function handleLocaleChange(lang: string) {
  setLocale(lang)
}

function handleLogout() {
  ElMessageBox.confirm(t('common.logoutConfirm'), t('common.tip'), { type: 'warning' })
    .then(() => {
      resetAuthVerified()
      userStore.logout()
      router.push('/login')
      ElMessage.success(t('common.logoutSuccess'))
    })
    .catch(() => {})
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.admin-navbar {
  position: sticky;
  top: 0;
  z-index: var(--z-sticky);
  height: 60px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-icon {
  font-size: 24px;
}
.logo-text {
  font-size: 18px;
  font-weight: 600;
  background: linear-gradient(135deg, #303133, #667eea);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  color: #606266;
}

.admin-body {
  display: flex;
  padding: 0;
}

.admin-sidebar {
  width: 200px;
  min-height: calc(100vh - 60px);
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.3);
  padding: 16px 0;
  flex-shrink: 0;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  color: #606266;
  font-size: 14px;
  text-decoration: none;
  transition: all 0.2s;
}

.sidebar-item:hover {
  background: rgba(102, 126, 234, 0.08);
  color: #667eea;
}

.sidebar-item.active {
  background: rgba(102, 126, 234, 0.12);
  color: #667eea;
  font-weight: 500;
  border-right: 3px solid #667eea;
}

.admin-content {
  flex: 1;
  padding: 24px;
  min-width: 0;
}
</style>
