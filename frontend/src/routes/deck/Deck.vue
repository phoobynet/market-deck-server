<script
  lang="ts"
  setup
>
import { useDeckStore } from '@/routes/deck/useDeckStore'
import { storeToRefs } from 'pinia'
import { onMounted, watch } from 'vue'
import { useMagicKeys } from '@vueuse/core'
import DeckSearchModal from '@/routes/deck/DeckSearchModal.vue'
import DeckCard from '@/routes/deck/DeckCard.vue'

const deckStore = useDeckStore()

const {
  snapshots,
  symbols,
  showModal,
} = storeToRefs(deckStore)

const keys = useMagicKeys()

const cmdK = keys['Cmd+K']
const esc = keys['Escape']

watch(cmdK, (value) => {
  if (value) {
    showModal.value = true
  }
})

watch(esc, (value) => {
  if (value) {
    showModal.value = false
  }
})

onMounted(async () => {
  await deckStore.getSymbols()

})
</script>

<template>
  <main
    class="w-full px-2 mt-2"
    v-if="symbols"
  >
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 w-full gap-1 pt-1">
      <DeckCard
        v-for="symbol in symbols"
        :key="symbol"
        :symbol="symbol"
      />
    </div>

    <Teleport to="body">
      <DeckSearchModal />
    </Teleport>
  </main>
</template>

<style
  lang="scss"
  scoped
>
</style>
