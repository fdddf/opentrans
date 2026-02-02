import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/style.css'
import { i18n } from './i18n'
import Toast, { type PluginOptions } from 'vue-toastification'
import 'vue-toastification/dist/index.css'

const app = createApp(App)

// Toast plugin configuration
const toastOptions: PluginOptions = {
  timeout: 3000,
  position: 'top-right',
  closeOnClick: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: false,
  icon: true,
  rtl: false,
  toastClassName: 'vue-toastification-custom',
}

app.use(router)
app.use(i18n)
app.use(Toast, toastOptions)
app.mount('#app')
