<script
    lang="ts"
    setup
>
import { computed } from 'vue'
import numeral from 'numeral'

const props = defineProps<{
  amount?: number
  showSign?: boolean
  currency?: string
  sign?: number
  change?: number
}>()

const signSymbol = computed(() => {
  if ((props.amount ?? 0) > 0) {
    return '+'
  }

  if ((props.amount ?? 0) < 0) {
    return '-'
  }

  return ''
})

const formatted = computed(() => {
  if (!props.amount) {
    return ''
  }

  return `${props.currency ?? ''}${signSymbol.value}${numeral(Math.abs(props.amount)).format('0,0.00')}`
})

</script>

<template>
  <span
      class="tabular-nums"
      :data-sign="sign">{{ formatted }}</span>
</template>
