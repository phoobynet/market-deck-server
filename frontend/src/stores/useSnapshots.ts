import { defineStore } from 'pinia'
import { Snapshots } from '@/types'

export interface SnapshotsState {
  snapshots: Snapshots
}

export const useSnapshots = defineStore('snapshots', {
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
