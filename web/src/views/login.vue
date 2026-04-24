<template>
  <div class="login-page">
    <div class="login-card">
      <div class="logo-area">
        <span class="logo-icon">📝</span>
        <h1>Todo App</h1>
        <p>个人待办管理平台</p>
      </div>
      <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleLogin">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" size="large" :prefix-icon="User" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" size="large" :prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" native-type="submit" class="login-btn">
            登 录
          </el-button>
        </el-form-item>
      </el-form>
      <div v-if="allowRegister" class="register-link">
        还没有账号？<router-link to="/register">去注册 →</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { login, checkRegister } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const allowRegister = ref(true)

const form = ref({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
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
    ElMessage.success('登录成功')
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
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

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.register-link {
  text-align: center;
  color: #909399;
  font-size: 14px;
  margin-top: 8px;
}
</style>
