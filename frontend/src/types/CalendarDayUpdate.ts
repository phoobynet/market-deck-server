import { CalendarDay } from './CalendarDay'
import { CurrentMarketCondition } from '@/types/CurrentMarketCondition'

export interface CalendarDayUpdate {
  condition: CurrentMarketCondition
  at: number
  previous: CalendarDay
  current: CalendarDay
  next: CalendarDay
}
