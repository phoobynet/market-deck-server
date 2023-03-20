<script
  lang="ts"
  setup
>
import { computed, provide } from 'vue'
import { useSnapshotsStore } from '@/stores/useSnapshotsStore'
import { storeToRefs } from 'pinia'
import { useAssetsStore } from '@/stores/useAssetsStore'
import { Asset, Snapshot } from '@/types'
import { AssetKey, SignKey, SignSymbolKey, SnapshotKey } from '@/components/ReportCard/injectionKeys'
import ReportCardSymbol from '@/components/ReportCard/ReportCardSymbol.vue'
import ReportCardName from '@/components/ReportCard/ReportCardName.vue'
import ReportCardLatestPrice from '@/components/ReportCard/ReportCardLatestPrice.vue'
import ReportCardDailyVolume from '@/components/ReportCard/ReportCardDailyVolume.vue'
import ReportCardExchange from '@/components/ReportCard/ReportCardExchange.vue'
import ReportCardChange from '@/components/ReportCard/ReportCardChange.vue'

const props = defineProps<{
  symbol: string
}>()

const snapshotsStore = useSnapshotsStore()
const assetsStore = useAssetsStore()

const { snapshots } = storeToRefs(snapshotsStore)

const snapshot = computed<Snapshot | undefined>(() => {
  if (!snapshots.value) {
    return undefined
  }

  return snapshots.value[props.symbol]
})

const asset = computed<Asset | undefined>(() => {
  return assetsStore.getBySymbol(props.symbol)
})

const sign = computed(() => {
  return snapshot?.value?.cs ?? 0
})

const signSymbol = computed(() => {
  if (!sign.value) {
    return ''
  }

  return sign.value > 0
    ? '+'
    : '-'
})

provide(AssetKey, asset)
provide(SnapshotKey, snapshot)
provide(SignKey, sign)
provide(SignSymbolKey, signSymbol)
</script>

<template>
  <div
    v-if="symbol && snapshot && asset"
    class="report-card"
    :data-sign="sign"
  >
    <ReportCardSymbol />
    <ReportCardName class="name" />
    <ReportCardLatestPrice class="latest-price" />
    <ReportCardChange class="change" />
    <div class="previous-close"></div>
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
    //grid-template-columns: repeat(3, 1fr);
    //grid-template-rows: auto auto auto repeat(2, 1fr);
    grid-template-areas:
      "name name name"
      "symbol latest-price latest-price"
      "exchange change change"
      "previous-close previous-close previous-close"
      "daily-volume . .";

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

    .change {
      grid-area: change;
      @apply self-start;
    }

    .daily-volume {
      grid-area: daily-volume;
    }

    .exchange {
      grid-area: exchange;
    }
  }
</style>
