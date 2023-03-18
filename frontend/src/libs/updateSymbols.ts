import { http } from '@/libs/http'

export const updateSymbols = async (symbols: string[]): Promise<void> => {
  await http.post(`/symbols?symbols=${symbols.join(',')}`, symbols)
}
