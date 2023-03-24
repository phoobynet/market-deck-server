<script
  lang="ts"
  setup
>
import { useCalendarDayUpdateStore } from '@/stores/useCalendarDayUpdateStore'
import { storeToRefs } from 'pinia'

const calendarDayUpdateStore = useCalendarDayUpdateStore()

const {
  marketTimeFormatted,
  conditionDescription,
} = storeToRefs(calendarDayUpdateStore)
</script>

<template>
  <div class="top-bar">
    <div class="condition">{{ conditionDescription }}</div>
    <div class="market-time">{{ marketTimeFormatted }}</div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .top-bar {
    @apply flex justify-between items-center justify-between md:px-4 px-1 bg-slate-900 h-6 md:h-10 bg-black;

    .condition {
      @apply uppercase text-xxs font-bold tracking-widest;
    }

    .market-time {
      @apply text-xxs tabular-nums;
    }
  }
</style>
