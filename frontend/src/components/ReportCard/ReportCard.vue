<script
  lang="ts"
  setup
>
import { computed, provide, ref, watch } from 'vue'
import { useSnapshotsStore } from '@/stores/useSnapshotsStore'
import { storeToRefs } from 'pinia'
import { useAssetsStore } from '@/stores/useAssetsStore'
import { Asset, Snapshot, SnapshotChange } from '@/types'
import {
  AssetKey, ChangeSincePreviousKey, HeightKey, SnapshotKey, WidthKey,
} from '@/components/ReportCard/injectionKeys'
import ReportCardSymbol from '@/components/ReportCard/ReportCardSymbol.vue'
import ReportCardName from '@/components/ReportCard/ReportCardName.vue'
import ReportCardLatestPrice from '@/components/ReportCard/ReportCardLatestPrice.vue'
import ReportCardDailyVolume from '@/components/ReportCard/ReportCardDailyVolume.vue'
import ReportCardExchange from '@/components/ReportCard/ReportCardExchange.vue'
import { useCalendarDayUpdateStore } from '@/stores'
import ReportCardDayRange from '@/components/ReportCard/ReportCardDayRange.vue'
import { useElementSize } from '@vueuse/core'

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

const reportCard = ref<HTMLDivElement>()

const {
  height,
  width,
} = useElementSize(reportCard)

watch(width,  () => {
  console.log('width', width.value)
})

provide(AssetKey, asset)
provide(SnapshotKey, snapshot)
provide(ChangeSincePreviousKey, changeSincePrevious)
provide(WidthKey, width)
provide(HeightKey, height)
</script>

<template>
  <div
    ref="reportCard"
    v-if="symbol && snapshot && asset"
    class="report-card"
    :data-sign="changeSincePrevious?.cs"
  >
    <ReportCardSymbol />
    <ReportCardName class="name" />
    <ReportCardLatestPrice class="latest-price" />
    <ReportCardDailyVolume class="daily-volume" />
    <ReportCardExchange class="exchange" />
    <ReportCardDayRange class="day-range" />
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
      ". . daily-volume"
      "day-range day-range day-range";

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

    .day-range {
      grid-area: day-range;
    }
  }
</style>
