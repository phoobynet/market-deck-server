<script
  lang="ts"
  setup
>
import { computed, inject } from 'vue'
import { ChangeSincePreviousKey, SnapshotKey } from '@/components/ReportCard/injectionKeys'
import { formatMoney } from '@/libs/helpers/formatMoney'

const snapshot = inject(SnapshotKey)
const changeSincePrevious = inject(ChangeSincePreviousKey)

const formatted = computed(() => {
  return formatMoney(snapshot?.value?.lt?.p ?? 0)
})
</script>

<template>
  <div
    class="latest-price"
    :data-sign="changeSincePrevious?.cs"
  >
    {{ formatted }}
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .latest-price {
    @apply tabular-nums text-2xl font-semibold;

    &[data-sign="1"] {
      @apply text-up;
    }

    &[data-sign="-1"] {
      @apply text-down;
    }
  }
</style>
