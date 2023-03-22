<script
  lang="ts"
  setup
>
import { useSnapshotsStore } from '@/stores'
import { useAssetsStore } from '@/stores/useAssetsStore'
import { storeToRefs } from 'pinia'
import { computed, onMounted, ref } from 'vue'
import { Asset } from '@/types'
import { cleanAssetName } from '@/libs/helpers/cleanAssetName'
import { Icon } from '@vicons/utils'
import { Close } from '@vicons/carbon'

const props = defineProps<{
  symbol: string,
}>()

const snapshotStore = useSnapshotsStore()
const assetStore = useAssetsStore()

const { snapshots } = storeToRefs(snapshotStore)

const asset = ref<Asset>()

const assetName = computed(() => {
  return cleanAssetName(asset.value?.n)
})

const snapshot = computed(() => {
  return snapshots.value?.[props.symbol]
})

const emit = defineEmits(['close'])

const close = () => {
  emit('close', props.symbol)
}

onMounted(() => {
  asset.value = assetStore.getBySymbol(props.symbol)
})

</script>

<template>
  <div class="card">
    <header class="h-8 flex items-center justify-between px-1">
      <div>
        {{ assetName }}
      </div>
      <Icon
        :size="25"
        class="cursor-pointer"
        @click="close"
      >
        <Close />
      </Icon>
    </header>
    <div class="font-bold tracking-wider">
      {{ symbol }}
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .card {
    @apply border rounded overflow-hidden;

    display: grid;

    @media (min-width: 768px) {
      grid-template-columns: auto auto;
    }
  }
</style>
