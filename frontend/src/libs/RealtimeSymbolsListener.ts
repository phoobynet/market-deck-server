import { baseUrl } from '@/libs/baseUrl'
import { useRealtimeSymbolsStore } from '@/stores'

export class RealtimeSymbolsListener {
  private realtimeSymbolsSource?: EventSource

  start () {
    this.realtimeSymbolsSource = new EventSource(`${baseUrl}/api/stream?stream=realtime_symbols`)

    const realtimeSymbolsStore = useRealtimeSymbolsStore()

    this.realtimeSymbolsSource.onmessage = (event) => {
      realtimeSymbolsStore.realtimeSymbols = JSON.parse(event.data)
    }

    this.realtimeSymbolsSource.onerror = (error) => {
      console.error(error)
    }
  }

  close () {
    this.realtimeSymbolsSource?.close()
  }
}
