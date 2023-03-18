<script
  lang="ts"
  setup
>
import DashboardSymbol from '@/routes/dashboard/DashboardSymbol.vue'
import { useRealtimeSymbolsStore } from '@/stores'
import { storeToRefs } from 'pinia'
import DashboardSearch from '@/routes/dashboard/DashboardSearch.vue'

const liveSymbolsStore = useRealtimeSymbolsStore()

const { symbols } = storeToRefs(liveSymbolsStore)
</script>

<template>
  <div class="container mx-auto max-w-[95vw] mt-5">
    <DashboardSearch />
    <main class="grid grid-cols-6 gap-1">
      <DashboardSymbol
        :symbol="symbol"
        v-for="symbol in symbols"
        :key="symbol"
      />
    </main>
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
