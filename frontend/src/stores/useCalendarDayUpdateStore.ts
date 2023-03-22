import { defineStore } from 'pinia'
import { CalendarDayUpdate, CurrentMarketCondition } from '@/types'
import { formatInTimeZone } from 'date-fns-tz'


const marketTimeZone = 'America/New_York'

export interface CalendarDayUpdateState {
  calendarDayUpdate: CalendarDayUpdate | undefined
}

export const useCalendarDayUpdateStore = defineStore('calendarDayUpdate', ({
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
    marketTimeFormatted: (state): string => {
      if (state.calendarDayUpdate === undefined) {
        return ''
      }

      return formatInTimeZone(new Date(state.calendarDayUpdate.at), marketTimeZone, 'eee eo MMM HH:mm:ss zzz')
    },
    condition: (state): CurrentMarketCondition => {
      if (state.calendarDayUpdate === undefined) {
        return CurrentMarketCondition.unknown
      }

      return state.calendarDayUpdate.condition
    },

    previousDate: (state): string => {
      if (state.calendarDayUpdate === undefined) {
        return ''
      }

      return state.calendarDayUpdate?.prev?.date ?? ''
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
          return 'After hours'
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
