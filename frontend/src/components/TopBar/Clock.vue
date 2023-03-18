<script
  lang="ts"
  setup
>
import { computed } from 'vue'
import { format, utcToZonedTime } from 'date-fns-tz'

const DATE_PART_FORMAT = 'E do LLL'
const TIME_PART_FORMAT = 'HH:mm:ss'

const props = defineProps<{
  label: string,
  timeUnixMs: number
  tz: string
}>()

type TimeParts = {
  date: string
  time: string
}

const timeParts = computed<TimeParts>(() => {
  const t = utcToZonedTime(new Date(props.timeUnixMs), props.tz)

  return {
    date: format(t, DATE_PART_FORMAT),
    time: format(t, TIME_PART_FORMAT),
  }
})

</script>

<template>
  <div class="time">
    <div class="label">{{ label }}</div>
    <div class="value">
      <div class="date-part">{{ timeParts.date }}</div>
      <div class="time-part">{{ timeParts.time }}</div>
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .time {
    @apply border border-slate-700 rounded text-sm leading-[0.5rem] overflow-hidden;

    @apply flex justify-between items-center gap-2;

    .label {
      @apply text-primary uppercase bg-slate-700 px-2 font-bold;
    }

    .value {
      @apply tabular-nums tracking-wider pr-2 flex gap-2 justify-between items-center;

      .date-part {
        @apply font-light;
      }

      .time-part {
        @apply font-semibold text-primary;
      }
    }
  }
</style>
