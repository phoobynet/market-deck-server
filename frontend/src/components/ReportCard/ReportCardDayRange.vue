<script
  lang="ts"
  setup
>
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { select, Selection } from 'd3-selection'
import { SnapshotKey, WidthKey } from '@/components/ReportCard/injectionKeys'
import { scaleLinear } from 'd3-scale'
import { symbol, symbolTriangle } from 'd3-shape'
import { formatMoneyNoSymbol } from '@/libs/helpers/formatMoney'


const dayRange = ref<HTMLDivElement>()

const snapshot = inject(SnapshotKey)
const width = inject(WidthKey)

const low = computed(() => snapshot?.value?.db?.l ?? 0)
const high = computed(() => snapshot?.value?.db?.h ?? 0)
const currentPrice = computed(() => snapshot?.value?.lt?.p ?? 0)

const scale = computed(() => {
  if (low.value && high.value) {
    return scaleLinear().domain([low.value, high.value]).range([0, width?.value ?? 0])
  }

  return undefined
})

let dayRangeBar: Selection<HTMLDivElement, null, null, undefined>

watch(() => width?.value, () => {
  if (dayRangeBar) {
    dayRangeBar
      .select('rect')
      .attr('width', width?.value ?? 0)
  }
})

watch(currentPrice, (newValue) => {
  if (dayRangeBar) {
    dayRangeBar
      .select('#currentPrice')
      .transition()
      .duration(500)
      .attr('transform', `translate(${scale.value?.(newValue) ?? 0}, 35)`)
  }
})

watch(low, (newValue) => {
  if (dayRangeBar) {
    dayRangeBar
      .select('#low')
      .text(newValue)
  }
})

watch(high, (newValue) => {
  if (dayRangeBar) {
    dayRangeBar
      .select('#high')
      .text(newValue)
  }
})

onMounted(() => {
  nextTick(() => {
    if (dayRange.value) {
      dayRangeBar = select(dayRange.value)

      dayRangeBar.append('text')
        .text(formatMoneyNoSymbol(low.value))
        .attr('x', 0)
        .attr('y', 15)
        .attr('stroke', 'rgb(180, 198, 239)')
        .attr('fill', 'rgb(180, 198, 239)')
        .style('font-variant-numeric', 'tabular-nums')
        .style('font-size', '0.8rem')
        .attr('id', 'low')

      const highText = dayRangeBar.append('text')
        .text(formatMoneyNoSymbol(high.value))
        .attr('x', 0)
        .attr('y', 15)
        .attr('stroke', 'rgb(180, 198, 239)')
        .attr('fill', 'rgb(180, 198, 239)')
        .style('font-variant-numeric', 'tabular-nums')
        .style('font-size', '0.8rem')
        .attr('id', 'high')

      const x = (width?.value ?? 0) - (highText.node()?.getComputedTextLength() ?? 0)

      highText.attr('x', x)

      dayRangeBar
        .append('rect')
        .attr('x', 0)
        .attr('y', 20)
        .attr('width', width?.value ?? 0)
        .attr('height', 5)
        .attr('fill', '#f96a6a')
        .attr('id', 'dayRangeBar')

      const triangle = symbol().type(symbolTriangle).size(75)

      dayRangeBar.append('path')
        .attr('d', triangle)
        .attr('stroke', 'red')
        .attr('fill', 'red')
        .attr('transform', `translate(${scale.value?.(currentPrice.value) ?? 0}, 35)`)
        .attr('id', 'currentPrice')
    }
  })
})
</script>

<template>
  <svg
    ref="dayRange"
    :id="snapshot?.lt?.S"
    height="60"
    :width="width"
  ></svg>
</template>


<style scoped>

</style>
