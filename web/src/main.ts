import { createApp, defineComponent, h } from 'vue'
import ElementPlus from 'element-plus'
import { ElConfigProvider } from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import i18n from './locales'
import { useLocale } from './composables/useLocale'
import './style.css'

// 用渲染函数包装 el-config-provider，避免修改 App.vue
const RootApp = defineComponent({
  name: 'RootApp',
  setup() {
    const { elementLocale } = useLocale()
    return { elementLocale }
  },
  render() {
    return h(ElConfigProvider, { locale: this.elementLocale }, () => h(App))
  },
})

const app = createApp(RootApp)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(ElementPlus)
app.use(createPinia())
app.use(i18n)
app.use(router)
app.mount('#app')
