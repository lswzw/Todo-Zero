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
          <el-select v-model="currentLang" size="small" style="width: 90px" @change="handleLocaleChange">
            <el-option v-for="opt in localeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
          <span class="username">{{ userStore.username }}</span>
          <el-button v-if="userStore.isAdmin" text type="warning" @click="$router.push('/admin')">{{
            t('home.adminPanel')
          }}</el-button>
          <el-button text @click="$router.push('/trash')">{{ t('home.trash') }}</el-button>
          <el-button text @click="showPasswordDialog = true">{{ t('auth.changePassword') }}</el-button>
          <el-button text type="danger" @click="handleLogout">{{ t('common.logout') }}</el-button>
        </div>
      </div>
    </header>

    <!-- 主内容 -->
    <main class="main-content">
      <!-- 统计卡片 -->
      <div class="stat-row">
        <div class="stat-card total">
          <div class="stat-icon">📊</div>
          <div class="stat-value">{{ stat.total }}</div>
          <div class="stat-label">{{ t('home.totalTasks') }}</div>
        </div>
        <div class="stat-card todo">
          <div class="stat-icon">⏳</div>
          <div class="stat-value">{{ stat.todo }}</div>
          <div class="stat-label">{{ t('home.todo') }}</div>
        </div>
        <div class="stat-card done">
          <div class="stat-icon">✅</div>
          <div class="stat-value">{{ stat.done }}</div>
          <div class="stat-label">{{ t('home.completed') }}</div>
        </div>
        <div class="stat-card rate">
          <div class="stat-icon">📈</div>
          <div class="stat-value">{{ stat.doneRate }}%</div>
          <div class="stat-label">{{ t('home.completionRate') }}</div>
        </div>
      </div>

      <!-- 任务列表区 -->
      <div class="task-section">
        <div class="section-header">
          <h2>{{ t('home.taskList') }}</h2>
          <div class="section-actions">
            <el-select
              v-model="filters.status"
              :placeholder="t('home.status')"
              clearable
              style="width: 100px"
              @change="onFilterChange"
            >
              <el-option :label="t('home.todo')" :value="0" />
              <el-option :label="t('home.completed')" :value="2" />
            </el-select>
            <el-select
              v-model="filters.categoryId"
              :placeholder="t('home.category')"
              clearable
              style="width: 100px"
              @change="onFilterChange"
            >
              <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id">
                <span :style="{ color: c.color || '#909399' }">●</span> {{ c.name }}
              </el-option>
            </el-select>
            <el-select
              v-model="filters.priority"
              :placeholder="t('home.priority')"
              clearable
              style="width: 100px"
              @change="onFilterChange"
            >
              <el-option :label="t('home.urgent')" :value="2" />
              <el-option :label="t('home.important')" :value="1" />
              <el-option :label="t('home.normal')" :value="3" />
            </el-select>
            <el-button :type="selectMode ? 'primary' : ''" @click="toggleSelectMode">
              {{ selectMode ? t('home.exitMultiSelect') : t('home.multiSelect') }}
            </el-button>
            <!-- 防抖预留：若将来改为 @input 实时搜索，需引入 debounce（lodash-es 或手写），避免频繁请求 -->
            <el-input
              v-model="filters.keyword"
              :placeholder="t('common.search')"
              clearable
              style="width: 180px"
              @clear="onFilterChange"
              @keyup.enter="onFilterChange"
            >
              <template #prefix
                ><el-icon><Search /></el-icon
              ></template>
            </el-input>
            <el-button type="primary" @click="openTaskDialog()">
              <el-icon><Plus /></el-icon> {{ t('home.newTask') }}
            </el-button>
            <el-dropdown @command="handleExport">
              <el-button>
                <el-icon><Download /></el-icon> {{ t('common.export') }}
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="json">{{ t('common.exportJson') }}</el-dropdown-item>
                  <el-dropdown-item command="csv">{{ t('common.exportCsv') }}</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button @click="showCategoryDialog = true">{{ t('home.categoryManage') }}</el-button>
          </div>
        </div>

        <!-- 批量操作栏 -->
        <div v-if="selectMode && selectedIds.length > 0" class="batch-bar">
          <span>{{ t('home.selected', { count: selectedIds.length }) }}</span>
          <el-button size="small" type="success" @click="handleBatch('complete')">{{
            t('home.batchComplete')
          }}</el-button>
          <el-button size="small" type="warning" @click="handleBatch('undo')">{{ t('home.batchUndo') }}</el-button>
          <el-button size="small" type="danger" @click="handleBatch('delete')">{{ t('home.batchDelete') }}</el-button>
          <el-button size="small" @click="selectedIds = []">{{ t('home.cancelSelect') }}</el-button>
        </div>

        <!-- 任务列表 -->
        <div v-if="tasks.length" class="task-list">
          <draggable
            v-model="tasks"
            item-key="id"
            handle=".drag-handle"
            ghost-class="ghost"
            animation="200"
            @end="onDragEnd"
          >
            <template #item="{ element: task }">
              <div class="task-item">
                <div class="drag-handle" :title="t('home.sortSaveFailed')">⠿</div>
                <div v-if="selectMode" class="task-check" @click="toggleSelect(task.id)">
                  <div :class="['check-dot', { active: selectedIds.includes(task.id) }]" />
                </div>
                <div class="task-status" @click="handleToggle(task)">
                  <div :class="['status-circle', { done: task.status === 2 }]">
                    <el-icon v-if="task.status === 2"><Check /></el-icon>
                  </div>
                </div>
                <div class="task-body" @click="router.push(`/task/${task.id}`)">
                  <div :class="['task-title', { 'line-through': task.status === 2 }]">{{ task.title }}</div>
                  <div v-if="task.content" class="task-content">{{ task.content }}</div>
                  <div class="task-meta">
                    <el-tag v-if="task.priority === 2" size="small" type="danger">{{ t('home.urgent') }}</el-tag>
                    <el-tag v-else-if="task.priority === 1" size="small" type="warning">{{
                      t('home.important')
                    }}</el-tag>
                    <el-tag v-else size="small" type="success">{{ t('home.normal') }}</el-tag>
                    <el-tag
                      v-if="task.categoryName"
                      size="small"
                      type="info"
                      :color="getCategoryColor(task.categoryId)"
                      style="border-color: transparent"
                      :style="{ color: getCategoryTextColor(task.categoryId) }"
                      >{{ task.categoryName }}</el-tag
                    >
                    <el-tag
                      v-for="tag in parseTags(task.tags)"
                      :key="tag"
                      size="small"
                      effect="plain"
                      class="task-tag"
                      >{{ tag }}</el-tag
                    >
                    <span
                      v-if="task.endTime"
                      class="task-time"
                      :class="{ overdue: isOverdue(task.endTime, task.status) }"
                      >{{ t('home.deadline', { time: task.endTime }) }}</span
                    >
                    <span v-else class="task-time">{{ task.createTime }}</span>
                  </div>
                </div>
                <div class="task-actions">
                  <el-button text size="small" @click="openTaskDialog(task)">{{ t('common.edit') }}</el-button>
                  <el-popconfirm :title="t('home.deleteConfirm')" @confirm="handleDelete(task.id)">
                    <template #reference>
                      <el-button text size="small" type="danger">{{ t('common.delete') }}</el-button>
                    </template>
                  </el-popconfirm>
                </div>
              </div>
            </template>
          </draggable>
        </div>
        <div v-else class="empty-state">
          <span class="empty-icon">📋</span>
          <p>{{ t('home.noTasks') }}</p>
        </div>

        <!-- 虚拟滚动预留：当前 pageSize=10 性能无问题；若未来增大分页或取消分页，需引入 vue-virtual-scroller 或 @tanstack/vue-virtual -->
        <!-- 分页 -->
        <div v-if="total > 0" class="pagination">
          <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="prev, pager, next" />
        </div>
      </div>
    </main>

    <!-- 新增/编辑任务弹窗 -->
    <el-dialog
      v-model="taskDialogVisible"
      :title="editingTask ? t('home.editTask') : t('home.newTask')"
      width="480px"
      destroy-on-close
    >
      <el-form ref="taskFormRef" :model="taskForm" :rules="taskRules" label-width="80px">
        <el-form-item :label="t('home.taskTitle')" prop="title">
          <el-input v-model="taskForm.title" maxlength="100" :placeholder="t('home.enterTitle')" />
        </el-form-item>
        <el-form-item :label="t('home.taskContent')" prop="content">
          <el-input
            v-model="taskForm.content"
            type="textarea"
            :rows="4"
            maxlength="1000"
            :placeholder="t('home.enterContent')"
          />
        </el-form-item>
        <el-form-item :label="t('home.priority')">
          <el-select v-model="taskForm.priority" style="width: 100%">
            <el-option :label="t('home.urgent')" :value="2" />
            <el-option :label="t('home.important')" :value="1" />
            <el-option :label="t('home.normal')" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('home.category')">
          <el-select v-model="taskForm.categoryId" clearable style="width: 100%">
            <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id">
              <span :style="{ color: c.color || '#909399' }">●</span> {{ c.name }}
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="t('home.startTime')">
          <el-date-picker
            v-model="taskForm.startTime"
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
            v-model="taskForm.endTime"
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
            v-model="taskForm.reminder"
            type="datetime"
            :placeholder="t('home.selectReminderTime')"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DD HH:mm"
            style="width: 100%"
            clearable
          />
        </el-form-item>
        <el-form-item :label="t('home.tags')">
          <el-input v-model="taskForm.tags" maxlength="200" :placeholder="t('home.tagsPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="taskDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmitTask">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 修改密码弹窗 -->
    <el-dialog v-model="showPasswordDialog" :title="t('auth.changePassword')" width="420px" destroy-on-close>
      <el-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules" label-width="90px">
        <el-form-item :label="t('auth.currentPassword')" prop="oldPassword">
          <el-input v-model="pwdForm.oldPassword" type="password" show-password />
        </el-form-item>
        <el-form-item :label="t('auth.newPassword')" prop="newPassword">
          <el-input v-model="pwdForm.newPassword" type="password" show-password />
        </el-form-item>
        <el-form-item :label="t('auth.confirmNewPassword')" prop="confirmPassword">
          <el-input v-model="pwdForm.confirmPassword" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="pwdLoading" @click="handleChangePassword">{{
          t('common.confirm')
        }}</el-button>
      </template>
    </el-dialog>

    <!-- 分类管理弹窗 -->
    <el-dialog v-model="showCategoryDialog" :title="t('home.categoryManage')" width="480px" destroy-on-close>
      <div class="category-manage">
        <div class="category-add-row">
          <el-input
            v-model="newCategoryName"
            :placeholder="t('home.categoryName')"
            maxlength="20"
            style="flex: 1"
            @keyup.enter="handleAddCategory"
          />
          <el-color-picker v-model="newCategoryColor" size="small" />
          <el-button type="primary" @click="handleAddCategory">{{ t('common.add') }}</el-button>
        </div>
        <div class="category-list">
          <div v-for="c in categories" :key="c.id" class="category-item">
            <el-color-picker v-model="c._color" size="small" @change="handleUpdateCategory(c)" />
            <el-input v-model="c._name" size="small" maxlength="20" style="flex: 1" @blur="handleUpdateCategory(c)" />
            <el-tag v-if="c.isSystem" size="small" type="info">{{ t('common.system') }}</el-tag>
            <el-button v-else text size="small" type="danger" @click="handleDeleteCategory(c)">{{
              t('common.delete')
            }}</el-button>
          </div>
          <div v-if="!categories.length" class="empty-state" style="padding: 20px 0">
            <p>{{ t('home.noCategories') }}</p>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 通知权限请求弹窗 -->
    <el-dialog
      v-model="showNotificationPermission"
      :title="t('home.notificationPermission')"
      width="400px"
      destroy-on-close
    >
      <div style="text-align: center; padding: 20px 0">
        <div style="font-size: 48px; margin-bottom: 16px">🔔</div>
        <p style="margin-bottom: 8px">{{ t('home.notificationPermissionDesc') }}</p>
        <p style="color: #909399; font-size: 14px">{{ t('home.notificationPermissionNote') }}</p>
      </div>
      <template #footer>
        <el-button @click="showNotificationPermission = false">{{ t('common.later') }}</el-button>
        <el-button type="primary" @click="handleNotificationPermission">{{ t('common.allow') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Search, Plus, Check, Download } from '@element-plus/icons-vue'
import draggable from 'vuedraggable'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { resetAuthVerified } from '@/router'
import { useLocale } from '@/composables/useLocale'
import { useNotification } from '@/composables/useNotification'
import {
  getTaskList,
  createTask,
  updateTask,
  toggleTask,
  deleteTask,
  batchTask,
  sortTask,
  getCategoryList,
  createCategory,
  updateCategory,
  deleteCategory,
  getStat,
  changePassword,
  exportTasks,
} from '@/api'
import type { TaskItem, TaskFormData, StatResp, CategoryItem } from '@/types'

const { t } = useI18n()
const { currentLocale, setLocale, localeOptions } = useLocale()
const currentLang = ref(currentLocale.value)

const router = useRouter()
const userStore = useUserStore()

const { permission, requestPermission, scheduleReminders } = useNotification()
const showNotificationPermission = ref(false)

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
  status: undefined,
  categoryId: undefined,
  priority: undefined,
  keyword: '',
})

