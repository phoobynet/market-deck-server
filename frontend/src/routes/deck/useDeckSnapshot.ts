import { useDeckStore } from '@/routes/deck/useDeckStore'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { Snapshot } from '@/types/Snapshot'
import { SnapshotChange } from '@/types/SnapshotChange'
import { formatPercentAbs } from '@/libs/helpers/formatPercent'
import { formatMoneyNoSymbol } from '@/libs/helpers/formatMoney'
import numeral from 'numeral'
import { sumBy } from 'lodash'
import { cleanAssetName } from '@/libs/helpers/cleanAssetName'

export const useDeckSnapshot = (symbol: string) => {
  const deckStore = useDeckStore()

  const { snapshots } = storeToRefs(deckStore)

  const snapshot = computed<Snapshot | undefined>(() => snapshots.value[symbol])
  const companyName = computed<string>(() => cleanAssetName(snapshot.value?.name))

  const change = computed<SnapshotChange | undefined>(() => snapshot.value?.change)
  const ytdChange = computed<SnapshotChange | undefined>(() => snapshot.value?.ytdChange)

  const sign = computed(() => change.value?.sign ?? '')

  const changePercentAbs = computed(() =>
    change.value?.changePercent !== undefined
      ? `${formatPercentAbs(Math.abs(change.value?.changePercent ?? 0))}`
      : '')

  const priceChange = computed(() =>
    change.value?.change !== undefined
      ? `${sign.value}${formatMoneyNoSymbol(change.value?.absoluteChange)}`
      : '')

  const ytdChangePercentAbs = computed(() =>
    ytdChange.value?.changePercent !== undefined
      ? `${formatPercentAbs(Math.abs(ytdChange.value?.changePercent ?? 0))}`
      : '')

  const ytdPriceChange = computed(() =>
    ytdChange.value?.change !== undefined
      ? `${sign.value}${formatMoneyNoSymbol(ytdChange.value?.absoluteChange)}`
      : '')

  const ytdSign = computed(() => ytdChange.value?.sign ?? '')

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
    companyName,
    changePercentAbs,
    priceChange,
    ytdChangePercentAbs,
    ytdPriceChange,
    ytdSign,
    prevClose,
    price,
    snapshot,
    dailyHigh,
    dailyLow,
    dailyVolume,
    avgVolume,
  }
}
