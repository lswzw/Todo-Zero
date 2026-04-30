<template>
  <div class="detail-page">
    <!-- 顶部导航 -->
    <header class="navbar">
      <div class="nav-inner">
        <div class="nav-left">
          <el-button text @click="router.push('/')">
            <el-icon><ArrowLeft /></el-icon> {{ t('taskDetail.backToList') }}
          </el-button>
        </div>
        <div class="nav-right">
          <el-button text :disabled="loading" @click="openEditDialog">{{ t('common.edit') }}</el-button>
          <el-button text :disabled="loading" @click="handleToggle">
            {{ task?.status === 2 ? t('taskDetail.markAsTodo') : t('taskDetail.markAsDone') }}
          </el-button>
          <el-popconfirm :title="t('taskDetail.deleteConfirm')" @confirm="handleDelete">
            <template #reference>
              <el-button text type="danger" :disabled="loading">{{ t('common.delete') }}</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>
    </header>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
      <p>{{ t('common.loading') }}</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="loadError" class="error-state">
      <span class="error-icon">😕</span>
      <p>{{ loadError }}</p>
      <el-button type="primary" @click="loadDetail">{{ t('common.retry') }}</el-button>
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
        <el-tag v-if="task.priority === 2" type="danger" effect="dark">{{ t('home.urgent') }}</el-tag>
        <el-tag v-else-if="task.priority === 1" type="warning" effect="dark">{{ t('home.important') }}</el-tag>
        <el-tag v-else type="success" effect="dark">{{ t('home.normal') }}</el-tag>

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
          {{ task.status === 2 ? t('home.completed') : t('home.todo') }}
        </el-tag>
      </div>

      <!-- 信息卡片 -->
      <div class="info-cards">
        <div class="info-card">
          <div class="info-label">{{ t('taskDetail.createTime') }}</div>
          <div class="info-value">{{ task.createTime || '—' }}</div>
        </div>
        <div class="info-card">
          <div class="info-label">{{ t('taskDetail.startTime') }}</div>
          <div class="info-value">{{ task.startTime || '—' }}</div>
        </div>
        <div class="info-card">
          <div class="info-label">{{ t('taskDetail.endTime') }}</div>
          <div class="info-value" :class="{ overdue: isOverdue }">{{ task.endTime || '—' }}</div>
        </div>
        <div class="info-card">
          <div class="info-label">{{ t('taskDetail.reminderTime') }}</div>
          <div class="info-value">{{ task.reminder || '—' }}</div>
        </div>
        <div class="info-card">
          <div class="info-label">{{ t('taskDetail.updateTime') }}</div>
          <div class="info-value">{{ task.updateTime || '—' }}</div>
        </div>
      </div>

      <!-- 内容区 -->
      <div v-if="task.content" class="content-section">
        <h3>{{ t('taskDetail.taskContent') }}</h3>
        <div class="content-body">{{ task.content }}</div>
      </div>
      <div v-else class="content-section empty-content">
        <h3>{{ t('taskDetail.taskContent') }}</h3>
        <p class="no-content">{{ t('taskDetail.noContent') }}</p>
      </div>
    </main>
  </div>

  <!-- 编辑弹窗 -->
  <el-dialog v-model="editDialogVisible" :title="t('home.editTask')" width="480px" destroy-on-close>
    <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="80px">
      <el-form-item :label="t('home.taskTitle')" prop="title">
        <el-input v-model="editForm.title" maxlength="100" :placeholder="t('home.enterTitle')" />
      </el-form-item>
      <el-form-item :label="t('home.taskContent')" prop="content">
        <el-input
          v-model="editForm.content"
          type="textarea"
          :rows="4"
          maxlength="1000"
          :placeholder="t('home.enterContent')"
        />
      </el-form-item>
      <el-form-item :label="t('home.priority')">
        <el-select v-model="editForm.priority" style="width: 100%">
          <el-option :label="t('home.urgent')" :value="2" />
          <el-option :label="t('home.important')" :value="1" />
          <el-option :label="t('home.normal')" :value="3" />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('home.category')">
        <el-select v-model="editForm.categoryId" clearable style="width: 100%">
          <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id">
            <span :style="{ color: c.color || '#909399' }">●</span> {{ c.name }}
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="t('home.startTime')">
        <el-date-picker
          v-model="editForm.startTime"
          type="datetime"
          :placeholder="t('home.selectStartTime')"
          format="YYYY-MM-DD HH:mm"
          value-format="YYYY-MM-DD HH:mm"
          style="width: 100%"
          clearable
        />
      </el-form-item>
      <el-form-item :label="t('home.endTime')">
        <el-date-picker
          v-model="editForm.endTime"
          type="datetime"
          :placeholder="t('home.selectEndTime')"
          format="YYYY-MM-DD HH:mm"
          value-format="YYYY-MM-DD HH:mm"
          style="width: 100%"
          clearable
        />
      </el-form-item>
      <el-form-item :label="t('home.reminderTime')">
        <el-date-picker
          v-model="editForm.reminder"
          type="datetime"
          :placeholder="t('home.selectReminderTime')"
          format="YYYY-MM-DD HH:mm"
          value-format="YYYY-MM-DD HH:mm"
          style="width: 100%"
          clearable
        />
      </el-form-item>
      <el-form-item :label="t('home.tags')">
        <el-input v-model="editForm.tags" maxlength="200" :placeholder="t('home.tagsPlaceholder')" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="editDialogVisible = false">{{ t('common.cancel') }}</el-button>
      <el-button type="primary" :loading="editSubmitting" @click="handleEditSubmit">{{
        t('common.confirm')
      }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { ArrowLeft, Check, Loading } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { getTaskDetail, updateTask, toggleTask, deleteTask, getCategoryList } from '@/api'
import type { TaskDetail, TaskFormData, CategoryItem } from '@/types'

const { t } = useI18n()

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
    { required: true, message: () => t('home.enterTitle'), trigger: 'blur' },
    { max: 100, message: () => t('home.titleMaxLength'), trigger: 'blur' },
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
    loadError.value = t('taskDetail.invalidTaskId')
    loading.value = false
    return
  }
  loading.value = true
  loadError.value = ''
  try {
    task.value = await getTaskDetail(id)
  } catch {
    loadError.value = t('taskDetail.loadFailed')
  } finally {
    loading.value = false
  }
}

