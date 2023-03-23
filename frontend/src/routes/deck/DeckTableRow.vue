<script
  lang="ts"
  setup
>
import { useSnapshot } from '@/routes/deck/useSnapshot'
import { cleanAssetName } from '@/libs/helpers/cleanAssetName'

const props = defineProps<{
  symbol: string
}>()

const {
  snapshot,
  price,
  priceChange,
  changePercentAbs,
  sign,
  prevClose,
  dailyHigh,
  dailyLow,
  dailyVolume,
  avgVolume,
  ytdChangePercentAbs,
  ytdPriceChange,
} = useSnapshot(props.symbol)

</script>

<template>
  <tr
    class="border-b border-slate-700 h-8"
    :data-sign="sign"
  >
    <td
      class="text-left tracking-wider text-2xl"
      data-name="symbol"
    >{{ snapshot?.symbol }}
    </td>
    <td
      class="text-left"
      data-name="exchange"
    >
      <div class="flex flex-col">
        <div>
          {{ cleanAssetName(snapshot?.name) }}
        </div>
        <div class="opacity-60 font-light">
          {{ snapshot?.exchange }}
        </div>
      </div>
    </td>
    <td
      class="text-right tabular-nums font-bold text-xl"
      data-name="price"
    >{{ price }}
    </td>
    <td
      class="text-right tabular-nums"
      data-name="prevClose"
    >{{ prevClose }}
    </td>
    <td
      class="text-right tabular-nums"
      data-name="change"
    >
      <div class="flex flex-row w-full">
        <div class="basis-1/2 font-semibold">{{ priceChange }}</div>
        <div class="basis-1/2">({{ changePercentAbs }})</div>
      </div>
    </td>
    <td
      class="text-right tabular-nums"
      data-name="dailyLow"
    >{{ dailyLow }}
    </td>
    <td
      class="text-right tabular-nums"
      data-name="dailyHigh"
    >{{ dailyHigh }}
    </td>
    <td
      class="text-right tabular-nums"
      data-name="dailyVolume"
    >{{ dailyVolume }}
    </td>
    <td
      class="text-right tabular-nums"
      data-name="avgVolume"
    >{{ avgVolume }}
    </td>
    <td
      class="text-right tabular-nums"
      data-name="change"
    >
      <div class="flex flex-row w-full">
        <div class="basis-1/2 font-semibold">{{ ytdPriceChange }}</div>
        <div class="basis-1/2">({{ ytdChangePercentAbs }})</div>
      </div>
    </td>
  </tr>
</template>

<style
  lang="scss"
  scoped
>
  [data-sign="+"] {
    [data-name="change"], [data-name="symbol"], [data-name="price"] {
      @apply text-up;
    }
  }

  [data-sign="-"] {
    [data-name="change"], [data-name="symbol"], [data-name="price"] {
      @apply text-down;
    }
  }
</style>
