<script
  lang="ts"
  setup
>
import { useSnapshot } from '@/routes/deck/useSnapshot'

const props = defineProps<{
  symbol: string
}>()

const {
  snapshot,
  price,
  priceChange,
  changePercent,
  sign,
  prevClose,
  dailyHigh,
  dailyLow,
  dailyVolume,
} = useSnapshot(props.symbol)

</script>

<template>
  <div v-if="snapshot">
    <div>{{ snapshot?.symbol }}</div>
    <div>{{ snapshot?.name }}</div>
    <div>{{ snapshot?.exchange }}</div>
    <div>Price: {{ price }}</div>
    <div>Previous Close: {{ prevClose }}</div>
    <div>Change: {{ priceChange }}</div>
    <div>Change %: {{ changePercent }}</div>
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
