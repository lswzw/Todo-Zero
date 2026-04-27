<template>
  <div class="detail-page">
    <!-- 顶部导航 -->
    <header class="navbar">
      <div class="nav-inner">
        <div class="nav-left">
          <el-button text @click="router.push('/')">
            <el-icon><ArrowLeft /></el-icon> 返回列表
          </el-button>
        </div>
        <div class="nav-right">
          <el-button text :disabled="loading" @click="openEditDialog">编辑</el-button>
          <el-button text :disabled="loading" @click="handleToggle">
            {{ task?.status === 2 ? '标为待办' : '标为完成' }}
          </el-button>
          <el-popconfirm title="确定删除该任务？" @confirm="handleDelete">
            <template #reference>
              <el-button text type="danger" :disabled="loading">删除</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>
    </header>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
      <p>加载中...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="loadError" class="error-state">
      <span class="error-icon">😕</span>
      <p>{{ loadError }}</p>
      <el-button type="primary" @click="loadDetail">重试</el-button>
    </div>

    <!-- 任务详情 -->
    <main v-else-if="task" class="main-content">
      <!-- 标题区 -->
      <div class="detail-header">
        <div class="status-indicator" :class="{ done: task.status === 2 }" @click="handleToggle">
          <el-icon v-if="task.status === 2"><Check /></el-icon>
        </div>
        <h1 :class="['detail-title', { 'line-through': task.status === 2 }]">{{ task.title }}</h1>
      </div>

      <!-- 元信息区 -->
      <div class="meta-row">
        <el-tag v-if="task.priority === 2" type="danger" effect="dark">紧急</el-tag>
        <el-tag v-else-if="task.priority === 1" type="warning" effect="dark">重要</el-tag>
        <el-tag v-else type="success" effect="dark">普通</el-tag>

        <el-tag
          v-if="task.categoryName"
          type="info"
          :color="categoryColor"
          style="border-color: transparent"
          :style="{ color: categoryTextColor }"
        >
          {{ task.categoryName }}
        </el-tag>

        <el-tag v-for="tag in parseTags(task.tags)" :key="tag" effect="plain" size="small" class="task-tag">{{
          tag
        }}</el-tag>

        <el-tag :type="task.status === 2 ? 'success' : 'warning'" size="small">
          {{ task.status === 2 ? '已完成' : '待办' }}
        </el-tag>
      </div>

      <!-- 信息卡片 -->
      <div class="info-cards">
        <div class="info-card">
          <div class="info-label">创建时间</div>
          <div class="info-value">{{ task.createTime || '—' }}</div>
        </div>
        <div class="info-card">
          <div class="info-label">开始时间</div>
          <div class="info-value">{{ task.startTime || '—' }}</div>
        </div>
        <div class="info-card">
          <div class="info-label">截止时间</div>
          <div class="info-value" :class="{ overdue: isOverdue }">{{ task.endTime || '—' }}</div>
        </div>
        <div class="info-card">
          <div class="info-label">提醒时间</div>
          <div class="info-value">{{ task.reminder || '—' }}</div>
        </div>
        <div class="info-card">
          <div class="info-label">更新时间</div>
          <div class="info-value">{{ task.updateTime || '—' }}</div>
        </div>
      </div>

      <!-- 内容区 -->
      <div v-if="task.content" class="content-section">
        <h3>任务内容</h3>
        <div class="content-body">{{ task.content }}</div>
      </div>
      <div v-else class="content-section empty-content">
        <h3>任务内容</h3>
        <p class="no-content">暂无详细内容</p>
      </div>
    </main>
  </div>

  <!-- 编辑弹窗 -->
  <el-dialog v-model="editDialogVisible" title="编辑任务" width="480px" destroy-on-close>
    <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="80px">
      <el-form-item label="标题" prop="title">
        <el-input v-model="editForm.title" maxlength="100" placeholder="请输入任务标题" />
      </el-form-item>
      <el-form-item label="内容" prop="content">
        <el-input
          v-model="editForm.content"
          type="textarea"
          :rows="4"
          maxlength="1000"
          placeholder="任务详细内容（选填）"
        />
      </el-form-item>
      <el-form-item label="优先级">
        <el-select v-model="editForm.priority" style="width: 100%">
          <el-option label="紧急" :value="2" />
          <el-option label="重要" :value="1" />
          <el-option label="普通" :value="3" />
        </el-select>
      </el-form-item>
      <el-form-item label="分类">
        <el-select v-model="editForm.categoryId" clearable style="width: 100%">
          <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id">
            <span :style="{ color: c.color || '#909399' }">●</span> {{ c.name }}
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="开始时间">
        <el-date-picker
          v-model="editForm.startTime"
          type="datetime"
          placeholder="选择开始时间"
          format="YYYY-MM-DD HH:mm"
          value-format="YYYY-MM-DD HH:mm"
          style="width: 100%"
          clearable
        />
      </el-form-item>
      <el-form-item label="截止时间">
        <el-date-picker
          v-model="editForm.endTime"
          type="datetime"
          placeholder="选择截止时间"
          format="YYYY-MM-DD HH:mm"
          value-format="YYYY-MM-DD HH:mm"
          style="width: 100%"
          clearable
        />
      </el-form-item>
      <el-form-item label="提醒时间">
        <el-date-picker
          v-model="editForm.reminder"
          type="datetime"
          placeholder="选择提醒时间"
          format="YYYY-MM-DD HH:mm"
          value-format="YYYY-MM-DD HH:mm"
          style="width: 100%"
          clearable
        />
      </el-form-item>
      <el-form-item label="标签">
        <el-input v-model="editForm.tags" maxlength="200" placeholder="多个标签用逗号分隔，如：工作,重要" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="editDialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="editSubmitting" @click="handleEditSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { ArrowLeft, Check, Loading } from '@element-plus/icons-vue'
