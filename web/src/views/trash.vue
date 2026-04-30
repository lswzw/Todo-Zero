<template>
  <div class="trash-page">
    <!-- 顶部导航 -->
    <header class="navbar">
      <div class="nav-inner">
        <div class="nav-left">
          <el-button text @click="router.push('/')">
            <el-icon><ArrowLeft /></el-icon> {{ t('trash.back') }}
          </el-button>
          <span class="page-title">{{ t('trash.trash') }}</span>
        </div>
        <div class="nav-right">
          <el-button
            v-if="selectMode && selectedIds.length > 0"
            type="success"
            size="small"
            @click="handleBatchRestore"
          >
            {{ t('trash.batchRestore', { count: selectedIds.length }) }}
          </el-button>
          <el-button
            v-if="selectMode && selectedIds.length > 0"
            type="danger"
            size="small"
            @click="handleBatchPermanentDelete"
          >
            {{ t('trash.batchPermanentDelete', { count: selectedIds.length }) }}
          </el-button>
          <el-button :type="selectMode ? 'primary' : ''" size="small" @click="toggleSelectMode">
            {{ selectMode ? t('home.exitMultiSelect') : t('home.multiSelect') }}
          </el-button>
        </div>
      </div>
    </header>

    <!-- 主内容 -->
    <main class="main-content">
      <div class="trash-section">
        <div v-if="tasks.length" class="task-list">
          <div v-for="task in tasks" :key="task.id" class="task-item">
            <div v-if="selectMode" class="task-check" @click="toggleSelect(task.id)">
              <div :class="['check-dot', { active: selectedIds.includes(task.id) }]" />
            </div>
            <div class="task-body">
              <div class="task-title">{{ task.title }}</div>
              <div v-if="task.content" class="task-content">{{ task.content }}</div>
              <div class="task-meta">
                <el-tag v-if="task.priority === 2" size="small" type="danger">{{ t('home.urgent') }}</el-tag>
                <el-tag v-else-if="task.priority === 1" size="small" type="warning">{{ t('home.important') }}</el-tag>
                <el-tag v-else size="small" type="success">{{ t('home.normal') }}</el-tag>
                <el-tag
                  v-if="task.categoryName && task.categoryName !== t('trash.uncategorized')"
                  size="small"
                  type="info"
                  >{{ task.categoryName }}</el-tag
                >
                <span class="task-time">{{ t('trash.deletedAt', { time: task.updateTime }) }}</span>
              </div>
            </div>
            <div class="task-actions">
              <el-button text size="small" type="success" @click="handleRestore(task.id)">{{
                t('trash.restore')
              }}</el-button>
              <el-popconfirm :title="t('trash.permanentDeleteConfirm')" @confirm="handlePermanentDelete(task.id)">
                <template #reference>
                  <el-button text size="small" type="danger">{{ t('trash.permanentDelete') }}</el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">
          <span class="empty-icon">🗑️</span>
          <p>{{ t('trash.empty') }}</p>
        </div>

        <!-- 分页 -->
        <div v-if="total > 0" class="pagination">
          <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="prev, pager, next" />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { getTrashList, restoreTask, permanentDeleteTask, batchTask } from '@/api'
import type { TrashItem } from '@/types'

const { t } = useI18n()

const router = useRouter()
const tasks = ref<TrashItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 多选
const selectMode = ref(false)
const selectedIds = ref<number[]>([])

watch(page, () => loadTasks())

onMounted(() => loadTasks())

async function loadTasks() {
  try {
    const res = await getTrashList({ page: page.value, pageSize: pageSize.value })
    tasks.value = res.list || []
    total.value = res.total || 0
  } catch {
    ElMessage.error(t('trash.loadFailed'))
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

async function handleRestore(id: number) {
  try {
    await restoreTask(id)
    ElMessage.success(t('trash.restored'))
    loadTasks()
  } catch {
    ElMessage.error(t('trash.restoreFailed'))
  }
}

async function handlePermanentDelete(id: number) {
  try {
    await permanentDeleteTask(id)
    ElMessage.success(t('trash.permanentDeleted'))
    loadTasks()
  } catch {
    ElMessage.error(t('trash.permanentDeleteFailed'))
  }
}

async function handleBatchRestore() {
  try {
    await batchTask({ ids: selectedIds.value, action: 'restore' })
    ElMessage.success(t('trash.batchRestoreSuccess'))
    selectedIds.value = []
    loadTasks()
  } catch {
    ElMessage.error(t('trash.batchRestoreFailed'))
  }
}

async function handleBatchPermanentDelete() {
  try {
    await ElMessageBox.confirm(t('trash.batchPermanentDeleteConfirm'), t('common.tip'), { type: 'warning' })
  } catch {
    return
  }
  try {
    // 并行永久删除
    const results = await Promise.allSettled(selectedIds.value.map((id) => permanentDeleteTask(id)))
    const failCount = results.filter((r) => r.status === 'rejected').length
    if (failCount > 0) {
      ElMessage.warning(t('trash.someDeleteFailed', { count: failCount }))
    } else {
      ElMessage.success(t('trash.permanentDeleted'))
    }
    selectedIds.value = []
    loadTasks()
  } catch {
    ElMessage.error(t('trash.batchPermanentDeleteFailed'))
  }
}
</script>

<style scoped>
.trash-page {
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

.page-title {
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
  gap: 8px;
}

.main-content {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px 20px;
}

.trash-section {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
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

.task-item:hover {
  background: rgba(255, 255, 255, 0.6);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
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

.task-body {
  flex: 1;
  min-width: 0;
}

.task-title {
  font-size: 15px;
  color: #909399;
  line-height: 1.5;
  text-decoration: line-through;
}

.task-content {
  font-size: 13px;
  color: #c0c4cc;
  margin-top: 4px;
  line-height: 1.4;
  text-decoration: line-through;
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
  display: inline-block;
  opacity: 0.6;
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

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
