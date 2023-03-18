import { Asset } from '@/types'
import { http } from '@/libs/http'

export const search = async (query: string) => {
  return http.get<Asset[]>(`/symbols/query?query=${query}`).then(r => r.data)
}
