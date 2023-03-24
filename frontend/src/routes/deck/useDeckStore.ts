import { defineStore } from 'pinia'
import { Snapshot } from '@/types/Snapshot'
import { http } from '@/libs/http'

export interface DeckState {
  snapshots: Record<string, Snapshot>
  showModal: boolean
}

export const useDeckStore = defineStore('deck', {
  state: (): DeckState => ({
    snapshots: {},
    showModal: false,
  }),

  actions: {
    async getSymbols () {
      return http.get<{ symbols: string[] }>('/symbols').then((res) => res.data.symbols)
    },
    async updateSymbols (symbols: string[] = []) {
      await http.post(`/symbols?symbols=${symbols.join(',')}`, symbols)
    },
    async deleteSymbol (symbol: string) {
      const newSymbols = Object.keys(this.snapshots).filter((key) => key !== symbol)
      await http.post(`/symbols?symbols=${newSymbols.join(',')}`, newSymbols)
    }
  },
  getters: {
    snapshotsList: (state) => Object.values(state.snapshots),
    symbols: (state) => Object.keys(state.snapshots),
  },
})

const source = new EventSource('http://localhost:3000/api/stream?stream=snapshots')

source.onmessage = (event) => {
  const { data } = JSON.parse(event.data) as { event: string, data: Record<string, Snapshot> }

  useDeckStore().snapshots = data
}
