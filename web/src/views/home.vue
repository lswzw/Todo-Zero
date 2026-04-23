<template>
  <div class="home-page">
    <!-- 顶部导航 -->
    <header class="navbar">
      <div class="nav-inner">
        <div class="nav-left">
          <span class="logo-icon">📝</span>
          <span class="logo-text">Todo App</span>
        </div>
        <div class="nav-right">
          <span class="username">{{ userStore.username }}</span>
          <el-button text @click="showPasswordDialog = true">修改密码</el-button>
          <el-button text type="danger" @click="handleLogout">退出登录</el-button>
        </div>
      </div>
    </header>

    <!-- 主内容 -->
    <main class="main-content">
      <!-- 统计卡片 -->
      <div class="stat-row">
        <div class="stat-card total">
          <div class="stat-value">{{ stat.total }}</div>
          <div class="stat-label">总任务</div>
        </div>
        <div class="stat-card todo">
          <div class="stat-value">{{ stat.todo }}</div>
          <div class="stat-label">待办</div>
        </div>
        <div class="stat-card done">
          <div class="stat-value">{{ stat.done }}</div>
          <div class="stat-label">已完成</div>
        </div>
        <div class="stat-card rate">
          <div class="stat-value">{{ stat.doneRate }}%</div>
          <div class="stat-label">完成率</div>
        </div>
      </div>

      <!-- 任务列表区 -->
      <div class="task-section">
        <div class="section-header">
          <h2>任务列表</h2>
          <div class="section-actions">
            <el-select v-model="filters.status" placeholder="状态" clearable style="width: 100px" @change="loadTasks">
              <el-option label="待办" :value="0" />
              <el-option label="已完成" :value="2" />
            </el-select>
            <el-select v-model="filters.categoryId" placeholder="分类" clearable style="width: 100px" @change="loadTasks">
              <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
            </el-select>
            <el-select v-model="filters.priority" placeholder="优先级" clearable style="width: 100px" @change="loadTasks">
              <el-option label="紧急" :value="2" />
              <el-option label="重要" :value="1" />
              <el-option label="普通" :value="3" />
            </el-select>
            <el-button :type="selectMode ? 'primary' : ''" @click="toggleSelectMode">
              {{ selectMode ? '退出多选' : '多选' }}
            </el-button>
            <el-input v-model="filters.keyword" placeholder="搜索" clearable style="width: 180px" @clear="loadTasks" @keyup.enter="loadTasks">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-button type="primary" @click="openTaskDialog()">
              <el-icon><Plus /></el-icon> 新增任务
            </el-button>
          </div>
        </div>

        <!-- 批量操作栏 -->
        <div v-if="selectMode && selectedIds.length > 0" class="batch-bar">
          <span>已选 {{ selectedIds.length }} 项</span>
          <el-button size="small" type="success" @click="handleBatch('complete')">批量完成</el-button>
          <el-button size="small" type="warning" @click="handleBatch('undo')">批量取消</el-button>
          <el-button size="small" type="danger" @click="handleBatch('delete')">批量删除</el-button>
          <el-button size="small" @click="selectedIds = []">取消选择</el-button>
        </div>

        <!-- 任务列表 -->
        <div v-if="tasks.length" class="task-list">
          <div v-for="task in tasks" :key="task.id" class="task-item">
            <div v-if="selectMode" class="task-check" @click="toggleSelect(task.id)">
              <div :class="['check-dot', { active: selectedIds.includes(task.id) }]" />
            </div>
            <div class="task-status" @click="handleToggle(task)">
              <div :class="['status-circle', { done: task.status === 2 }]">
                <el-icon v-if="task.status === 2"><Check /></el-icon>
              </div>
            </div>
            <div class="task-body">
              <div :class="['task-title', { 'line-through': task.status === 2 }]">{{ task.title }}</div>
              <div v-if="task.content" class="task-content">{{ task.content }}</div>
              <div class="task-meta">
                <el-tag v-if="task.priority === 2" size="small" type="danger">紧急</el-tag>
                <el-tag v-else-if="task.priority === 1" size="small" type="warning">重要</el-tag>
                <el-tag v-else size="small" type="success">普通</el-tag>
                <el-tag v-if="task.categoryName" size="small" type="info">{{ task.categoryName }}</el-tag>
                <span class="task-time">{{ task.createTime }}</span>
              </div>
            </div>
            <div class="task-actions">
              <el-button text size="small" @click="openTaskDialog(task)">编辑</el-button>
              <el-popconfirm title="确定删除该任务？" @confirm="handleDelete(task.id)">
                <template #reference>
                  <el-button text size="small" type="danger">删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">
          <span class="empty-icon">📋</span>
          <p>暂无任务</p>
        </div>

        <!-- 分页 -->
        <div v-if="total > 0" class="pagination">
          <el-pagination
            v-model:current-page="page"
            :page-size="pageSize"
            :total="total"
            layout="prev, pager, next"
            @current-change="loadTasks"
          />
        </div>
      </div>
    </main>

    <!-- 新增/编辑任务弹窗 -->
    <el-dialog v-model="taskDialogVisible" :title="editingTask ? '编辑任务' : '新增任务'" width="480px" destroy-on-close>
      <el-form ref="taskFormRef" :model="taskForm" :rules="taskRules" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="taskForm.title" maxlength="100" placeholder="请输入任务标题" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="taskForm.content" type="textarea" :rows="4" maxlength="1000" placeholder="任务详细内容（选填）" />
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="taskForm.priority" style="width: 100%">
            <el-option label="紧急" :value="2" />
            <el-option label="重要" :value="1" />
            <el-option label="普通" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="taskForm.categoryId" clearable style="width: 100%">
            <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="taskDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmitTask">确定</el-button>
      </template>
    </el-dialog>

    <!-- 修改密码弹窗 -->
    <el-dialog v-model="showPasswordDialog" title="修改密码" width="420px" destroy-on-close>
      <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-width="90px">
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input v-model="pwdForm.oldPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="pwdForm.newPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="pwdForm.confirmPassword" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" :loading="pwdLoading" @click="handleChangePassword">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Search, Plus, Check } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { resetAuthVerified } from '@/router'
import {
  getTaskList, createTask, updateTask, toggleTask, deleteTask, batchTask,
  getCategoryList, getStat, changePassword,
} from '@/api'
import type { TaskItem, StatResp, CategoryItem } from '@/types'

