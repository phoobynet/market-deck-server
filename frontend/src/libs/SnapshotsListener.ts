import { baseUrl } from '@/libs/baseUrl'
import { useSnapshotsStore } from '@/stores'
import { Snapshots } from '@/types'

export class SnapshotsListener {
  private source?: EventSource

  start () {
    this.source = new EventSource(`${baseUrl}/api/stream?stream=snapshots`)

    const snapshotsStore = useSnapshotsStore()

    this.source.onmessage = (event) => {
      const { data } = JSON.parse(event.data) as { event: string, data: Snapshots }

      snapshotsStore.snapshots = data
    }

    this.source.onerror = (error) => {
      console.error(error)
    }
  }

  close () {
    this.source?.close()
  }
}