// 多选
const selectMode = ref(false)
const selectedIds = ref<number[]>([])

// 任务弹窗
const taskDialogVisible = ref(false)
const editingTask = ref<TaskItem | null>(null)
const submitting = ref(false)
const taskFormRef = ref<FormInstance>()
const taskForm = ref<TaskFormData>({
  title: '',
  content: '',
  priority: 3,
  categoryId: undefined,
  startTime: '',
  endTime: '',
  reminder: '',
  tags: '',
})
const taskRules = {
  title: [
    { required: true, message: () => t('home.enterTitle'), trigger: 'blur' },
    { max: 100, message: () => t('home.titleMaxLength'), trigger: 'blur' },
  ],
}

// 密码弹窗
const showPasswordDialog = ref(false)
const pwdLoading = ref(false)

// 分类管理弹窗
const showCategoryDialog = ref(false)
const newCategoryName = ref('')
const newCategoryColor = ref('#409eff')
const pwdFormRef = ref<FormInstance>()
const pwdForm = ref({ oldPassword: '', newPassword: '', confirmPassword: '' })
const validatePwdConfirm = (_rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (value !== pwdForm.value.newPassword) callback(new Error(t('auth.passwordMismatch')))
  else callback()
}
const pwdRules = {
  oldPassword: [{ required: true, message: () => t('auth.enterCurrentPassword'), trigger: 'blur' }],
  newPassword: [
    { required: true, message: () => t('auth.enterNewPassword'), trigger: 'blur' },
    { min: 6, max: 20, message: () => t('auth.passwordLength'), trigger: 'blur' },
    { pattern: /^(?=.*[a-zA-Z])(?=.*\d)/, message: () => t('auth.passwordComplexity'), trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: () => t('auth.confirmNewPasswordPrompt'), trigger: 'blur' },
    { validator: validatePwdConfirm, trigger: 'blur' },
  ],
}

