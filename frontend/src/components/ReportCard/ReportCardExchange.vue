<script
  lang="ts"
  setup
>
import { AssetKey } from '@/components/ReportCard/injectionKeys'
import { inject } from 'vue'

const asset = inject(AssetKey)
</script>

<template>
  <div class="exchange">{{ asset?.x}}</div>
</template>

<style
  lang="scss"
  scoped
>
  .exchange {
    @apply text-xs opacity-60;
  }
</style>
