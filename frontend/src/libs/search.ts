import axios from 'axios'
import { baseUrl } from '@/libs/baseUrl'
import { Asset } from '@/types'

const http = axios.create({
  baseURL: `${baseUrl}/api`,
})

export const search = async (query: string) => {
  return http.get<Asset[]>(`/symbols/query?query=${query}`).then(r => r.data)
}