function handleLocaleChange(lang: string) {
  setLocale(lang)
}

function onFilterChange() {
  page.value = 1
  loadTasks()
}

watch(page, () => loadTasks())

onMounted(() => {
  loadCategories()
  Promise.all([loadTasks(), loadStat()])
  if (permission.value === 'default') {
    showNotificationPermission.value = true
  }
})

async function loadStat() {
  try {
    stat.value = await getStat()
  } catch {
    ElMessage.error(t('home.loadStatFailed'))
  }
}

async function loadCategories() {
  try {
    const res = await getCategoryList()
    categories.value = (res.list || []).map((c) => ({ ...c, _name: c.name, _color: c.color || '#409eff' }))
  } catch {
    ElMessage.error(t('home.loadCategoriesFailed'))
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
    scheduleReminders(tasks.value.map((t) => ({ id: t.id, title: t.title, reminder: t.reminder })))
  } catch {
    ElMessage.error(t('home.loadTasksFailed'))
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
    ElMessage.success(task.status === 0 ? t('home.markedDone') : t('home.markedTodo'))
    await Promise.all([loadTasks(), loadStat()])
  } catch {
    ElMessage.error(t('home.toggleFailed'))
  }
}

async function handleDelete(id: number) {
  try {
    await deleteTask(id)
    ElMessage.success(t('home.deleted'))
    await Promise.all([loadTasks(), loadStat()])
  } catch {
    ElMessage.error(t('home.deleteTaskFailed'))
  }
}

