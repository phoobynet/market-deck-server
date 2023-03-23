<script
  lang="ts"
  setup
>
import { useDeckStore } from '@/routes/deck/useDeckStore'
import { useAssetsStore } from '@/stores/useAssetsStore'
import { storeToRefs } from 'pinia'
import { onMounted, ref } from 'vue'
import { debouncedWatch } from '@vueuse/core'

const assetsStore = useAssetsStore()
const deckStore = useDeckStore()

const {
  assets,
  hasAssets,
} = storeToRefs(assetsStore)

const query = ref<string>('')
const symbols = ref<string[]>([])
const searchInput = ref<HTMLInputElement>()

// debouncedWatch(symbols, async (newValue) => {
//   await deckStore.updateSymbols(newValue)
// }, {
//   immediate: true,
//   debounce: 500,
// })

onMounted(async () => {
  symbols.value = await deckStore.getSymbols()
})

</script>

<template>
  <div>
    <div>
      <input
        ref="searchInput"
        type="text"
        @change="query = $event"
      >
    </div>
    <div class="flex flex-row flex-wrap gap-1">
      <div
        v-for="symbol in symbols"
        :key="symbols"
      >
        {{ symbol }}
      </div>
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
