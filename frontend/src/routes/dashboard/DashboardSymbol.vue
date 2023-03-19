<script
  lang="ts"
  setup
>
import { computed } from 'vue'
import Money from '@/components/formatting/Money.vue'
import { useSnapshots } from '@/stores/useSnapshots'
import { storeToRefs } from 'pinia'

const props = defineProps<{
  symbol: string
}>()

const snapshotsStore = useSnapshots()

const { snapshots } = storeToRefs(snapshotsStore)

const liveSymbol = computed(() => {
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

const assetName = computed<string>(() => {
  return nameCleaner(liveSymbol.value?.asset?.name)
})
</script>

<template>
  <div
    v-if="symbol && liveSymbol"
    class="dashboard-symbol"
  >
    <div class="symbol">{{ liveSymbol?.trade?.S }}</div>
    <div class="name">{{ assetName }}</div>
    <div class="price">
      <Money
        :amount="liveSymbol?.trade?.p"
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
  .dashboard-symbol {
    @apply border px-2 py-1 rounded-md border-slate-600;

    .symbol {
      @apply text-xl tracking-widest;
    }

    .name {
      @apply font-light text-sm overflow-hidden;
    }
  }
</style>