async function handleBatch(action: string) {
  try {
    await batchTask({ ids: selectedIds.value, action })
    ElMessage.success(t('home.batchSuccess'))
    selectedIds.value = []
    await Promise.all([loadTasks(), loadStat()])
  } catch {
    ElMessage.error(t('home.batchFailed'))
  }
}

function openTaskDialog(task?: TaskItem) {
  editingTask.value = task || null
  if (task) {
    taskForm.value = {
      title: task.title,
      content: task.content || '',
      priority: task.priority,
      categoryId: task.categoryId || undefined,
      startTime: task.startTime || '',
      endTime: task.endTime || '',
      reminder: task.reminder || '',
      tags: task.tags || '',
    }
  } else {
    taskForm.value = {
      title: '',
      content: '',
      priority: 3,
      categoryId: undefined,
      startTime: '',
      endTime: '',
      reminder: '',
      tags: '',
    }
  }
  taskDialogVisible.value = true
}

function parseTags(tags: string): string[] {
  if (!tags) return []
  return tags
    .split(',')
    .map((t) => t.trim())
    .filter(Boolean)
}

function isOverdue(endTime: string, status: number): boolean {
  if (status === 2 || !endTime) return false
  return new Date(endTime) < new Date()
}

async function handleSubmitTask() {
  await taskFormRef.value?.validate()
  submitting.value = true
  try {
    if (editingTask.value) {
      await updateTask(editingTask.value.id, taskForm.value)
      ElMessage.success(t('home.updateSuccess'))
    } else {
      await createTask(taskForm.value)
      ElMessage.success(t('home.createSuccess'))
    }
    taskDialogVisible.value = false
    await Promise.all([loadTasks(), loadStat()])
  } catch {
    ElMessage.error(editingTask.value ? t('home.updateFailed') : t('home.createFailed'))
  } finally {
    submitting.value = false
  }
}

