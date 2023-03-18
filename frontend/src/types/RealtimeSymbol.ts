import { Bar } from '@/types/Bar'
import { Trade } from '@/types/Trade'
import { Quote } from '@/types/Quote'
import { Asset } from '@/types/Asset'

export interface RealtimeSymbol {
  asset: Asset
  bar: Bar
  trade: Trade
  quote: Quote
  prevDailyBar: Bar
  intradayBars: Bar[]
  dailyBars: Bar[]
}
