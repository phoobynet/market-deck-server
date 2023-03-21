<script
  lang="ts"
  setup
>
import { ChangeSincePreviousKey, SnapshotKey } from '@/components/ReportCard/injectionKeys'
import { computed, inject } from 'vue'
import { formatMoneyNoSymbol } from '@/libs/helpers/formatMoney'
import { formatPercentAbs } from '@/libs/helpers/formatPercent'
import ReportCardChangeIndicator from '@/components/ReportCard/ReportCardChangeIndicator.vue'

const snapshot = inject(SnapshotKey)
const changeSincePrevious = inject(ChangeSincePreviousKey)

const priceChange = computed(() => {
  return formatMoneyNoSymbol(changeSincePrevious?.value?.ca ?? 0)
})

const percentChange = computed(() => {
  return formatPercentAbs(changeSincePrevious?.value?.cp ?? 0)
})

const signSymbol = computed(() => {
  if (changeSincePrevious?.value?.cs == 1) {
    return '+'
  } else if (changeSincePrevious?.value?.cs == -1) {
    return '-'
  } else {
    return ''
  }
})
</script>

<template>
  <div
    class="change"
    :data-sign="changeSincePrevious?.cs"
  >
    <div>
      <ReportCardChangeIndicator :sign="changeSincePrevious?.cs" />
    </div>
    <div class="price-change"><span class="tracking-widest">{{ signSymbol }}</span>{{ priceChange }}</div>
    <div class="percent-change">({{ percentChange }})</div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .change {
    @apply flex items-center justify-end gap-1 tabular-nums;

    &[data-sign="1"] {
      @apply text-up;
    }

    &[data-sign="-1"] {
      @apply text-down;
    }
  }

</style>
