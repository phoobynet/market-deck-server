<script
  lang="ts"
  setup
>
import { computed } from 'vue'
import Money from '@/components/formatting/Money.vue'
import { useSnapshotsStore } from '@/stores/useSnapshotsStore'
import { storeToRefs } from 'pinia'
import { useAssetsStore } from '@/stores/useAssetsStore'
import { Asset, Trade } from '@/types'

const props = defineProps<{
  symbol: string
}>()

const snapshotsStore = useSnapshotsStore()
const assetsStore = useAssetsStore()

const { snapshots } = storeToRefs(snapshotsStore)

const snapshot = computed(() => {
  if (!snapshots.value) {
    return
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
</script>

<template>
  <div
    v-if="symbol && snapshot && asset"
    class="dashboard-report-card"
  >
    <div class="symbol">{{ asset?.S }}</div>
    <div class="name">{{ assetName }}</div>
    <pre>{{JSON.stringify(snapshot, null, 2)}}</pre>
    <div class="price">
      <Money
        :amount="latestTrade?.p"
        :show-sign="false"
        currency="$"
        :sexy="true"
      ></Money>
    </div>
    <div class="previous-close-change">
      <Money
        :amount="0"
        :show-sign="true"
        :sexy="true"
      />
    </div>
    <div class="previous-close">
      <Money
        :amount="0"
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
  .dashboard-report-card {
    @apply border px-2 py-1 rounded-md border-slate-600;

    .symbol {
      @apply text-xl tracking-widest;
    }

    .name {
      @apply font-light text-sm overflow-hidden;
    }
  }
</style>
