import { ref, watch } from 'vue'

export function useStorage<T>(
  key: string,
  defaultValue: T,
): {
  value: ReturnType<typeof ref<T>>
  remove: () => void
} {
  const storedValue = localStorage.getItem(key)
  const value = ref<T>(storedValue !== null ? (parseStorageValue(storedValue) as T) : defaultValue)

  watch(
    value,
    (newValue) => {
      if (newValue === null || newValue === undefined) {
        localStorage.removeItem(key)
      } else {
        localStorage.setItem(key, stringifyStorageValue(newValue))
      }
    },
    { deep: true },
  )

  function remove() {
    localStorage.removeItem(key)
    value.value = defaultValue
  }

  return { value, remove }
}

function stringifyStorageValue(value: unknown): string {
  if (typeof value === 'string') {
    return value
  }
  if (typeof value === 'number' || typeof value === 'boolean') {
    return String(value)
  }
  return JSON.stringify(value)
}

function parseStorageValue(value: string): unknown {
  if (value === 'true') return true
  if (value === 'false') return false
  if (!isNaN(Number(value))) return Number(value)
  try {
    return JSON.parse(value)
  } catch {
    return value
  }
}
