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

const formatted = computed(() => {
  if (!props.amount) {
    return ''
  }

  let sign = ''
  if (props.sign) {
    sign = props.sign === -1
      ? '-'
      : '+'
  }

  return `${props.currency ?? ''}${sign}${numeral(Math.abs(props.amount)).format('0,0.00')}`
})

const classes = computed(() => {
  return {
    'tabular-nums': true,
  }
})

const sign = computed(() => {
  if (!props.amount) {
    return ''
  }

  return props.amount! > 0
    ? '+'
    : '-'
})
</script>

<template>
  <span
    :class="classes"
    :data-sign="sign"
  >{{ formatted }}</span>
</template>
