import { baseUrl } from '@/libs/baseUrl'
import { useCalendarDayUpdate } from '@/stores'

export class CalendarDayUpdateListener {
  private calendarDayUpdateSource?: EventSource

  start () {
    this.calendarDayUpdateSource = new EventSource(`${baseUrl}/api/stream?stream=calendar_day_update`)

    const calendarDayUpdateStore = useCalendarDayUpdate()

    this.calendarDayUpdateSource.onmessage = (event) => {
      calendarDayUpdateStore.calendarDayUpdate = JSON.parse(event.data)
    }

    this.calendarDayUpdateSource.onerror = (error) => {
      console.error(error)
    }
  }

  close () {
    this.calendarDayUpdateSource?.close()
  }
}
