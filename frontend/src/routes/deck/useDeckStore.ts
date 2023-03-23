import { defineStore } from 'pinia'
import { SnapshotLite } from '@/types/SnapshotLite'
import { http } from '@/libs/http'

export interface DeckState {
  snapshots: Record<string, SnapshotLite>
}

export const useDeckStore = defineStore('deck', {
  state: (): DeckState => ({
    snapshots: {},
  }),

  actions: {
    async getSymbols () {
      return http.get<{ symbols: string[] }>('/symbols').then((res) => res.data.symbols)
    },
    async updateSymbols (symbols: string[] = []) {
      await http.post(`/symbols?symbols=${symbols.join(',')}`, symbols)
    },
  },
  getters: {
    snapshotsList: (state) => Object.values(state.snapshots),
    symbols: (state) => Object.keys(state.snapshots),
  },
})

const source = new EventSource('http://localhost:3000/api/stream?stream=snapshots_lite')

source.onmessage = (event) => {
  const { data } = JSON.parse(event.data) as { event: string, data: Record<string, SnapshotLite> }

  useDeckStore().snapshots = data
}
