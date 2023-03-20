import { Bar } from '@/types/Bar'
import { Trade } from '@/types/Trade'
import { Quote } from '@/types/Quote'

export interface Snapshot {
  // latest bar
  lb: Bar
  // latest trade
  lt: Trade
  // latest quote
  lq: Quote
  // Previous daily bar
  pdb: Bar
  // Daily bar
  db: Bar
  // previous close
  pc: number
  // change amount
  c: number
  // change percent
  cp: number
  // change sign
  cs: number
  // change absolute
  ca: number
}
