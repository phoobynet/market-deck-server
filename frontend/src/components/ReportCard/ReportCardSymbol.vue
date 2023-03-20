<script
  lang="ts"
  setup
>
import { computed, inject } from 'vue'
import { AssetKey } from '@/components/ReportCard/injectionKeys'

const asset = inject(AssetKey)
const symbol = computed(() => asset?.value?.S || '')
</script>

<template>
  <div
    class="text-2xl font-semibold tracking-widest"
  >{{ symbol }}
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