import { getTaskDetail, updateTask, toggleTask, deleteTask, getCategoryList } from '@/api'
import type { TaskDetail, TaskFormData, CategoryItem } from '@/types'

const route = useRoute()
const router = useRouter()

const task = ref<TaskDetail | null>(null)
const loading = ref(true)
const loadError = ref('')
const categories = ref<CategoryItem[]>([])

// 编辑弹窗
const editDialogVisible = ref(false)
const editSubmitting = ref(false)
const editFormRef = ref<FormInstance>()
const editForm = ref<TaskFormData>({
  title: '',
  content: '',
  priority: 3,
  categoryId: undefined,
  startTime: '',
  endTime: '',
  reminder: '',
  tags: '',
})
const editRules = {
  title: [
    { required: true, message: '请输入任务标题', trigger: 'blur' },
    { max: 100, message: '标题最长100字符', trigger: 'blur' },
  ],
}

const isOverdue = computed(() => {
  if (!task.value || task.value.status === 2 || !task.value.endTime) return false
  return new Date(task.value.endTime) < new Date()
})

const categoryColor = computed(() => {
  if (!task.value?.categoryId) return undefined
  const cat = categories.value.find((c) => c.id === task.value!.categoryId)
  return cat?.color || undefined
})

const categoryTextColor = computed(() => {
  if (!task.value?.categoryId) return '#909399'
  const cat = categories.value.find((c) => c.id === task.value!.categoryId)
  if (!cat?.color) return '#909399'
  const hex = cat.color.replace('#', '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  const lum = (0.299 * r + 0.587 * g + 0.114 * b) / 255
  return lum > 0.6 ? '#303133' : '#ffffff'
})

function parseTags(tags: string): string[] {
  if (!tags) return []
  return tags
    .split(',')
    .map((t) => t.trim())
    .filter(Boolean)
}

onMounted(() => {
  loadCategories()
  loadDetail()
})

async function loadCategories() {
  try {
    const res = await getCategoryList()
    categories.value = res.list || []
  } catch {
    // 非关键数据，静默降级
  }
}

async function loadDetail() {
  const id = Number(route.params.id)
  if (!id) {
    loadError.value = '无效的任务ID'
    loading.value = false
    return
  }
  loading.value = true
  loadError.value = ''
  try {
    task.value = await getTaskDetail(id)
  } catch {
    loadError.value = '加载任务详情失败'
  } finally {
    loading.value = false
  }
}

async function handleToggle() {
  if (!task.value) return
  try {
    await toggleTask(task.value.id)
    ElMessage.success(task.value.status === 2 ? '已标记待办' : '已标记完成')
    await loadDetail()
  } catch {
    ElMessage.error('切换状态失败')
  }
}

async function handleDelete() {
  if (!task.value) return
  try {
    await deleteTask(task.value.id)
    ElMessage.success('已删除')
    router.push('/')
  } catch {
    ElMessage.error('删除任务失败')
  }
}

function openEditDialog() {
  if (!task.value) return
  editForm.value = {
    title: task.value.title,
    content: task.value.content || '',
    priority: task.value.priority,
    categoryId: task.value.categoryId || undefined,
    startTime: task.value.startTime || '',
    endTime: task.value.endTime || '',
    reminder: task.value.reminder || '',
    tags: task.value.tags || '',
  }
  editDialogVisible.value = true
}

async function handleEditSubmit() {
  await editFormRef.value?.validate()
  if (!task.value) return
  editSubmitting.value = true
  try {
    await updateTask(task.value.id, editForm.value)
    ElMessage.success('修改成功')
    editDialogVisible.value = false
    await loadDetail()
  } catch {
    ElMessage.error('修改任务失败')
  } finally {
    editSubmitting.value = false
  }
}
</script>

<style scoped>
.detail-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: #fff;
  height: 56px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.nav-inner {
  max-width: 720px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.nav-left,
.nav-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 120px 20px;
  color: #909399;
}

.error-icon {
  font-size: 48px;
}

.main-content {
  max-width: 720px;
  margin: 0 auto;
  padding: 24px 20px;
}

.detail-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 16px;
}

.status-indicator {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  border: 2px solid #c0c4cc;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  flex-shrink: 0;
  margin-top: 2px;
  transition: all 0.2s;
}

.status-indicator.done {
  background: #67c23a;
  border-color: #67c23a;
  color: #fff;
  font-size: 14px;
}

.detail-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  line-height: 1.4;
  word-break: break-word;
}

.detail-title.line-through {
  text-decoration: line-through;
  color: #909399;
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.task-tag {
  border-radius: 12px;
  font-size: 12px;
}

.info-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 12px;
  margin-bottom: 24px;
}

.info-card {
  background: #fff;
  border-radius: 10px;
  padding: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.info-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 6px;
}

.info-value {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.info-value.overdue {
  color: #f56c6c;
}

.content-section {
  background: #fff;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.content-section h3 {
  font-size: 15px;
  color: #303133;
  margin-bottom: 12px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

.content-body {
  font-size: 14px;
  color: #606266;
  line-height: 1.8;
  white-space: pre-wrap;
  word-break: break-word;
}

.no-content {
  color: #c0c4cc;
  font-size: 14px;
}

@media (max-width: 768px) {
  .info-cards {
    grid-template-columns: repeat(2, 1fr);
  }
  .detail-title {
    font-size: 20px;
  }
}
</style>
