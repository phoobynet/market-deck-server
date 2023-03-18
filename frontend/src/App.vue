<script
  lang="ts"
  setup
>
import { onBeforeUnmount, onMounted } from 'vue'
import TopBar from '@/components/TopBar/TopBar.vue'
import { CalendarDayUpdateListener } from '@/libs/CalendarDayUpdateListener'
import { RealtimeSymbolsListener } from '@/libs/RealtimeSymbolsListener'

const calendarDayUpdateListener = new CalendarDayUpdateListener()
const realtimeSymbolsListener = new RealtimeSymbolsListener()

onMounted(async () => {
  calendarDayUpdateListener.start()
  realtimeSymbolsListener.start()
})

onBeforeUnmount(() => {
  calendarDayUpdateListener.close()
  realtimeSymbolsListener.close()
})
</script>
<template>
  <div>
    <nav>
      <TopBar />
    </nav>
    <RouterView></RouterView>
  </div>
</template>