const router = useRouter()
const userStore = useUserStore()

// 统计
const stat = ref<StatResp>({ total: 0, done: 0, todo: 0, doneRate: 0 })

// 分类
const categories = ref<CategoryItem[]>([])

// 任务列表
const tasks = ref<TaskItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const filters = ref<Record<string, number | string | undefined>>({
  status: undefined, categoryId: undefined, priority: undefined, keyword: '',
})

// 多选
const selectMode = ref(false)
const selectedIds = ref<number[]>([])

// 任务弹窗
const taskDialogVisible = ref(false)
const editingTask = ref<TaskItem | null>(null)
const submitting = ref(false)
const taskFormRef = ref<FormInstance>()
const taskForm = ref({ title: '', content: '', priority: 3, categoryId: undefined as number | undefined })
const taskRules = {
  title: [
    { required: true, message: '请输入任务标题', trigger: 'blur' },
    { max: 100, message: '标题最长100字符', trigger: 'blur' },
  ],
}

// 密码弹窗
const showPasswordDialog = ref(false)
const pwdLoading = ref(false)
const pwdFormRef = ref<FormInstance>()
const pwdForm = ref({ oldPassword: '', newPassword: '', confirmPassword: '' })
const validatePwdConfirm = (_rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (value !== pwdForm.value.newPassword) callback(new Error('两次输入的密码不一致'))
  else callback()
}
const pwdRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度6-20位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validatePwdConfirm, trigger: 'blur' },
  ],
}

onMounted(() => {
  loadStat()
  loadCategories()
  loadTasks()
})

async function loadStat() {
  try {
    stat.value = await getStat()
  } catch {
    ElMessage.error('加载统计数据失败')
  }
}

async function loadCategories() {
  try {
    const res = await getCategoryList()
    categories.value = res.list || []
  } catch {
    ElMessage.error('加载分类失败')
  }
}

async function loadTasks() {
  try {
    const params: Record<string, unknown> = { page: page.value, pageSize: pageSize.value }
    if (filters.value.status !== undefined && filters.value.status !== '') params.status = filters.value.status
    if (filters.value.categoryId) params.categoryId = filters.value.categoryId
    if (filters.value.priority) params.priority = filters.value.priority
    if (filters.value.keyword) params.keyword = filters.value.keyword
    const res = await getTaskList(params)
    tasks.value = res.list || []
    total.value = res.total || 0
  } catch {
    ElMessage.error('加载任务列表失败')
  }
}

function toggleSelectMode() {
  selectMode.value = !selectMode.value
  selectedIds.value = []
}

function toggleSelect(id: number) {
  const idx = selectedIds.value.indexOf(id)
  if (idx > -1) selectedIds.value.splice(idx, 1)
  else selectedIds.value.push(id)
}

