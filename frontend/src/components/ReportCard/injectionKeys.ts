import { InjectionKey, Ref } from 'vue'
import { Asset, Snapshot } from '@/types'

export const AssetKey = Symbol('Asset') as InjectionKey<Ref<Asset | undefined>>
export const SnapshotKey = Symbol('Snapshot') as InjectionKey<Ref<Snapshot | undefined>>
export const SignKey = Symbol('Sign') as InjectionKey<Ref<number | undefined>>
export const SignSymbolKey = Symbol('SignSymbol') as InjectionKey<Ref<string>>
