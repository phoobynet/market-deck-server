<script
  lang="ts"
  setup
>
import { useDeckStore } from '@/routes/deck/useDeckStore'
import { storeToRefs } from 'pinia'
import { onMounted, ref, watch } from 'vue'
import { useMagicKeys } from '@vueuse/core'
import DeckSearchModal from '@/routes/deck/DeckSearchModal.vue'

const deckStore = useDeckStore()

const {
  snapshots,
  symbols,
} = storeToRefs(deckStore)

const showModal = ref<boolean>(false)

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
  <main class="w-full px-2 mt-2" v-if="symbols">
    <table class="text-xs w-full">
      <thead>
        <tr class="border-b border-primary">
          <th class="text-left w-20">Symbol</th>
          <th class="text-left w-32">Company</th>
          <th class="text-right w-32">Price</th>
          <th class="text-right w-32">Prev. Close</th>
          <th class="text-right w-40">Change</th>
          <th class="text-right w-24">Low</th>
          <th class="text-right w-24">High</th>
          <th class="text-right w-24">Volume</th>
          <th class="text-right w-24">Avg. Volume</th>
          <th class="text-right w-24">YTD Change</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="symbol in symbols"
          :key="symbol"
          is="vue:deck-table-row"
          :symbol="symbol"
        />
      </tbody>
    </table>
    <Teleport to="body">
      <DeckSearchModal v-if="showModal" />
    </Teleport>
  </main>
</template>

<style
  lang="scss"
  scoped
>
</style>
