import { baseUrl } from '@/libs/baseUrl'
import { RealtimeSymbols } from '@/types'
import { useRealtimeSymbolsStore } from '@/stores'

export class RealtimeSymbolsListener {
  private realtimeSymbolsSource?: EventSource

  start () {
    this.realtimeSymbolsSource = new EventSource(`${baseUrl}/api/stream?stream=realtime_symbols`)

    const realtimeSymbolsStore = useRealtimeSymbolsStore()

    this.realtimeSymbolsSource.onmessage = (event) => {
      const { data } = JSON.parse(event.data) as { data: RealtimeSymbols }

      realtimeSymbolsStore.$patch({ realtimeSymbols: data })
    }

    this.realtimeSymbolsSource.onerror = (error) => {
      console.error(error)
    }
  }

  close () {
    this.realtimeSymbolsSource?.close()
  }
}
