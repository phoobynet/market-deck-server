<script
  lang="ts"
  setup
>
import { useDeckSnapshot } from '@/routes/deck/useDeckSnapshot'

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
} = useDeckSnapshot(props.symbol)

</script>

<template>
  <div v-if="snapshot">
    <div>{{ snapshot?.symbol }}</div>
    <div>{{ snapshot?.name }}</div>
    <div>{{ snapshot?.exchange }}</div>
    <div>Price: {{ price }}</div>
    <div>Previous Close: {{ prevClose }}</div>
    <div>Change: {{ priceChange }}</div>
    <div>Change %: {{ changePercentAbs }}</div>
    <div>Daily High: {{ dailyHigh }}</div>
    <div>Daily Low: {{ dailyLow }}</div>
    <div>Daily Volume: {{ dailyVolume }}</div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
