import { CalendarDay } from './CalendarDay'
import { CurrentMarketCondition } from '@/types/CurrentMarketCondition'

export interface CalendarDayUpdate {
  condition: CurrentMarketCondition
  at: number
  prev: CalendarDay
  current: CalendarDay
  next: CalendarDay
}
