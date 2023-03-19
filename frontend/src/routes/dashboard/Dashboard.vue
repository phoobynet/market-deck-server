<script
  lang="ts"
  setup
>
import DashboardSymbol from '@/routes/dashboard/DashboardSymbol.vue'
import { useSnapshots } from '@/stores'
import { storeToRefs } from 'pinia'
import { onMounted, ref, watch } from 'vue'
import { useAssetsStore } from '@/stores/useAssetsStore'
import Tags from '@/components/Tags.vue'
import { debouncedWatch } from '@vueuse/core'
import { updateSymbols } from '@/libs/updateSymbols'
import { getSymbols } from '@/libs/getSymbols'

const snapshotsStore = useSnapshots()

const { symbols } = storeToRefs(snapshotsStore)
const assetsStore = useAssetsStore()

const {
  assets,
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

debouncedWatch(tags, async (newValue) => {
  await updateSymbols(newValue)
}, {
  immediate: true,
  debounce: 500,
})

onMounted(async () => {
  tags.value = await getSymbols()
})
</script>

<template>
  <div
    class="mx-auto mx-2 mt-2 md:mx-4"
    v-if="!loading"
  >
    <Tags
      :options="options"
      :tags="tags"
      @change="tags = $event"
      placeholder="Enter symbol and press Space or Enter"
    />
    <main class="grid grid-cols-6 gap-1 mt-3">
      <transition-group
        enter-active-class="animate__animated animate__fadeIn animate__faster"
        leave-active-class="animate__animated animate__fadeOut animate__faster"
      >
        <DashboardSymbol
          :symbol="symbol"
          v-for="symbol in symbols"
          :key="symbol"
        />
      </transition-group>
    </main>
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
