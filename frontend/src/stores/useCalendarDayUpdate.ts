import { defineStore } from 'pinia'
import { CalendarDayUpdate, CurrentMarketCondition } from '@/types'
import { format, utcToZonedTime } from 'date-fns-tz'

const TIME_FORMAT = 'E do LLL HH:mm:ss'
const USER_TIMEZONE = Intl.DateTimeFormat().resolvedOptions().timeZone
const marketTimeZone = 'America/New_York'

export interface CalendarDayUpdateState {
  calendarDayUpdate: CalendarDayUpdate | undefined
}

export const useCalendarDayUpdate = defineStore('calendarDayUpdate', ({
  state: (): CalendarDayUpdateState => ({
    calendarDayUpdate: undefined,
  }),

  getters: {
    timeUnixMs: (state): number => {
      if (state.calendarDayUpdate === undefined) {
        return 0
      }

      return state.calendarDayUpdate.at
    },
    marketTimeZone: (): string => marketTimeZone,
    localTimeZone: (): string => USER_TIMEZONE,
    marketTimeFormatted: (state): string => {
      if (state.calendarDayUpdate === undefined) {
        return ''
      }

      return format(utcToZonedTime(new Date(state.calendarDayUpdate.at), marketTimeZone), TIME_FORMAT)
    },
    localTimeFormatted: (state): string => {
      if (state.calendarDayUpdate === undefined) {
        return ''
      }

      return format(utcToZonedTime(new Date(state.calendarDayUpdate.at), USER_TIMEZONE), TIME_FORMAT)
    },
    condition: (state): CurrentMarketCondition => {
      if (state.calendarDayUpdate === undefined) {
        return CurrentMarketCondition.unknown
      }

      return state.calendarDayUpdate.condition
    },

    conditionDescription: (state): string => {
      let c
      if (state.calendarDayUpdate === undefined) {
        c = CurrentMarketCondition.unknown
      } else {
        c = state.calendarDayUpdate.condition
      }

      switch (c) {
        case CurrentMarketCondition.closed_today:
          return 'closed today'
        case CurrentMarketCondition.open:
          return 'open'
        case CurrentMarketCondition.pre_market:
          return 'Pre-market'
        case CurrentMarketCondition.post_market:
          return 'Post-market'
        case CurrentMarketCondition.closed_opening_later:
          return 'Opening later'
        case CurrentMarketCondition.closed_for_the_day:
          return 'Closed for today'
        default:
          return 'unknown'
      }
    },
  },
}))
