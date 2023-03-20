<script
  lang="ts"
  setup
>
import { computed } from 'vue'
import Money from '@/components/formatting/Money.vue'
import { useSnapshotsStore } from '@/stores/useSnapshotsStore'
import { storeToRefs } from 'pinia'
import { useAssetsStore } from '@/stores/useAssetsStore'
import { Asset, Snapshot, Trade } from '@/types'
import Percentage from '@/components/formatting/Percentage.vue'

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

const nameCleaner = (name: string | undefined): string => {
  name = (name ?? '').trim()

  return name
    .replace('Common Stock', '')
    .replace('American Depositary Shares', '(ADS)')
    .replace('American Depositary Receipts', '(ADR)')
}

const asset = computed<Asset | undefined>(() => {
  return assetsStore.getBySymbol(props.symbol)
})

const assetName = computed<string>(() => {
  return nameCleaner(asset.value?.n)
})

const latestTrade = computed<Trade | undefined>(() => {
  return snapshot.value?.lt
})

const previousClose = computed<number>(() => {
  return snapshot.value?.pc ?? 0
})

const previousClosePercentChange = computed<number>(() => {
  return snapshot.value?.cp ?? 0
})

const change = computed<number>(() => {
  return snapshot.value?.c ?? 0
})
</script>

<template>
  <div
    v-if="symbol && snapshot && asset"
    class="report-card"
  >
    <div class="symbol">{{ symbol }}</div>
    <div class="name">{{ assetName }}</div>
    <div class="latest-price">
      <Money
        :amount="latestTrade?.p"
        :show-sign="false"
        currency="$"
        :sexy="true"
      ></Money>
    </div>
    <div class="previous-close-change">
      <Money
        :amount="change"
        :show-sign="true"
        :sexy="true"
      />
    </div>
    <div class="previous-close-change-percentage">
      <Percentage :amount="previousClosePercentChange"></Percentage>
    </div>
    <div class="previous-close">
      <Money
        :amount="previousClose"
        :show-sign="false"
        currency="$"
      ></Money>
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .report-card {
    @apply border px-2 py-1 rounded-md border-slate-600 grid;

    grid-template-areas:
      "symbol price price"
      "name name name"
      "previous-close previous-close-change previous-close-change-percentage";

    grid-template-columns: repeat(3, 1fr);
    grid-template-rows: repeat(3, 1fr);

    .symbol {
      grid-area: symbol;
      @apply text-xl tracking-widest;
    }

    .name {
      grid-area: name;
      @apply font-light text-sm overflow-hidden;
    }

    .latest-price {
      grid-area: price;
      @apply justify-self-end text-xl;
    }

    .previous-close {
      grid-area: previous-close;
    }

    .previous-close-change {
      grid-area: previous-close-change;
    }

    .previous-close-change-percentage {
      grid-area: previous-close-change-percentage;
    }
  }
</style>
