import { createRouter, createWebHashHistory } from 'vue-router'
import Deck from '@/routes/deck/Deck.vue'

export const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      component: Deck,
      name: 'Deck',
    },
  ],
})
