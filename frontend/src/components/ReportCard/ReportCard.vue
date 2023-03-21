<script
  lang="ts"
  setup
>
import { computed, provide } from 'vue'
import { useSnapshotsStore } from '@/stores/useSnapshotsStore'
import { storeToRefs } from 'pinia'
import { useAssetsStore } from '@/stores/useAssetsStore'
import { Asset, Snapshot, SnapshotChange } from '@/types'
import { AssetKey, ChangeSincePreviousKey, SnapshotKey } from '@/components/ReportCard/injectionKeys'
import ReportCardSymbol from '@/components/ReportCard/ReportCardSymbol.vue'
import ReportCardName from '@/components/ReportCard/ReportCardName.vue'
import ReportCardLatestPrice from '@/components/ReportCard/ReportCardLatestPrice.vue'
import ReportCardDailyVolume from '@/components/ReportCard/ReportCardDailyVolume.vue'
import ReportCardExchange from '@/components/ReportCard/ReportCardExchange.vue'
import { useCalendarDayUpdateStore } from '@/stores'

const props = defineProps<{
  symbol: string
}>()

const snapshotsStore = useSnapshotsStore()
const assetsStore = useAssetsStore()
const calendarDayUpdateStore = useCalendarDayUpdateStore()

const { snapshots } = storeToRefs(snapshotsStore)
const { previousDate } = storeToRefs(calendarDayUpdateStore)

const snapshot = computed<Snapshot | undefined>(() => {
  if (!snapshots.value) {
    return undefined
  }

  return snapshots.value[props.symbol]
})

const changeSincePrevious = computed<SnapshotChange | undefined>(() => {
  if (!previousDate.value ?? !snapshot.value) {
    return undefined
  }

  return snapshot.value?.changes[previousDate.value]
})

const asset = computed<Asset | undefined>(() => {
  return assetsStore.getBySymbol(props.symbol)
})

provide(AssetKey, asset)
provide(SnapshotKey, snapshot)
provide(ChangeSincePreviousKey, changeSincePrevious)
</script>

<template>
  <div
    v-if="symbol && snapshot && asset"
    class="report-card"
    :data-sign="changeSincePrevious?.cs"
  >
    <ReportCardSymbol />
    <ReportCardName class="name" />
    <ReportCardLatestPrice class="latest-price" />
    <ReportCardDailyVolume class="daily-volume" />
    <ReportCardExchange class="exchange" />
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .report-card {
    @apply border px-2 py-1 rounded-md border-slate-600 transition-all duration-300 ease-in-out;

    &[data-sign="1"] {
      @apply border-up;
    }

    &[data-sign="-1"] {
      @apply border-down;
    }

    display: grid;
    grid-template-areas:
      "name name exchange"
      "symbol latest-price latest-price"
      "previous-close previous-close previous-close"
      ". .daily-volume"
      "changes changes changes";

    .symbol {
      grid-area: symbol;
      @apply text-xl tracking-widest;
    }

    .name {
      grid-area: name;
      @apply self-start;
    }

    .latest-price {
      grid-area: latest-price;
      @apply justify-self-end;
    }

    .previous-close {
      grid-area: previous-close;
    }

    .changes {
      grid-area: changes;
      @apply self-start;
    }

    .daily-volume {
      grid-area: daily-volume;
      @apply justify-self-end;
    }

    .exchange {
      grid-area: exchange;
      @apply justify-self-end;
    }
  }
</style>
