import type { SnapshotLiteChange } from './SnapshotLiteChange'
import { Bar } from '@/types/Bar'

export interface SnapshotLite {
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
  change: SnapshotLiteChange
  volumes: Array<{
    date: string,
    vol: number
  }>
  monthlyBars: Array<Bar>
  ytdChange: SnapshotLiteChange
}
