<script
  lang="ts"
  setup
>
import { storeToRefs } from 'pinia'
import { useCalendarDayUpdateStore } from '@/stores/useCalendarDayUpdateStore'

const calendarDayUpdateStore = useCalendarDayUpdateStore()

const {
  condition,
  conditionDescription,
} = storeToRefs(calendarDayUpdateStore)
</script>

<template>
  <div
    class="calendar-day-update"
  >
    <div class="label">condition</div>
    <div
      class="value"
      :data-status="condition"
    >{{ conditionDescription }}
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .calendar-day-update {
    @apply border border-slate-700 rounded text-sm leading-[0.5rem] overflow-hidden flex justify-between items-center gap-2;

    .label {
      @apply text-primary uppercase bg-slate-700 px-2 font-bold;
    }

    .value {
      @apply uppercase font-semibold tracking-wider px-2;

      &[data-status='open'] {
        @apply text-green-400;
      }

      &[data-status='closed_for_the_day'], &[data-status='closed_today'] {
        @apply text-slate-600;
      }

      &[data-status='closed_opening_later'] {
        @apply text-slate-700;
      }

      &[data-status='pre_market'], &[data-status='post_market'] {
        @apply text-purple-500;
      }
    }
  }
</style>
