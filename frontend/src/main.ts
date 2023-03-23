import { createApp } from 'vue'
import App from '@/App.vue'
import { router } from '@/router'
import { createPinia } from 'pinia'
import '@/styles/styles.scss'
import 'animate.css/animate.css'

import DeckTableRow from '@/routes/deck/DeckTableRow.vue'

const pinia = createPinia()

const app = createApp(App)

// special requirement for built-in components
app.component('deck-table-row', DeckTableRow)
app.use(router)
app.use(pinia)
app.mount('#app')