async function handleToggle(task: TaskItem) {
  try {
    await toggleTask(task.id)
    ElMessage.success(task.status === 0 ? '已标记完成' : '已标记待办')
    loadTasks()
    loadStat()
  } catch {
    // 错误已由拦截器处理
  }
}

async function handleDelete(id: number) {
  try {
    await deleteTask(id)
    ElMessage.success('已删除')
    loadTasks()
    loadStat()
  } catch {
    // 错误已由拦截器处理
  }
}

async function handleBatch(action: string) {
  try {
    await batchTask({ ids: selectedIds.value, action })
    ElMessage.success('操作成功')
    selectedIds.value = []
    loadTasks()
    loadStat()
  } catch {
    // 错误已由拦截器处理
  }
}

function openTaskDialog(task?: TaskItem) {
  editingTask.value = task || null
  if (task) {
    taskForm.value = { title: task.title, content: task.content || '', priority: task.priority, categoryId: task.categoryId || undefined }
  } else {
    taskForm.value = { title: '', content: '', priority: 3, categoryId: undefined }
  }
  taskDialogVisible.value = true
}

async function handleSubmitTask() {
  await taskFormRef.value?.validate()
  submitting.value = true
  try {
    if (editingTask.value) {
      await updateTask(editingTask.value.id, taskForm.value)
      ElMessage.success('修改成功')
    } else {
      await createTask(taskForm.value)
      ElMessage.success('创建成功')
    }
    taskDialogVisible.value = false
    loadTasks()
    loadStat()
  } catch {
    // 错误已由拦截器处理
  } finally {
    submitting.value = false
  }
}

async function handleChangePassword() {
  await pwdFormRef.value?.validate()
  pwdLoading.value = true
  try {
    await changePassword({ oldPassword: pwdForm.value.oldPassword, newPassword: pwdForm.value.newPassword })
    ElMessage.success('密码修改成功，请重新登录')
    showPasswordDialog.value = false
    resetAuthVerified()
    userStore.logout()
    router.push('/login')
  } catch {
    // 错误已由拦截器处理
  } finally {
    pwdLoading.value = false
  }
}

function handleLogout() {
  ElMessageBox.confirm('确定退出登录？', '提示', { type: 'warning' }).then(() => {
    resetAuthVerified()
    userStore.logout()
    router.push('/login')
  }).catch(() => {})
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: #fff;
  height: 60px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.nav-inner {
  max-width: 960px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.nav-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-icon { font-size: 24px; }
.logo-text { font-size: 20px; font-weight: 600; color: #303133; }

.nav-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  color: #606266;
  font-size: 14px;
  margin-right: 8px;
}

.main-content {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px 20px;
}

.stat-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
}

.stat-card.total .stat-value { color: #303133; }
.stat-card.todo .stat-value { color: #e6a23c; }
.stat-card.done .stat-value { color: #67c23a; }
.stat-card.rate .stat-value { color: #667eea; }

.task-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 20px;
}

.section-header h2 {
  font-size: 18px;
  color: #303133;
}

.section-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.batch-bar {
  background: #ecf5ff;
  padding: 10px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #409eff;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-item {
  display: flex;
  align-items: flex-start;
  padding: 14px 12px;
  border-radius: 8px;
  transition: background 0.2s;
}

.task-item:hover {
  background: #f5f7fa;
}

.task-check {
  padding: 4px 8px 0 0;
  cursor: pointer;
}

.check-dot {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 2px solid #c0c4cc;
  transition: all 0.2s;
}

.check-dot.active {
  background: #667eea;
  border-color: #667eea;
}

.task-status {
  padding: 2px 8px 0 0;
  cursor: pointer;
}

.status-circle {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 2px solid #c0c4cc;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.status-circle.done {
  background: #67c23a;
  border-color: #67c23a;
  color: #fff;
  font-size: 12px;
}

.task-body {
  flex: 1;
  min-width: 0;
}

.task-title {
  font-size: 15px;
  color: #303133;
  line-height: 1.5;
}

.task-title.line-through {
  text-decoration: line-through;
  color: #909399;
}

.task-content {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
  line-height: 1.4;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
  flex-wrap: wrap;
}

.task-time {
  font-size: 12px;
  color: #c0c4cc;
}

.task-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: 8px;
  flex-shrink: 0;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #909399;
}

.empty-icon {
  font-size: 48px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

@media (max-width: 768px) {
  .stat-row { grid-template-columns: repeat(2, 1fr); }
  .section-actions { flex-direction: column; align-items: stretch; }
}
</style>
