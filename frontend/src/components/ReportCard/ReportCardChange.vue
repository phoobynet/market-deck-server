<script
  lang="ts"
  setup
>
import { SignKey, SignSymbolKey, SnapshotKey } from '@/components/ReportCard/injectionKeys'
import { computed, inject } from 'vue'
import { formatMoneyNoSymbol } from '@/libs/helpers/formatMoney'
import { formatPercent } from '@/libs/helpers/formatPercent'
import ReportCardChangeIndicator from '@/components/ReportCard/ReportCardChangeIndicator.vue'

const snapshot = inject(SnapshotKey)
const sign = inject(SignKey)
const signSymbol = inject(SignSymbolKey)

const priceChange = computed(() => {
  return formatMoneyNoSymbol(snapshot?.value?.ca ?? 0)
})

const percentChange = computed(() => {
  return formatPercent(snapshot?.value?.cp ?? 0)
})
</script>

<template>
  <div
    class="change"
    :data-sign="sign"
  >
    <div>
      <ReportCardChangeIndicator :sign="sign" />
    </div>
    <div class="price-change">{{signSymbol}} {{ priceChange }}</div>
    <div class="percent-change">({{ percentChange }})</div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .change {
    @apply flex items-center justify-end gap-2 tabular-nums;

    &[data-sign="1"] {
      @apply text-up;
    }

    &[data-sign="-1"] {
      @apply text-down;
    }
  }

</style>