async function handleChangePassword() {
  await pwdFormRef.value?.validate()
  pwdLoading.value = true
  try {
    await changePassword({ oldPassword: pwdForm.value.oldPassword, newPassword: pwdForm.value.newPassword })
    ElMessage.success(t('auth.passwordChanged'))
    showPasswordDialog.value = false
    resetAuthVerified()
    userStore.logout()
    router.push('/login')
  } catch {
    ElMessage.error(t('auth.changePasswordFailed'))
  } finally {
    pwdLoading.value = false
  }
}

// 分类颜色辅助函数
function getCategoryColor(categoryId: number): string | undefined {
  const cat = categories.value.find((c) => c.id === categoryId)
  return cat?.color || undefined
}

function getCategoryTextColor(categoryId: number): string {
  const cat = categories.value.find((c) => c.id === categoryId)
  if (!cat?.color) return '#909399'
  // 简单亮度检测：浅色背景用深色文字
  const hex = cat.color.replace('#', '')
  const r = parseInt(hex.substring(0, 2), 16)
  const g = parseInt(hex.substring(2, 4), 16)
  const b = parseInt(hex.substring(4, 6), 16)
  const lum = (0.299 * r + 0.587 * g + 0.114 * b) / 255
  return lum > 0.6 ? '#303133' : '#ffffff'
}

async function handleAddCategory() {
  if (!newCategoryName.value.trim()) {
    ElMessage.warning(t('home.enterCategoryName'))
    return
  }
  try {
    await createCategory({ name: newCategoryName.value.trim(), color: newCategoryColor.value })
    newCategoryName.value = ''
    newCategoryColor.value = '#409eff'
    await loadCategories()
    ElMessage.success(t('home.categoryAdded'))
  } catch {
    ElMessage.error(t('home.addCategoryFailed'))
  }
}

async function handleUpdateCategory(c: CategoryItem & { _name: string; _color: string }) {
  if (!c._name.trim()) {
    ElMessage.warning(t('home.categoryNameEmpty'))
    c._name = c.name
    return
  }
  // 检查是否有变化
  if (c._name === c.name && c._color === c.color) return
  try {
    await updateCategory(c.id, { name: c._name.trim(), color: c._color })
    await loadCategories()
  } catch {
    ElMessage.error(t('home.updateCategoryFailed'))
    await loadCategories()
  }
}

