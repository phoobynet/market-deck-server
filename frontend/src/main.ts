import { createApp } from 'vue'
import App from '@/App.vue'
import { router } from '@/router'
import { createPinia } from 'pinia'
import '@/styles/styles.scss'
import 'animate.css/animate.css'

const pinia = createPinia()

const app = createApp(App)

// special requirement for built-in components
app.use(router)
app.use(pinia)
app.mount('#app')
