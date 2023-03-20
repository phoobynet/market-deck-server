<script
  lang="ts"
  setup
>
import { useSnapshotsStore } from '@/stores'
import { storeToRefs } from 'pinia'
import { onMounted, ref, watch } from 'vue'
import { useAssetsStore } from '@/stores/useAssetsStore'
import Tags from '@/components/Tags.vue'
import { debouncedWatch } from '@vueuse/core'
import { updateSymbols } from '@/libs/updateSymbols'
import { getSymbols } from '@/libs/getSymbols'
import ReportCard from '@/components/ReportCard/ReportCard.vue'

const snapshotsStore = useSnapshotsStore()

const { symbols, snapshots } = storeToRefs(snapshotsStore)
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
      options.set(asset.S, asset.n)
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
    <main class="grid lg:grid-cols-5 md:grid-cols-3 sm:grid-cols-2 xs:grid-cols-1 gap-1 mt-3">
      <transition-group
        enter-active-class="animate__animated animate__fadeIn animate__faster"
        leave-active-class="animate__animated animate__fadeOut animate__faster"
      >
        <ReportCard
          v-for="symbol of symbols"
          :key="symbol"
          :symbol="symbol"
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
