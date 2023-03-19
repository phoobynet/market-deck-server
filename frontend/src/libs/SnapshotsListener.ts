import { baseUrl } from '@/libs/baseUrl'
import { useSnapshots } from '@/stores'

export class SnapshotsListener {
  private source?: EventSource

  start () {
    this.source = new EventSource(`${baseUrl}/api/stream?stream=realtime_symbols`)

    const snapshotsStore = useSnapshots()

    this.source.onmessage = (event) => {
      snapshotsStore.snapshots = JSON.parse(event.data)
    }

    this.source.onerror = (error) => {
      console.error(error)
    }
  }

  close () {
    this.source?.close()
  }
}
