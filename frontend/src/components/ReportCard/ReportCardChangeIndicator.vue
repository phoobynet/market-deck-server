<script
  lang="ts"
  setup
>
import { Icon } from '@vicons/utils'
import { CaretDown, CaretUp } from '@vicons/fa'

defineProps<{
  sign?: number
}>()
</script>

<template>
  <div class="flex flex-col gap-0 transition-all">
    <div v-show="sign === 1">
      <Icon>
        <CaretUp />
      </Icon>
    </div>
    <div v-show="sign === -1" class="translate-y-1">
      <Icon>
        <CaretDown />
      </Icon>
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