async function handleToggle() {
  if (!task.value) return
  try {
    await toggleTask(task.value.id)
    ElMessage.success(task.value.status === 2 ? t('home.markedTodo') : t('home.markedDone'))
    await loadDetail()
  } catch {
    ElMessage.error(t('home.toggleFailed'))
  }
}

async function handleDelete() {
  if (!task.value) return
  try {
    await deleteTask(task.value.id)
    ElMessage.success(t('home.deleted'))
    router.push('/')
  } catch {
    ElMessage.error(t('home.deleteTaskFailed'))
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
    ElMessage.success(t('home.updateSuccess'))
    editDialogVisible.value = false
    await loadDetail()
  } catch {
    ElMessage.error(t('home.updateFailed'))
  } finally {
    editSubmitting.value = false
  }
}
</script>

<style scoped>
.detail-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.navbar {
  position: sticky;
  top: 0;
  z-index: var(--z-sticky);
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
  height: 56px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
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
  transition: all 0.3s ease;
}

.status-indicator:hover {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.15);
}

.status-indicator.done {
  background: linear-gradient(135deg, #67c23a, #38f9d7);
  border-color: #67c23a;
  color: #fff;
  font-size: 14px;
  animation: check-pop 0.4s ease;
}

@keyframes check-pop {
  0% {
    transform: scale(1);
  }
  40% {
    transform: scale(1.2);
  }
  70% {
    transform: scale(0.9);
  }
  100% {
    transform: scale(1);
  }
}

.detail-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  line-height: 1.4;
  word-break: break-word;
  transition: color 0.3s ease;
}

.detail-title.line-through {
  text-decoration: line-through;
  color: #909399;
  text-decoration-color: #909399;
  text-decoration-thickness: 2px;
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
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  transition:
    transform 0.3s ease,
    box-shadow 0.3s ease;
}

.info-card:hover {
  transform: translateY(-4px) scale(1.02);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
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
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
}

.content-section h3 {
  font-size: 15px;
  background: linear-gradient(135deg, #303133, #667eea);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 12px;
  padding-bottom: 10px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
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
