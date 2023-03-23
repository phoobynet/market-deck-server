import { defineStore } from 'pinia'
import type { Snapshots } from '@/types'

export interface SnapshotsState {
  snapshots: Snapshots
}

export const useSnapshotsStore = defineStore('snapshots', {
  state: (): SnapshotsState => ({
    snapshots: {},
  }),

  getters: {
    symbols (state: SnapshotsState): string[] {
      if (!state.snapshots) return []

      return Object.keys(state.snapshots)
    },
  },
})

