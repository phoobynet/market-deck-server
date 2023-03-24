<script
  lang="ts"
  setup
>
import { debouncedWatch, onClickOutside } from '@vueuse/core'
import { ref, watch } from 'vue'
import Tags from '@/components/Tags.vue'
import { storeToRefs } from 'pinia'
import { useAssetsStore } from '@/stores/useAssetsStore'
import { useDeckStore } from '@/routes/deck/useDeckStore'

const assetsStore = useAssetsStore()
const deckStore = useDeckStore()

const {
  assets,
  hasAssets,
} = storeToRefs(assetsStore)

const {
  showModal,
  symbols,
} = storeToRefs(deckStore)

const tags = ref<string[]>([])
const options = new Map<string, string>()

const loading = ref<boolean>(true)
const modal = ref<HTMLDivElement>()

onClickOutside(modal, () => {
  showModal.value = false
})

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
  await deckStore.updateSymbols(newValue)
}, {
  immediate: true,
  debounce: 500,
})

watch(symbols, (newValue) => {
  tags.value = newValue
}, {
  immediate: true,
})
</script>

<template>
  <div
    class="deck-search-modal"
    v-if="showModal"
  >
    <div
      class="bg-slate-900 h-28 min-w-[60vw] border rounded-lg border-slate-700 px-4 py-2 flex flex-col gap-2"
      ref="modal"
    >

      <h2 class="pl-1 text-1xl">
        Asset Search
      </h2>
      <Tags
        :options="options"
        :tags="tags"
        @change="tags = $event"
        placeholder="Search..."
      />
      <div class="text-slate-600 text-xs"><i><strong>Escape</strong></i> to close</div>
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .deck-search-modal {
    @apply fixed top-0 left-0 w-screen h-screen z-50 flex justify-center pt-32;
    background-color: rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(5px);
  }
</style>
