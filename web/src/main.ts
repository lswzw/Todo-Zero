import { createApp, defineComponent, h } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import i18n from './locales'
import { useLocale } from './composables/useLocale'
import ripple from './directives/ripple'
import './styles/main.css'
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/message-box/style/css'

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

app.directive('ripple', ripple)
app.use(createPinia())
app.use(i18n)
app.use(router)
app.mount('#app')
