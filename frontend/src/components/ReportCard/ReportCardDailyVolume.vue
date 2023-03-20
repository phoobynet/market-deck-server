<script
  lang="ts"
  setup
>
import { SnapshotKey } from '@/components/ReportCard/injectionKeys'
import { computed, inject } from 'vue'
import { formatLargeNumber } from '@/libs/helpers/formatLargeNumber'

const snapshot = inject(SnapshotKey)

const volume = computed(() => {
  return formatLargeNumber(snapshot?.value?.db?.v ?? 0)
})

</script>

<template>
  <div class="text-xs tabular-nums">
    <span class="opacity-60">Volume:</span> {{ volume }}
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
