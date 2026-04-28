import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'

export function useLocale() {
  const { locale } = useI18n()

  const elementLocale = computed(() => (locale.value === 'en' ? en : zhCn))

  function setLocale(lang: string) {
    locale.value = lang
    localStorage.setItem('locale', lang)
    document.documentElement.lang = lang === 'zh-CN' ? 'zh-CN' : 'en'
  }

  const currentLocale = computed(() => locale.value)

  const localeOptions = [
    { label: '中文', value: 'zh-CN' },
    { label: 'English', value: 'en' },
  ]

  return { locale, elementLocale, setLocale, currentLocale, localeOptions }
}
