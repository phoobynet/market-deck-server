import axios from 'axios'
import { baseUrl } from '@/libs/baseUrl'

export const http = axios.create({
  baseURL: `${baseUrl}/api`,
})
