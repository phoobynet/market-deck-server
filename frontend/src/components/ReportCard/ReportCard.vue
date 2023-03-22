<script
  lang="ts"
  setup
>
import { provide, ref } from 'vue'
import {
  AssetKey, ChangeSincePreviousKey, CurrentPriceKey, HeightKey, IntradayHighKey, IntradayLowKey, PercentChangeKey,
  PriceChangeKey, SignSymbolKey, SnapshotKey, SymbolKey, WidthKey,
} from '@/components/ReportCard/injectionKeys'
import ReportCardSymbol from '@/components/ReportCard/ReportCardSymbol.vue'
import ReportCardName from '@/components/ReportCard/ReportCardName.vue'
import ReportCardCurrentPrice from '@/components/ReportCard/ReportCardCurrentPrice.vue'
import ReportCardDailyVolume from '@/components/ReportCard/ReportCardDailyVolume.vue'
import ReportCardExchange from '@/components/ReportCard/ReportCardExchange.vue'
import ReportCardDayRange from '@/components/ReportCard/ReportCardDayRange/ReportCardDayRange.vue'
import { useElementSize } from '@vueuse/core'
import { useReportCard } from '@/components/ReportCard/useReportCard'

const props = defineProps<{
  symbol: string
}>()

const {
  snapshot,
  changeSincePrevious,
  asset,
  intradayHigh,
  intradayLow,
  currentPrice,
  priceChange,
  percentChange,
  signSymbol,
} = useReportCard(props.symbol)

const reportCard = ref<HTMLDivElement>()

const {
  height,
  width,
} = useElementSize(reportCard)

provide(AssetKey, asset)
provide(SnapshotKey, snapshot)
provide(ChangeSincePreviousKey, changeSincePrevious)
provide(WidthKey, width)
provide(HeightKey, height)
provide(IntradayHighKey, intradayHigh)
provide(IntradayLowKey, intradayLow)
provide(CurrentPriceKey, currentPrice)
provide(SymbolKey, props.symbol)
provide(PriceChangeKey, priceChange)
provide(PercentChangeKey, percentChange)
provide(SymbolKey, props.symbol)
provide(SignSymbolKey, signSymbol)
</script>

<template>
  <div
    ref="reportCard"
    v-if="symbol && snapshot && asset"
    class="report-card"
    :data-sign="changeSincePrevious?.cs"
  >
    <ReportCardSymbol />
    <ReportCardName class="name" />
    <ReportCardCurrentPrice class="latest-price" />
    <ReportCardDailyVolume class="daily-volume" />
    <ReportCardExchange class="exchange" />
<!--    <ReportCardDayRange class="day-range" />-->
    <pre class="preview text-[11px]">
      {{ JSON.stringify(snapshot, null, 2) }}
    </pre>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .report-card {
    @apply border px-2 py-1 rounded-md border-slate-600 transition-all duration-300 ease-in-out;

    &[data-sign="1"] {
      @apply border-up bg-gradient-to-t from-[#00BFA63C] to-transparent;
    }

    &[data-sign="-1"] {
      @apply border-down bg-gradient-to-t from-[#EB6B6B33] to-transparent;;
    }

    display: grid;
    grid-template-areas:
      "name name exchange"
      "symbol latest-price latest-price"
      "previous-close previous-close previous-close"
      ". . daily-volume"
      "day-range day-range day-range"
      "preview preview preview";

    .symbol {
      grid-area: symbol;
      @apply text-xl tracking-widest;
    }

    .name {
      grid-area: name;
      @apply self-start;
    }

    .latest-price {
      grid-area: latest-price;
      @apply justify-self-end;
    }

    .previous-close {
      grid-area: previous-close;
    }

    .changes {
      grid-area: changes;
      @apply self-start;
    }

    .daily-volume {
      grid-area: daily-volume;
      @apply justify-self-end;
    }

    .exchange {
      grid-area: exchange;
      @apply justify-self-end;
    }

    .day-range {
      grid-area: day-range;
    }

    .preview {
      grid-area: preview;
      display: none;
    }
  }
</style>
