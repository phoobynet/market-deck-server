<script
  lang="ts"
  setup
>
import { computed, inject } from 'vue'
import { SnapshotKey } from '@/components/ReportCard/injectionKeys'
import { orderBy } from 'lodash'
import { SnapshotChange } from '@/types'

const snapshot = inject(SnapshotKey)

const dates = computed(() => {
  if (!snapshot?.value?.changes) {
    return []
  }

  return orderBy(Object.keys(snapshot.value.changes), ['dates'])
})

const changes = computed<SnapshotChange[]>(() => {
  if (!snapshot?.value?.changes) {
    return []
  }

  return orderBy(snapshot?.value?.changes, ['since'], ['asc'])
})
</script>

<template>
  <div>
    <template
      v-for="change in changes"
      :key="change.since"
    >
      <pre class="text-[10px]">{{ JSON.stringify(change, null, 2) }}</pre>
    </template>
  </div>
</template>

<style
  lang="scss"
  scoped
>
</style>
