import { createRouter, createWebHashHistory } from 'vue-router'
import Dashboard from '@/routes/dashboard/Dashboard.vue'

export const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      component: Dashboard,
      name: 'Dashboard',
    },
  ],
})
