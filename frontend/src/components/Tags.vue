<template>
  <div class="w-full">
    <div class="tags" :class="{ 'has-focus': hasFocus }">
      <transition-group
        enter-active-class="animate__animated animate__fadeIn animate__faster"
        leave-active-class="animate__animated animate__fadeOut animate__faster"
      >
        <tag
          v-for="tag in tags"
          :key="tag"
          :deletable="true"
          :tag="tag"
          @deleted="onRemoveTag"
        />
      </transition-group>
      <div class="relative flex-grow">
        <input
          ref="inputRef"
          v-model="tagInputValue"
          @keydown.esc="onEscape"
          @keydown.enter="onProcessTags"
          @keydown.delete="onDelete"
          @keydown.tab="onProcessTags"
          @keydown.down="onDown"
          @keydown.up="onUp"
          @blur="onProcessTags"
          @input="onInput"
          class="tags-input"
          :class="{ 'has-focus': hasFocus }"
          :placeholder="placeholder"
          v-bind="attrs"
          @focusin="hasFocus = true"
          @focusout="hasFocus = false"
        />
        <div
          class="absolute flex flex-col space-y-1 text-xs overflow-auto bg-slate-800 border border-orange-400"
          :style="{ left: optionsPosition }"
          v-if="tagInputValue.trim().length > 0"
        >
          <div
            v-for="(option, i) in optionsMatchingFilter"
            class="px-2"
            :class="{
              'bg-orange-400 text-black': i === selectedPosition,
            }"
          >
            {{ option.label }}
          </div>
        </div>
      </div>
      <div
        class="opacity-50 pr-2 cursor-pointer"
        @click.prevent="clearAll"
        v-if="false"
      >
        <icon>
          <close></close>
        </icon>
      </div>
    </div>
    <div
      ref="measureTextDiv"
      style="
        position: absolute;
        visibility: hidden;
        height: auto;
        width: auto;
        white-space: nowrap;
      "
    >
      {{ tagInputValue }}
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, PropType, ref, useAttrs, watch } from 'vue'
import { uniq } from 'lodash'
import Tag from '@/components/Tag.vue'
import { Icon } from '@vicons/utils'
import { Close } from '@vicons/ionicons5'
import Fuse from 'fuse.js'
import { debouncedWatch } from '@vueuse/core'

let fuse: Fuse<TagInputOption>

const props = defineProps({
  tags: {
    type: Array as PropType<string[]>,
    required: false,
    default: [],
  },
  placeholder: {
    type: String as PropType<string>,
    required: false,
    default: 'Enter tag and press Space or Enter',
  },
  options: {
    type: Map as PropType<Map<string, string>>,
    required: true,
  },
})

interface TagInputOption {
  tag: string
  label: string
}

const inputRef = ref<HTMLInputElement>()
const measureTextDiv = ref<HTMLDivElement>()

const attrs = useAttrs()

const emit = defineEmits(['change'])

const tagInputValue = ref<string>('')
const optionsPosition = ref<string>('0px')
const optionsMatchingFilter = ref<TagInputOption[]>()
const tagInputOptions = ref<TagInputOption[]>([])
const selectedPosition = ref<number>(-1)
const selectedValue = ref<string>('')

const search = (value: string): TagInputOption[] => {
  if (!props.options) {
    return []
  }

  if (!fuse) {
    throw new Error('Fuse not initialized!')
  }

  const results = fuse.search(value, {
    limit: 20,
  })

  return results.flatMap((value) => value.item)
}

debouncedWatch(
  tagInputValue,
  (newValue) => {
    optionsMatchingFilter.value = search(newValue)
  },
  {
    debounce: 500,
  },
)

const onInput = (event: InputEvent) => {
  optionsPosition.value = (measureTextDiv.value?.clientWidth ?? 0) + 'px'
}

const onRemoveTag = (tag: string) => {
  const t = props.tags.filter((x) => x !== tag)
  emit('change', t)
}

const onEscape = () => {
  tagInputValue.value = ''
}

const hasFocus = ref<boolean>()

const onProcessTags = () => {
  let addedTags: string[] = []

  if (selectedValue.value) {
    addedTags.push(selectedValue.value)
    selectedValue.value = ''
  } else {
    addedTags = tagInputValue.value
      .trim()
      .split(' ')
      .map((v) => (v || '').trim().toUpperCase())
      .filter((v) => !!v)
      .filter((v) => props.options.has(v))
  }

  if (addedTags.length > 0) {
    emit('change', uniq([...props.tags, ...addedTags]))
  }

  tagInputValue.value = ''
}

watch(optionsMatchingFilter, (newValue) => {
  if (!newValue?.length) {
    selectedPosition.value = -1
  }
})

const onDown = () => {
  if (optionsMatchingFilter.value?.length) {
    let nextPosition = selectedPosition.value + 1

    if (nextPosition <= optionsMatchingFilter.value.length - 1) {
      selectedPosition.value = nextPosition
      selectedValue.value =
        optionsMatchingFilter.value[selectedPosition.value].tag
    }
  } else {
    selectedPosition.value = -1
  }
}

const onUp = () => {
  if (optionsMatchingFilter.value?.length) {
    if (selectedPosition.value > 0) {
      selectedPosition.value = selectedPosition.value - 1
      selectedValue.value =
        optionsMatchingFilter.value[selectedPosition.value].tag
    }
  } else {
    selectedPosition.value = -1
  }
}

const onDelete = () => {
  if (!tagInputValue.value) {
    const t = [...props.tags]

    t.pop()

    emit('change', t)
  }
}

const clearAll = () => {
  tagInputValue.value = ''
  onProcessTags()
  emit('change', [])
}

onMounted(() => {
  const input = inputRef.value as HTMLInputElement
  const style = window
    .getComputedStyle(input, null)
    .getPropertyValue('font-size')
  const div = measureTextDiv.value as HTMLDivElement
  div.style.fontSize = style

  if (props.options) {
    tagInputOptions.value = Array.from(props.options.entries()).map(
      (entry: string[]) => ({
        tag: entry[0],
        label: `${entry[0]} - ${entry[1]}`,
      }),
    )
    fuse = new Fuse<TagInputOption>(tagInputOptions.value, {
      isCaseSensitive: false,
      keys: [
        {
          name: 'tag',
          weight: 1,
        },
        { name: 'label', weight: 0.5 },
      ],
      shouldSort: true,
      includeScore: true,
    })
  }
})
</script>

<style lang="scss" scoped>
  .tags {
    @apply flex flex-row flex-wrap space-x-2 w-full items-center h-8 bg-slate-800 pl-1 border border-slate-900;

    &.hover {
      @apply border-orange-400;
    }

    &.has-focus {
      @apply border border-orange-400 bg-slate-700;
    }
  }

  .tags-input {
    @apply outline-none pl-0.5 w-full h-8 bg-slate-800 w-full h-full outline-none placeholder-slate-500;

    &.has-focus {
      @apply bg-slate-700;
    }
  }
</style>
