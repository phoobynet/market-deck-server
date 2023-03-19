import { Bar } from '@/types/Bar'
import { Trade } from '@/types/Trade'
import { Quote } from '@/types/Quote'
import { Asset } from '@/types/Asset'

export interface Snapshot {
  a: Asset
  lb: Bar
  lt: Trade
  lq: Quote
}
