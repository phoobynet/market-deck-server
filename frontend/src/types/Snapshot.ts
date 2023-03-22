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
  changes: Record<string, SnapshotChange>
  ibars: Bar[]
  // intraday high
  ih: number

  // intraday low
  il: number
}

export interface SnapshotChange {
  since: number
  label: string
  c: number
  // change percent
  cp: number
  // change sign
  cs: number
  // change absolute
  ca: number
}
