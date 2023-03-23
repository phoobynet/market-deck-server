import { createRouter, createWebHashHistory } from 'vue-router'
import Dashboard from '@/routes/dashboard/Dashboard.vue'
import Deck from '@/routes/deck/Deck.vue'

export const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/dashboard',
      component: Dashboard,
      name: 'Dashboard',
    },
    {
      path: '/',
      component: Deck,
      name: 'Deck',
    },
  ],
})