async function handleDeleteCategory(c: CategoryItem) {
  if (c.isSystem) return
  try {
    await ElMessageBox.confirm(t('home.deleteCategoryConfirm', { name: c.name }), t('common.tip'), { type: 'warning' })
  } catch {
    return // 用户取消
  }
  try {
    await deleteCategory(c.id)
    await loadCategories()
    ElMessage.success(t('home.categoryDeleted'))
  } catch {
    ElMessage.error(t('home.deleteCategoryFailed'))
  }
}

async function onDragEnd() {
  const orders = tasks.value.map((t, index) => ({
    id: t.id,
    sortOrder: index + 1,
  }))
  try {
    await sortTask({ orders })
  } catch {
    ElMessage.error(t('home.sortSaveFailed'))
    loadTasks()
  }
}

function handleExport(format: string) {
  const params: Record<string, unknown> = { format }
  if (filters.value.status !== undefined && filters.value.status !== '') params.status = filters.value.status
  if (filters.value.categoryId) params.categoryId = filters.value.categoryId
  if (filters.value.priority) params.priority = filters.value.priority
  if (filters.value.keyword) params.keyword = filters.value.keyword
  exportTasks(params)
    .then((blob: unknown) => {
      const url = window.URL.createObjectURL(blob as Blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `tasks.${format}`
      a.click()
      window.URL.revokeObjectURL(url)
    })
    .catch(() => {
      ElMessage.error(t('common.exportFailed'))
    })
}

async function handleNotificationPermission() {
  await requestPermission()
  showNotificationPermission.value = false
}

function handleLogout() {
  ElMessageBox.confirm(t('common.logoutConfirm'), t('common.tip'), { type: 'warning' })
    .then(() => {
      resetAuthVerified()
      userStore.logout()
      router.push('/login')
    })
    .catch(() => {})
}
</script>

<style scoped>
.home-page {
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
  height: 60px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
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

.logo-icon {
  font-size: 24px;
}
.logo-text {
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

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
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  font-size: 24px;
  margin-bottom: 8px;
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

.stat-card.total .stat-value {
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.stat-card.todo .stat-value {
  background: linear-gradient(135deg, #e6a23c, #f5c242);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.stat-card.done .stat-value {
  background: linear-gradient(135deg, #67c23a, #38f9d7);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}
.stat-card.rate .stat-value {
  background: linear-gradient(135deg, #4facfe, #00f2fe);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.task-section {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
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
  background: linear-gradient(135deg, #303133, #667eea);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.section-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.batch-bar {
  background: rgba(102, 126, 234, 0.08);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(102, 126, 234, 0.15);
  padding: 10px 16px;
  border-radius: 12px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #667eea;
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
  border-radius: 12px;
  transition: all 0.3s ease;
}

.task-item.ghost {
  opacity: 0.4;
  background: rgba(102, 126, 234, 0.08);
}

.drag-handle {
  cursor: grab;
  color: #c0c4cc;
  font-size: 16px;
  padding: 4px 4px 0 0;
  user-select: none;
  line-height: 1;
}

.drag-handle:active {
  cursor: grabbing;
}

.task-item:hover {
  background: rgba(255, 255, 255, 0.6);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  transform: translateX(4px);
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
  transition: all 0.3s ease;
  cursor: pointer;
}

.status-circle:hover {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.15);
}

.status-circle.done {
  background: linear-gradient(135deg, #67c23a, #38f9d7);
  border-color: #67c23a;
  color: #fff;
  font-size: 12px;
}

.task-body {
  flex: 1;
  min-width: 0;
  cursor: pointer;
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

.task-time.overdue {
  color: #f56c6c;
  font-weight: 500;
}

.task-tag {
  border-radius: 12px;
  font-size: 11px;
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
  display: inline-block;
  opacity: 0.6;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}

.category-manage {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.category-add-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

@media (max-width: 768px) {
  .stat-row {
    grid-template-columns: repeat(2, 1fr);
  }
  .section-actions {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
