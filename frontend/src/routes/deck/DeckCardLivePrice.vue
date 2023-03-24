<script
  lang="ts"
  setup
>
import { computed, ref, watch } from 'vue'
import { formatMoneyNoSymbol } from '@/libs/helpers/formatMoney'

const props = defineProps<{
  price?: number
}>()

const formatted = computed(() => {
  return formatMoneyNoSymbol(props.price)
})

const price = ref<HTMLSpanElement>()

const priceChangeDirection = ref<string>('')

const id = ref<number>(0)

watch(() => props.price, (newValue, oldValue) => {
  if (newValue !== undefined && oldValue !== undefined) {

    if (formatMoneyNoSymbol(newValue) !== formatMoneyNoSymbol(oldValue)) {
      id.value = 0
    } else {
      id.value = id.value + 1
    }

    priceChangeDirection.value = newValue > oldValue
      ? `up`
      : newValue < oldValue
        ? `down`
        : ''
  }
})


</script>

<template>
  <span
    :key="formatted"
    :class="`${priceChangeDirection}_${id}`"
    ref="price"
  >
    {{ formatted }}
  </span>
</template>

<style
  lang="scss"
  scoped
>

  @keyframes up {
    0% {
      color: #0df70d;
    }
    100% {
      color: inherit;
    }
  }

  @keyframes down {
    0% {
      color: #fe6d6d;
    }
    100% {
      color: inherit;
    }
  }


  span[class^="up_"] {
    animation: up .9s;
  }

  span[class^="down_"] {
    animation: down .9s;
  }
</style>
