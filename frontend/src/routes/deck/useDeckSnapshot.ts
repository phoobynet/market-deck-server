import { useDeckStore } from '@/routes/deck/useDeckStore'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { SnapshotLite } from '@/types/SnapshotLite'
import { SnapshotLiteChange } from '@/types/SnapshotLiteChange'
import { formatPercentAbs } from '@/libs/helpers/formatPercent'
import { formatMoneyNoSymbol } from '@/libs/helpers/formatMoney'
import numeral from 'numeral'
import { sumBy } from 'lodash'

export const useDeckSnapshot = (symbol: string) => {
  const deckStore = useDeckStore()

  const { snapshots } = storeToRefs(deckStore)

  const snapshot = computed<SnapshotLite | undefined>(() => snapshots.value[symbol])

  const change = computed<SnapshotLiteChange | undefined>(() => snapshot.value?.change)
  const ytdChange = computed<SnapshotLiteChange | undefined>(() => snapshot.value?.ytdChange)

  const sign = computed(() => change.value?.sign ?? '')

  const changePercentAbs = computed(() =>
    change.value?.changePercent
      ? `${formatPercentAbs(Math.abs(change.value?.changePercent ?? 0))}`
      : '')

  const priceChange = computed(() =>
    change.value?.change
      ? `${sign.value}${formatMoneyNoSymbol(change.value?.absoluteChange)}`
      : '')

  const ytdChangePercentAbs = computed(() =>
    ytdChange.value?.changePercent
      ? `${formatPercentAbs(Math.abs(ytdChange.value?.changePercent ?? 0))}`
      : '')

  const ytdPriceChange = computed(() =>
    ytdChange.value?.change
      ? `${sign.value}${formatMoneyNoSymbol(ytdChange.value?.absoluteChange)}`
      : '')

  const prevClose = computed(() => formatMoneyNoSymbol(snapshot?.value?.prevClose))

  const price = computed(() => formatMoneyNoSymbol(snapshot?.value?.price))

  const dailyHigh = computed(() => formatMoneyNoSymbol(snapshot?.value?.dailyHigh))
  const dailyLow = computed(() => formatMoneyNoSymbol(snapshot?.value?.dailyLow))
  const dailyVolume = computed(() => numeral(snapshot?.value?.dailyVolume).format('0,0a'))

  const avgVolume = computed<string>(() => {
    const volumes = snapshot?.value?.volumes ?? []

    if (volumes.length === 0) {
      return '?'
    }

    return numeral(sumBy(volumes, 'vol') / volumes.length).format('0,0a')
  })

  return {
    sign,
    changePercentAbs,
    priceChange,
    ytdChangePercentAbs,
    ytdPriceChange,
    prevClose,
    price,
    snapshot,
    dailyHigh,
    dailyLow,
    dailyVolume,
    avgVolume,
  }
}
