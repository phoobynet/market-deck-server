import { Bar } from '@/types/Bar'
import { Trade } from '@/types/Trade'
import { Quote } from '@/types/Quote'

export interface Snapshot {
  lb: Bar
  lt: Trade
  lq: Quote
}
