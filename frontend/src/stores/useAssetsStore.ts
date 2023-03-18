import { defineStore } from 'pinia'
import { Asset } from '@/types'
import { http } from '@/libs/http'

export interface AssetsState {
  assets: Asset[]
  fetching: boolean
}

export const useAssetsStore = defineStore('assets', {
  state: (): AssetsState => ({
    assets: [],
    fetching: true,
  }),
  actions: {
    async fetch (): Promise<void> {
      this.fetching = true
      try {
        this.assets = await http.get<Asset[]>('/assets').then(r => r.data)
      } finally {
        this.fetching = false
      }
    },
  },
  getters: {
    hasAssets (state): boolean {
      return state.assets.length > 0
    },
  },
})
