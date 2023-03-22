import { useCalendarDayUpdateStore, useSnapshotsStore } from '@/stores'
import { useAssetsStore } from '@/stores/useAssetsStore'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { Asset, Bar, Snapshot, SnapshotChange } from '@/types'
import { formatMoneyNoSymbol } from '@/libs/helpers/formatMoney'
import { formatPercentAbs } from '@/libs/helpers/formatPercent'

export const useReportCard = (symbol: string) => {
  const snapshotsStore = useSnapshotsStore()
  const assetsStore = useAssetsStore()
  const calendarDayUpdateStore = useCalendarDayUpdateStore()

  const { snapshots } = storeToRefs(snapshotsStore)
  const { previousDate } = storeToRefs(calendarDayUpdateStore)

  const snapshot = computed<Snapshot | undefined>(() => snapshots.value?.[symbol])

  const changeSincePrevious = computed<SnapshotChange | undefined>(() => {
    if (!previousDate.value ?? !snapshot.value) {
      return undefined
    }

    return snapshot.value?.changes[previousDate.value]
  })

  const asset = computed<Asset | undefined>(() => assetsStore.getBySymbol(symbol))

  const intradayBars = computed<Bar[]>(() => snapshot.value?.ibars ?? [])

  const intradayHigh = computed<number>(() => snapshot.value?.ih ?? 0)
  const intradayLow = computed<number>(() => snapshot.value?.il ?? 0)
  const currentPrice = computed<number>(() => snapshot.value?.lt?.p ?? 0)

  const priceChange = computed(() => formatMoneyNoSymbol(changeSincePrevious?.value?.ca ?? 0))
  const percentChange = computed(() => formatPercentAbs(changeSincePrevious?.value?.cp ?? 0))

  const signSymbol = computed<string>(() => {
    if (!changeSincePrevious.value?.cs) {
      return ''
    } else if (changeSincePrevious.value.cs === 1) {
      return '+'
    }

    return '-'
  })

  return {
    asset,
    snapshot,
    changeSincePrevious,
    intradayHigh,
    intradayLow,
    currentPrice,
    intradayBars,
    priceChange,
    percentChange,
    signSymbol,
  }
}
