<script
  lang="ts"
  setup
>
import { computed, inject } from 'vue'
import { AssetKey } from '@/components/ReportCard/injectionKeys'
import { cleanAssetName } from '@/libs/helpers/cleanAssetName'

const asset = inject(AssetKey)
const name = computed(() => cleanAssetName(asset?.value?.n))
</script>

<template>
  <div
    class="text-xs opacity-60 truncate"
  >{{ name }}
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
