import { baseUrl } from '@/libs/baseUrl'
import { useCalendarDayUpdateStore } from '@/stores'
import { CalendarDayUpdate } from '@/types'

export class CalendarDayUpdateListener {
  private calendarDayUpdateSource?: EventSource

  start () {
    this.calendarDayUpdateSource = new EventSource(`${baseUrl}/api/stream?stream=calendar_day_update`)

    const calendarDayUpdateStore = useCalendarDayUpdateStore()

    this.calendarDayUpdateSource.onmessage = (event) => {
      const { data } = JSON.parse(event.data) as { event: string, data: CalendarDayUpdate }

      calendarDayUpdateStore.calendarDayUpdate = data
    }

    this.calendarDayUpdateSource.onerror = (error) => {
      console.error(error)
    }
  }

  close () {
    this.calendarDayUpdateSource?.close()
  }
}
