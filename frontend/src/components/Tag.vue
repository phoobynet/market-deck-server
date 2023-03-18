<template>
  <div class="tag">
    <span>
      {{ tag }}
    </span>
    <template v-if="deletable">
      <icon
        class="text-lg cursor-pointer"
        @click="emit('deleted', tag)"
      >
        <CloseCircleOutline />
      </icon>
    </template>
  </div>
</template>

<script
  lang="ts"
  setup
>
import { Icon } from '@vicons/utils'
import { CloseCircleOutline } from '@vicons/ionicons5'
import { PropType } from 'vue'

defineProps({
  tag: {
    type: String as PropType<string>,
    required: true,
  },
  deletable: {
    type: Boolean as PropType<boolean>,
    required: false,
    default: false,
  },
})

const emit = defineEmits(['deleted'])
</script>

<style
  lang="scss"
  scoped
>
  .tag {
    @apply pl-1 pr-0.5 my-0.5 text-base font-semibold border border-accent flex flex-row items-center space-x-1 rounded overflow-hidden;
  }
</style>
