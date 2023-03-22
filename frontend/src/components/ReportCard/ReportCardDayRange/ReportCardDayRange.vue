<script
  lang="ts"
  setup
>
import { computed, inject, onBeforeUnmount, onMounted, Ref, ref, watch, watchEffect } from 'vue'
import { select, Selection } from 'd3-selection'
import {
  CurrentPriceKey, IntradayHighKey, IntradayLowKey, SignSymbolKey, SnapshotKey, SymbolKey, WidthKey,
} from '@/components/ReportCard/injectionKeys'
import { scaleLinear } from 'd3-scale'
import 'd3-transition'
import { RangeBar } from '@/components/ReportCard/ReportCardDayRange/RangeBar'
import { RangeBarPrice } from '@/components/ReportCard/ReportCardDayRange/RangeBarPrice'
import { CurrentPriceIndicator } from '@/components/ReportCard/ReportCardDayRange/CurrentPriceIndicator'

const dayRange = ref<SVGElement>()

const snapshot = inject(SnapshotKey)
const width = inject(WidthKey) as Ref<number>
const intradayHigh = inject(IntradayHighKey) as Ref<number>
const intradayLow = inject(IntradayLowKey) as Ref<number>
const currentPrice = inject(CurrentPriceKey) as Ref<number>
const symbol = inject(SymbolKey) as string
const signSymbol = inject(SignSymbolKey) as Ref<string>

let svg: Selection<SVGElement, null, null, undefined>
let rangeBar: RangeBar
let rangeBarLow: RangeBarPrice
let rangeBarHigh: RangeBarPrice
let currentPriceIndicator: CurrentPriceIndicator

const scale = computed(() =>
  scaleLinear()
    .domain([intradayLow.value, intradayHigh.value])
    .range([0, (width?.value ?? 0)]))

function updateCurrentPriceIndicator (price: number) {
  if (currentPriceIndicator) {
    currentPriceIndicator.priceScaled = (scale?.value ?? (() => 0))(price)
  }
}

function initialRender () {
  svg = select(dayRange.value!)

  rangeBar = new RangeBar(svg)
  rangeBar.x = 2
  rangeBar.width = (width?.value ?? 0) - 4
  rangeBar.signSymbol = signSymbol.value || ''

  rangeBarLow = new RangeBarPrice(svg)
  rangeBarLow.x = 4
  rangeBarLow.y = 14
  rangeBarLow.price = intradayLow.value

  rangeBarHigh = new RangeBarPrice(svg)
  rangeBarHigh.y = 14
  rangeBarHigh.x = (width.value - rangeBarHigh.computedWidth) - 5
  rangeBarHigh.price = intradayHigh.value

  currentPriceIndicator = new CurrentPriceIndicator(svg)
  currentPriceIndicator.priceScaled = scale.value?.(currentPrice.value) ?? 0
}

watchEffect(() => {
  if (rangeBar && rangeBarHigh) {
    rangeBar.width = (width?.value ?? 0) - 4
  }

  if (rangeBarLow) {
    rangeBarLow.price = intradayLow.value
  }

  if (rangeBarHigh) {
    rangeBarHigh.price = intradayHigh.value
    rangeBarHigh.x = (width.value - rangeBarHigh.computedWidth) - 5
  }

  updateCurrentPriceIndicator(currentPrice.value)
})

onMounted(() => {
  initialRender()
})

onBeforeUnmount(() => {
  svg?.remove()
})
</script>

<template>
  <svg
    ref="dayRange"
    :id="symbol"
    height="60"
    :width="width"
  ></svg>
</template>


<style scoped>

</style>
