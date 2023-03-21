<script
  lang="ts"
  setup
>
import { computed, inject } from 'vue'
import { SnapshotKey } from '@/components/ReportCard/injectionKeys'
import { orderBy } from 'lodash'
import { SnapshotChange } from '@/types'
import { formatPercent } from '@/libs/helpers/formatPercent'

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

const changeLabel = (label: string): string => {
  return label
    .replace('ago', '')
    .replace('year', 'y')
    .replace('months', 'm')
    .replace('month', 'm')
    .replace('weeks', 'w')
    .replace('week', 'w')
    .replace(' ', '')
    .trim()
    .toUpperCase()
}

const signSymbol = (sign: number): string => {
  if (sign === 1) {
    return '+'
  } else if (sign === -1) {
    return '-'
  }

  return ''
}
</script>

<template>
  <table class="w-full text-xxs">
    <thead>
      <tr>
        <th class="text-left">When</th>
        <th class="text-right">Change</th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="change of changes"
        :key="change.since"
      >
        <td>{{ change.label }}</td>
        <td class="text-right tabular-nums">{{ signSymbol(change.cs) }}{{ formatPercent(Math.abs(change.cp)) }}</td>
      </tr>
    </tbody>
  </table>
</template>

<style
  lang="scss"
  scoped
>
  .change {
    @apply flex gap-0.5 items-center justify-between text-xs leading-none text-white px-0.5 py-0 rounded-sm;


    &[data-sign="1"] {
      @apply bg-up-dark text-white;
    }

    &[data-sign="-1"] {
      @apply bg-down-dark text-white;
    }

    .label {
      @apply leading-none;
    }

    .percentage {
      @apply tabular-nums font-bold;
    }
  }

</style>
