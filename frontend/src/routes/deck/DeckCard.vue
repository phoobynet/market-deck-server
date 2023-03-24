<script
  lang="ts"
  setup
>
import { useDeckSnapshot } from '@/routes/deck/useDeckSnapshot'
import { Icon } from '@vicons/utils'
import { ChartAverage } from '@vicons/carbon'
import { ArrowBigDown, ArrowBigTop } from '@vicons/tabler'
import DeckCardLivePrice from '@/routes/deck/DeckCardLivePrice.vue'

const props = defineProps<{
  symbol: string
}>()

const {
  snapshot,
  price,
  companyName,
  priceChange,
  changePercentAbs,
  sign,
  prevClose,
  dailyHigh,
  dailyLow,
  dailyVolume,
  avgVolume,
  ytdChangePercentAbs,
  ytdSign,
} = useDeckSnapshot(props.symbol)

</script>

<template>
  <div
    v-if="snapshot"
    class="deck-card"
  >
    <div class="company-info">
      <div class="symbol">
        <div class="symbol">{{ snapshot?.symbol }}</div>
        <div class="exchange">{{ snapshot?.exchange }}</div>
      </div>
      <div class="name">{{ companyName }}</div>
    </div>
    <div class="current-price">
      <div class="currency-symbol">$</div>
      <div class="amount">
        <DeckCardLivePrice :price="snapshot?.price" />
      </div>
    </div>
    <div class="previous-close">
      <div class="label">Prev Close</div>
      <div class="value">{{ prevClose }}</div>
    </div>
    <div
      class="price-change"
      :data-sign="sign"
    >
      <div class="change-amount">
        <Icon class="translate-y-0.5">
          <ArrowBigDown v-if="ytdSign === '-'" />
          <ArrowBigTop v-else-if="ytdSign === '+'" />
        </Icon>
        {{ priceChange }}
      </div>
      <div class="change-percent">
        ({{ changePercentAbs }})
      </div>
    </div>
    <div class="day-range hidden">
      day range
    </div>

    <dl class="key-info">
      <div>
        <dt>Vol</dt>
        <dd>{{ dailyVolume }}</dd>
      </div>
      <div>
        <dt>
          Avg Vol
        </dt>
        <dd>{{ avgVolume }}</dd>
      </div>
      <div>
        <dt>YTD Chg</dt>
        <dd :data-sign="ytdSign">
          <Icon
            size="12"
            class="translate-y-[1px]"
          >
            <ArrowBigDown v-if="ytdSign === '-'" />
            <ArrowBigTop v-else-if="ytdSign === '+'" />
          </Icon>
          {{ ytdChangePercentAbs }}
        </dd>
      </div>
    </dl>
  </div>
</template>

<style
  lang="scss"
  scoped
>
  .deck-card {
    @apply bg-slate-900 p-2 grid;

    grid-template-columns: 1fr 1fr;
    grid-template-rows: repeat(4, auto);

    grid-template-areas:
      "company-info current-price"
      "previous-close price-change"
      "key-info key-info"
      "day-range day-range";

    .company-info {
      grid-area: company-info;
      @apply flex flex-col justify-items-start;

      .symbol {
        @apply flex text-lg;

        :first-child {
          @apply text-orange-400;
        }

        :last-child {
          &::before {
            content: ":"
          }
        }

        .exchange {
          @apply font-light text-slate-400 pr-1;
        }

        .symbol {
          @apply font-bold tracking-widest;
        }
      }

      .name {
        @apply font-normal text-slate-300 text-xxxs truncate;
        grid-area: name;
      }
    }

    .current-price {
      @apply flex gap-0.5 justify-end;
      grid-area: current-price;

      .currency-symbol {
        @apply translate-y-0.5;
      }

      .amount {
        @apply text-2xl tabular-nums;
      }
    }

    .price-change {
      @apply text-xxs flex gap-1 justify-end tabular-nums;
      grid-area: price-change;

      &[data-sign="+"] {
        @apply text-up;
      }

      &[data-sign="-"] {
        @apply text-down;
      }
    }

    .previous-close {
      @apply flex text-xxs justify-start gap-1 items-center tabular-nums;
      grid-area: previous-close;

      .label {
        @apply text-slate-400;
        &::after {
          content: ':'
        }
      }
    }

    .day-range {
      grid-area: day-range;
      @apply w-full bg-slate-700 text-xxs;
    }

    .key-info {
      grid-area: key-info;
      @apply w-full grid grid-cols-4 items-center gap-1 justify-between mt-2 text-sm;

      grid-template-columns: repeat(3, 1fr);
      grid-template-rows: 2rem auto;

      dt {
        @apply uppercase font-light text-center bg-slate-700 text-[0.7rem] leading-snug;
      }

      dd {
        @apply tabular-nums text-center text-orange-400;

        &[data-sign="+"] {
          @apply text-up;
        }

        &[data-sign="-"] {
          @apply text-down;
        }
      }
    }
  }
</style>
