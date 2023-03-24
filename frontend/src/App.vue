<script
  lang="ts"
  setup
>
import { onBeforeUnmount, onMounted } from 'vue'
import TopBar from '@/components/TopBar/TopBar.vue'
import { CalendarDayUpdateListener } from '@/libs/CalendarDayUpdateListener'
import { useAssetsStore } from '@/stores/useAssetsStore'

const assetsStore = useAssetsStore()

const calendarDayUpdateListener = new CalendarDayUpdateListener()

onMounted(async () => {
  await assetsStore.fetch()
  calendarDayUpdateListener.start()
})

onBeforeUnmount(() => {
  calendarDayUpdateListener.close()
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
