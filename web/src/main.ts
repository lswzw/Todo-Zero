import { createApp, defineComponent, h } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import i18n from './locales'
import { useLocale } from './composables/useLocale'
import './styles/main.css'

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

app.use(createPinia())
app.use(i18n)
app.use(router)
app.mount('#app')
