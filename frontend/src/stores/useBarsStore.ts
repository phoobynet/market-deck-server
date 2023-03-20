import { Bar } from '@/types'
import { defineStore } from 'pinia'

export interface BarState {
  bars: Record<string, Bar[]>
}

export const useBarsStore = defineStore('bars', {
  state: () => ({
    bars: {},
  }),
  actions: {
    async fetch(symbol: string): Promise<void> {

    }
  }
})
