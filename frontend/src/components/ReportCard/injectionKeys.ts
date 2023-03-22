import { InjectionKey, Ref } from 'vue'
import { Asset, Snapshot, SnapshotChange } from '@/types'

export const SymbolKey = Symbol('Symbol') as InjectionKey<string>
export const AssetKey = Symbol('Asset') as InjectionKey<Ref<Asset | undefined>>
export const SnapshotKey = Symbol('Snapshot') as InjectionKey<Ref<Snapshot | undefined>>
export const ChangeSincePreviousKey = Symbol('ChangeSincePrevious') as InjectionKey<Ref<SnapshotChange | undefined>>
export const WidthKey = Symbol('Width') as InjectionKey<Ref<number | undefined>>
export const HeightKey = Symbol('Height') as InjectionKey<Ref<number | undefined>>
export const IntradayHighKey = Symbol('IntradayHigh') as InjectionKey<Ref<number>>
export const IntradayLowKey = Symbol('IntradayLow') as InjectionKey<Ref<number>>
export const CurrentPriceKey = Symbol('CurrentPrice') as InjectionKey<Ref<number>>
export const PriceChangeKey = Symbol('PriceChange') as InjectionKey<Ref<string>>
export const PercentChangeKey = Symbol('PercentChange') as InjectionKey<Ref<string>>
export const SignSymbolKey = Symbol('SignSymbol') as InjectionKey<Ref<string>>

