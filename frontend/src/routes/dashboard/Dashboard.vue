<script
  lang="ts"
  setup
>
import DashboardSymbol from '@/routes/dashboard/DashboardSymbol.vue'
import { useRealtimeSymbolsStore } from '@/stores'
import { storeToRefs } from 'pinia'
import { ref, watch } from 'vue'
import { useAssetsStore } from '@/stores/useAssetsStore'
import Tags from '@/components/Tags.vue'

const liveSymbolsStore = useRealtimeSymbolsStore()

const { symbols } = storeToRefs(liveSymbolsStore)
const assetsStore = useAssetsStore()

const {
  assets,
  fetching,
  hasAssets,
} = storeToRefs(assetsStore)

const tags = ref<string[]>([])
const options = new Map<string, string>()

const loading = ref<boolean>(true)

watch(hasAssets, (newValue) => {
  if (newValue && options.size === 0) {
    for (const asset of assets.value) {
      options.set(asset.symbol, asset.name)
    }
    loading.value = false
  }
}, {
  immediate: true,
})
</script>

<template>
  <div
    class="container mx-auto max-w-[95vw] mt-5"
    v-if="!loading"
  >
    <div>
      <Tags
        :options="options"
        :tags="tags"
        @change="tags = $event"
        placeholder="Enter symbol and press Space or Enter"
      />
    </div>
    <main class="grid grid-cols-6 gap-1">
      <pre>{{ JSON.stringify(tags, null, 2)}}</pre>
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
