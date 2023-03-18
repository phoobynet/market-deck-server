import { defineStore } from 'pinia'
import { RealtimeSymbols } from '@/types'

export interface RealtimeSymbolsState {
  realtimeSymbols: RealtimeSymbols
}

export const useRealtimeSymbolsStore = defineStore('realtimeSymbols', {
  state: (): RealtimeSymbolsState => ({
    realtimeSymbols: {},
  }),

  getters: {
    symbols(state: RealtimeSymbolsState): string[] {
      if (!state.realtimeSymbols) return []

      return Object.keys(state.realtimeSymbols)
    }
  }
})
