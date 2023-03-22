<script
  lang="ts"
  setup
>
import { Icon } from '@vicons/utils'
import { CaretDown, CaretUp, Minus } from '@vicons/fa'

defineProps<{
  signSymbol?: string
}>()
</script>

<template>
  <div class="flex flex-col gap-0 transition-all items-center">
    <div v-show="signSymbol === '+'">
      <Icon :size="30">
        <CaretUp />
      </Icon>
    </div>
    <div
      v-show="signSymbol === '-'"
    >
      <Icon :size="30">
        <CaretDown />
      </Icon>
    </div>
    <div
      v-show="signSymbol === ''"
      class="opacity-50 uppercase"
    >
      <Icon :size="30">
        <Minus />
      </Icon>
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
