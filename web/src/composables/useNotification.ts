import { ref, onMounted, onUnmounted } from 'vue'

export interface ReminderTask {
  id: number
  title: string
  reminder: string
}

export function useNotification() {
  const isSupported = ref('Notification' in window && 'serviceWorker' in navigator && 'PushManager' in window)
  const permission = ref<NotificationPermission>('default')
  const reminderTasks = ref<ReminderTask[]>([])
  let checkInterval: number | null = null

  async function requestPermission(): Promise<boolean> {
    if (!isSupported.value) return false

    const result = await Notification.requestPermission()
    permission.value = result
    return result === 'granted'
  }

  function showNotification(task: ReminderTask) {
    if (permission.value !== 'granted' || !isSupported.value) return

    const notification = new Notification('任务提醒', {
      body: task.title,
      icon: '/favicon.ico',
      tag: `task-${task.id}`,
      requireInteraction: true,
    })

    notification.addEventListener('click', () => {
      window.location.href = `/task/${task.id}`
      notification.close()
    })

    notification.addEventListener('close', () => {
      const index = reminderTasks.value.findIndex((t) => t.id === task.id)
      if (index > -1) {
        reminderTasks.value.splice(index, 1)
      }
    })
  }

  function scheduleReminders(tasks: ReminderTask[]) {
    reminderTasks.value = tasks.filter((t) => t.reminder)

    if (checkInterval) {
      clearInterval(checkInterval)
    }

    checkInterval = window.setInterval(() => {
      const now = new Date()
      const nowStr = now.toISOString().slice(0, 19).replace('T', ' ')

      reminderTasks.value.forEach((task) => {
        if (task.reminder && task.reminder <= nowStr) {
          showNotification(task)
          const index = reminderTasks.value.findIndex((t) => t.id === task.id)
          if (index > -1) {
            reminderTasks.value.splice(index, 1)
          }
        }
      })
    }, 60000)
  }

  function clearReminders() {
    if (checkInterval) {
      clearInterval(checkInterval)
      checkInterval = null
    }
    reminderTasks.value = []
  }

  onMounted(() => {
    if (isSupported.value) {
      permission.value = Notification.permission
    }
  })

  onUnmounted(() => {
    clearReminders()
  })

  return {
    isSupported,
    permission,
    reminderTasks,
    requestPermission,
    showNotification,
    scheduleReminders,
    clearReminders,
  }
}
