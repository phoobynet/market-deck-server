<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { select } from 'd3'

const dayRange = ref<HTMLDivElement>()

onMounted(() => {
  select(dayRange)
})
</script>

<template>
  <div ref="dayRange"></div>
</template>


<style scoped>

</style>