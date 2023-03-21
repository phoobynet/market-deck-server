import { InjectionKey, Ref } from 'vue'
import { Asset, Snapshot, SnapshotChange } from '@/types'

export const AssetKey = Symbol('Asset') as InjectionKey<Ref<Asset | undefined>>
export const SnapshotKey = Symbol('Snapshot') as InjectionKey<Ref<Snapshot | undefined>>
export const ChangeSincePreviousKey = Symbol('ChangeSincePrevious') as InjectionKey<Ref<SnapshotChange | undefined>>
