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
import { formatMoney } from '@/libs/helpers/formatMoney'

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

const formattedPrice = computed(() => {
  return formatMoney(snapshot.value?.lt?.p)
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
    <header class="h-6 flex items-center justify-between pl-1">
      <div class="text-xs">
        {{ assetName }}
      </div>
      <Icon
          :size="25"
          class="cursor-pointer"
          @click="close"
      >
        <Close/>
      </Icon>
    </header>
    <div class="content">
      <div class="font-bold tracking-wider">
        {{ symbol }}
      </div>
      <div>{{ formattedPrice }}</div>
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

  .content {
    @apply p-1;
    grid-template-columns: auto auto;
    grid-template-rows: auto auto;
  }
}
</style>
