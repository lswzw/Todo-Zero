<template>
  <div class="register-page">
    <div class="register-card">
      <template v-if="allowRegister">
        <div class="logo-area">
          <span class="logo-icon">📝</span>
          <h1>Todo App</h1>
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
            <el-button type="primary" size="large" :loading="loading" native-type="submit" class="register-btn">
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-card {
  width: 420px;
  background: #fff;
  border-radius: 16px;
  padding: 48px 40px 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.logo-area {
  text-align: center;
  margin-bottom: 36px;
}

.logo-icon {
  font-size: 48px;
}

.logo-area h1 {
  font-size: 28px;
  color: #303133;
  margin: 8px 0 4px;
}

.logo-area p {
  color: #909399;
  font-size: 14px;
}

.register-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.locale-switch {
  display: flex;
  justify-content: center;
  margin: 8px 0;
}

.login-link {
  text-align: center;
  color: #909399;
  font-size: 14px;
  margin-top: 8px;
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
  color: #303133;
  margin: 16px 0 8px;
}

.closed-area p {
  color: #909399;
  font-size: 14px;
}

.back-link {
  display: inline-block;
  margin-top: 24px;
  color: var(--primary);
  font-size: 14px;
}
</style>
