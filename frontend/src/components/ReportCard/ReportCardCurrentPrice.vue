<script
  lang="ts"
  setup
>
import { computed, inject, ref } from 'vue'
import {
  ChangeSincePreviousKey, CurrentPriceKey, PercentChangeKey, PriceChangeKey, SignSymbolKey, SnapshotKey,
} from '@/components/ReportCard/injectionKeys'
import { formatMoney } from '@/libs/helpers/formatMoney'
import ReportCardChangeIndicator from '@/components/ReportCard/ReportCardChangeIndicator.vue'

const snapshot = inject(SnapshotKey)
const changeSincePrevious = inject(ChangeSincePreviousKey)
const currentPrice = inject(CurrentPriceKey)

const priceChange = inject(PriceChangeKey)
const percentChange = inject(PercentChangeKey)
const signSymbol = inject(SignSymbolKey)

const formatted = computed(() => {
  return formatMoney(currentPrice?.value ?? 0)
})

const priceChangeRef = ref<HTMLDivElement>()
</script>

<template>
  <div
    class="current-price"
    :data-sign="signSymbol"
  >
    <div class="formatted">
      {{ formatted }}
    </div>
    <div class="change">
      <div class="indicator">
        <ReportCardChangeIndicator :sign-symbol="signSymbol" />
      </div>
      <div
        class="price"
        ref="priceChangeRef"
      >
        {{ priceChange }}
      </div>
      <div class="percent">
        ({{ percentChange }})
      </div>
    </div>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .current-price {
    @apply tabular-nums font-semibold flex gap-2 flex-col items-center justify-end transition-all;

    .formatted {
      @apply text-lg sm:text-xl md:text-2xl lg:text-3xl;
    }

    &[data-sign="+"] {
      @apply text-up;

      .change {
        .indicator {
          @apply translate-y-1.5;
        }
      }
    }

    &[data-sign="-"] {
      @apply text-down;

      .change {
        .indicator {
          @apply translate-y-2;
        }
      }
    }

    .change {
      @apply flex gap-2 items-center justify-between -translate-y-0.5;
      .price {
        @apply text-lg font-normal self-end leading-relaxed;
      }

      .percent {
        @apply text-lg font-normal self-end leading-relaxed;
      }
    }
  }
</style>
