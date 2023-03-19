import { defineStore } from 'pinia'
import { Asset } from '@/types'
import { http } from '@/libs/http'

export interface AssetsState {
  assets: Asset[]
  assetsMap: Record<string, Asset>
  fetching: boolean
}

export const useAssetsStore = defineStore('assets', {
  state: (): AssetsState => ({
    assets: [],
    assetsMap: {},
    fetching: true,
  }),
  actions: {
    async fetch (): Promise<void> {
      this.fetching = true
      try {
        const assets = await http.get<Asset[]>('/assets').then(r => r.data)

        this.assetsMap = assets.reduce((acc, asset) => {
          return {
            ...acc,
            [asset.S]: asset,
          }
        }, {})

        this.assets = assets
      } finally {
        this.fetching = false
      }
    },

    getBySymbol (symbol: string): Asset | undefined {
      return this.assetsMap[symbol]
    },
  },
  getters: {
    hasAssets (state): boolean {
      return state.assets.length > 0
    },
  },
})
