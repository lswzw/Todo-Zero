<template>
  <GradientBg variant="primary" :show-decor="true">
    <div class="login-page">
      <div class="login-card">
        <div class="logo-area">
          <span class="logo-icon">📝</span>
          <h1 class="gradient-text">Todo App</h1>
          <p>{{ t('auth.subtitle') }}</p>
        </div>
        <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleLogin">
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              :placeholder="t('auth.enterUsernameHint')"
              size="large"
              :prefix-icon="User"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              :placeholder="t('auth.enterPasswordHint')"
              size="large"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          <el-form-item>
            <el-button v-ripple type="primary" size="large" :loading="loading" native-type="submit" class="login-btn">
              {{ t('auth.login') }}
            </el-button>
          </el-form-item>
        </el-form>
        <div class="locale-switch">
          <div class="lang-toggle">
            <button
              v-for="opt in localeOptions"
              :key="opt.value"
              :class="['lang-btn', { active: currentLang === opt.value }]"
              @click="handleLocaleChange(opt.value)"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>
        <div v-if="allowRegister" class="register-link">
          {{ t('auth.noAccount') }}<router-link to="/register">{{ t('auth.goToRegister') }}</router-link>
        </div>
      </div>
    </div>
  </GradientBg>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { login, checkRegister } from '@/api'
import { useUserStore } from '@/stores/user'
import { useLocale } from '@/composables/useLocale'
import GradientBg from '@/components/GradientBg.vue'

const { t } = useI18n()
const { currentLocale, setLocale, localeOptions } = useLocale()

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const allowRegister = ref(true)
const currentLang = ref(currentLocale.value)

const form = ref({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: () => t('auth.enterUsername'), trigger: 'blur' }],
  password: [{ required: true, message: () => t('auth.enterPassword'), trigger: 'blur' }],
}

function handleLocaleChange(lang: string) {
  setLocale(lang)
}

onMounted(async () => {
  try {
    const res = await checkRegister()
    allowRegister.value = res.allowRegister
  } catch {
    // 非关键，忽略
  }
})

const handleLogin = async () => {
  await formRef.value?.validate()
  loading.value = true
  try {
    const res = await login(form.value)
    userStore.setLogin(res, form.value.username)
    ElMessage.success(t('auth.loginSuccess'))
    router.push('/')
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-card {
  width: 420px;
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 24px;
  padding: 48px 40px 40px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  animation: card-appear 0.6s ease-out;
}

@keyframes card-appear {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.logo-area {
  text-align: center;
  margin-bottom: 36px;
}

.logo-icon {
  font-size: 48px;
  display: inline-block;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-6px);
  }
}

.logo-area h1 {
  font-size: 28px;
  margin: 8px 0 4px;
}

.logo-area p {
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.locale-switch {
  display: flex;
  justify-content: center;
  margin: 8px 0;
}

.lang-toggle {
  display: inline-flex;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 3px;
  gap: 2px;
}

.lang-btn {
  padding: 5px 16px;
  border: none;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.7);
  background: transparent;
  transition: all 0.3s ease;
  font-family: inherit;
}

.lang-btn:hover {
  color: rgba(255, 255, 255, 0.9);
}

.lang-btn.active {
  background: rgba(255, 255, 255, 0.25);
  color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.register-link {
  text-align: center;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  margin-top: 8px;
}

.register-link a {
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
}

.register-link a:hover {
  color: #fff;
}
</style>
