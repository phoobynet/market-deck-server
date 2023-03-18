import { http } from '@/libs/http'

export const getSymbols = async (): Promise<string[]> => {
  return http.get<{ symbols: string[] }>('/symbols').then((res) => res.data.symbols)
}
