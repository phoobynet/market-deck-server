<script
  lang="ts"
  setup
>
import Condition from '@/components/CalendarDayCondition.vue'
import { useCalendarDayUpdateStore } from '@/stores'
import { storeToRefs } from 'pinia'
import Clock from '@/components/TopBar/Clock.vue'

const calendarDayUpdateStore = useCalendarDayUpdateStore()

const {
  timeUnixMs,
  marketTimeZone,
  localTimeZone,
} = storeToRefs(calendarDayUpdateStore)

</script>
<template>
  <div class="top-bar">
    <Condition />
    <div class="flex gap-2 items-center">
      <Clock
        :time-unix-ms="timeUnixMs"
        :tz="marketTimeZone"
        label="market"
      />
      <Clock
        :time-unix-ms="timeUnixMs"
        :tz="localTimeZone"
        label="local"
      />
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .top-bar {
    @apply flex justify-between items-center justify-between px-4 bg-slate-900 h-10;
  }
</style>
