<script
  lang="ts"
  setup
>
import { Asset } from '@/types'
import { ref } from 'vue'
import { debouncedWatch } from '@vueuse/core'

const query = ref<string>('')
const assets = ref<Asset[]>([])
const searching = ref<boolean>(false)

const onSearch = async (query: string): Promise<void> => {
  searching.value = true

  try {
    assets.value = await SearchAssets(query)
  } catch (e) {
    console.log(e)
  } finally {
    searching.value = false
  }
}

debouncedWatch(query, async (newValue) => {
  await onSearch(newValue)
}, {
  debounce: 500,
})

</script>

<template>
  <div>
    <input
      type="text"
      v-model="query"
      class="input input-sm w-full input-bordered"
    >
    <div>
      <ul>
        <li
          v-for="asset in assets"
          :key="asset.symbol"
        >
          {{ asset.symbol }} - {{ asset.name }}
        </li>
      </ul>
    </div>
  </div>
</template>
