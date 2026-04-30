<template>
  <GradientBg variant="primary" :show-decor="true">
    <div class="register-page">
      <div class="register-card">
        <template v-if="allowRegister">
          <div class="logo-area">
            <span class="logo-icon">📝</span>
            <h1 class="gradient-text">Todo App</h1>
            <p>{{ t('auth.subtitle') }}</p>
          </div>
          <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleRegister">
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
            <el-form-item prop="confirmPassword">
              <el-input
                v-model="form.confirmPassword"
                type="password"
                :placeholder="t('auth.enterConfirmPassword')"
                size="large"
                :prefix-icon="Lock"
                show-password
              />
            </el-form-item>
            <el-form-item>
              <el-button
                v-ripple
                type="primary"
                size="large"
                :loading="loading"
                native-type="submit"
                class="register-btn"
              >
                {{ t('auth.register') }}
              </el-button>
            </el-form-item>
          </el-form>
          <div class="locale-switch">
            <el-select v-model="currentLang" size="small" style="width: 100px" @change="handleLocaleChange">
              <el-option v-for="opt in localeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
            </el-select>
          </div>
          <div class="login-link">
            {{ t('auth.hasAccount') }}<router-link to="/login">{{ t('auth.goToLogin') }}</router-link>
          </div>
        </template>
        <template v-else>
          <div class="closed-area">
            <span class="closed-icon">🔒</span>
            <h2>{{ t('auth.registerClosed') }}</h2>
            <p>{{ t('auth.registerClosedDesc') }}</p>
            <router-link to="/login" class="back-link">{{ t('auth.backToLogin') }}</router-link>
          </div>
        </template>
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
import { register, checkRegister } from '@/api'
import { useLocale } from '@/composables/useLocale'
import GradientBg from '@/components/GradientBg.vue'

const { t } = useI18n()
const { currentLocale, setLocale, localeOptions } = useLocale()

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const allowRegister = ref(true)
const currentLang = ref(currentLocale.value)

const form = ref({ username: '', password: '', confirmPassword: '' })

const validateConfirm = (_rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (value !== form.value.password) {
    callback(new Error(t('auth.passwordMismatch')))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: () => t('auth.enterUsername'), trigger: 'blur' },
    { min: 3, max: 20, message: () => t('auth.usernameLength'), trigger: 'blur' },
  ],
  password: [
    { required: true, message: () => t('auth.enterPassword'), trigger: 'blur' },
    { min: 6, max: 20, message: () => t('auth.passwordLength'), trigger: 'blur' },
    { pattern: /^(?=.*[a-zA-Z])(?=.*\d)/, message: () => t('auth.passwordComplexity'), trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: () => t('auth.enterConfirmPassword'), trigger: 'blur' },
    { validator: validateConfirm, trigger: 'blur' },
  ],
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

const handleRegister = async () => {
  await formRef.value?.validate()
  loading.value = true
  try {
    await register({ username: form.value.username, password: form.value.password })
    ElMessage.success(t('auth.registerSuccess'))
    router.push('/login')
  } catch {
    // 错误已由拦截器处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.register-card {
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

.register-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
}

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.locale-switch {
  display: flex;
  justify-content: center;
  margin: 8px 0;
}

.login-link {
  text-align: center;
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  margin-top: 8px;
}

.login-link a {
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
}

.login-link a:hover {
  color: #fff;
}

.closed-area {
  text-align: center;
  padding: 40px 0;
}

.closed-icon {
  font-size: 48px;
}

.closed-area h2 {
  font-size: 20px;
  color: rgba(255, 255, 255, 0.95);
  margin: 16px 0 8px;
}

.closed-area p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
}

.back-link {
  display: inline-block;
  margin-top: 24px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
  font-weight: 500;
}

.back-link:hover {
  color: #fff;
}
</style>
