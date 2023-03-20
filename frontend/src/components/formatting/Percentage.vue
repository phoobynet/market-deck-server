<script
  lang="ts"
  setup
>
import numeral from 'numeral'
import { computed } from 'vue'

const props = defineProps<{
  amount?: number
}>()

const formatted = computed<string>(() => {
  if (props.amount === undefined) {
    return ''
  }

  return numeral(props.amount).format('0.00%')
})
</script>

<template>
  <span class="tabular-nums">{{formatted}}</span>
</template>

<style
  lang="scss"
  scoped
>
</style>
