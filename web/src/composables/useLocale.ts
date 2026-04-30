import { computed, ref, readonly } from 'vue'
import { useI18n } from 'vue-i18n'
import { useStorage } from './useStorage'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'

// 支持的语言列表，未来添加新语言只需在此扩展
export const LOCALE_OPTIONS = [
  { label: '中文', value: 'zh-CN', alias: '中' },
  { label: 'English', value: 'en', alias: 'En' },
] as const

export type LocaleValue = (typeof LOCALE_OPTIONS)[number]['value']

const STORAGE_KEY = 'locale'

// 全局单例：持久化存储
const { value: storedLocale } = useStorage<string>(STORAGE_KEY, 'zh-CN')

// 全局单例：响应式 locale，确保所有组件共享同一状态
const syncedLocale = ref(storedLocale.value)

// 监听其他标签页的语言切换
if (typeof window !== 'undefined') {
  window.addEventListener('storage', (e) => {
    if (e.key === STORAGE_KEY && e.newValue) {
      syncedLocale.value = e.newValue
    }
  })
}

/**
 * 语言切换核心方法
 * 统一管理语言状态、持久化存储、Element Plus 语言同步、HTML lang 属性更新
 */
export function useLocale() {
  const { locale, t } = useI18n()

  // 初始化同步：确保 vue-i18n 与全局单例一致
  if (syncedLocale.value !== locale.value) {
    locale.value = syncedLocale.value
  }

  const elementLocale = computed(() => (locale.value === 'en' ? en : zhCn))

  const currentLocale = computed(() => locale.value)

  const currentLocaleOption = computed(
    () => LOCALE_OPTIONS.find((opt) => opt.value === locale.value) || LOCALE_OPTIONS[0],
  )

  /**
   * 切换语言
   * @param lang 目标语言代码
   */
  function setLocale(lang: string) {
    if (lang === locale.value) return
    locale.value = lang
    syncedLocale.value = lang
    storedLocale.value = lang
    document.documentElement.lang = lang === 'zh-CN' ? 'zh-CN' : 'en'
  }

  /**
   * 切换到下一个语言（循环切换）
   */
  function toggleLocale() {
    const currentIndex = LOCALE_OPTIONS.findIndex((opt) => opt.value === locale.value)
    const nextIndex = (currentIndex + 1) % LOCALE_OPTIONS.length
    setLocale(LOCALE_OPTIONS[nextIndex].value)
  }

  /**
   * 获取当前语言的显示标签
   */
  function getLocaleLabel(lang: LocaleValue): string {
    return LOCALE_OPTIONS.find((opt) => opt.value === lang)?.label || lang
  }

  /**
   * 获取当前语言的别名（简短显示）
   */
  function getLocaleAlias(lang: LocaleValue): string {
    return LOCALE_OPTIONS.find((opt) => opt.value === lang)?.alias || lang
  }

  return {
    locale: readonly(locale),
    currentLocale,
    currentLocaleOption,
    currentLocaleLabel: computed(() => currentLocaleOption.value.label),
    currentLocaleAlias: computed(() => currentLocaleOption.value.alias),
    elementLocale,
    localeOptions: LOCALE_OPTIONS,
    setLocale,
    toggleLocale,
    getLocaleLabel,
    getLocaleAlias,
    t,
  }
}
