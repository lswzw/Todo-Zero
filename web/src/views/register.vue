<template>
  <div class="register-page">
    <div class="register-card">
      <template v-if="allowRegister">
        <div class="logo-area">
          <span class="logo-icon">📝</span>
          <h1>Todo App</h1>
          <p>个人待办管理平台</p>
        </div>
        <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleRegister">
          <el-form-item prop="username">
            <el-input v-model="form.username" placeholder="请输入用户名(3-20位)" size="large" :prefix-icon="User" />
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="form.password" type="password" placeholder="请输入密码(6-20位)" size="large" :prefix-icon="Lock" show-password />
          </el-form-item>
          <el-form-item prop="confirmPassword">
            <el-input v-model="form.confirmPassword" type="password" placeholder="请确认密码" size="large" :prefix-icon="Lock" show-password />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="large" :loading="loading" native-type="submit" class="register-btn">
              注 册
            </el-button>
          </el-form-item>
        </el-form>
        <div class="login-link">
          已有账号？<router-link to="/login">去登录 →</router-link>
        </div>
      </template>
      <template v-else>
        <div class="closed-area">
          <span class="closed-icon">🔒</span>
          <h2>注册功能暂未开放</h2>
          <p>请联系管理员开通账号</p>
          <router-link to="/login" class="back-link">← 返回登录</router-link>
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
import { register, checkRegister } from '@/api'

const router = useRouter()
const formRef = ref()
const loading = ref(false)
const allowRegister = ref(true)

const form = ref({ username: '', password: '', confirmPassword: '' })

const validateConfirm = (_rule: any, value: string, callback: any) => {
  if (value !== form.value.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度3-20位', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度6-20位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirm, trigger: 'blur' },
  ],
}

onMounted(async () => {
  try {
    const res = await checkRegister() as any
    allowRegister.value = res.allowRegister
  } catch {}
})

const handleRegister = async () => {
  await formRef.value?.validate()
  loading.value = true
  try {
    await register({ username: form.value.username, password: form.value.password })
    ElMessage.success('注册成功')
    setTimeout(() => router.push('/login'), 800)
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
