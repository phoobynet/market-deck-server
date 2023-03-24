import type { SnapshotChange } from './SnapshotChange'
import { Bar } from '@/types/Bar'

export interface Snapshot {
  class: string
  symbol: string
  name: string
  exchange: string
  price: number
  prevClose: number
  prevCloseDate: string
  dailyHigh: number
  dailyLow: number
  dailyVolume: number
  change: SnapshotChange
  volumes: Array<{
    date: string,
    vol: number
  }>
  monthlyBars: Array<Bar>
  ytdChange: SnapshotChange
}
