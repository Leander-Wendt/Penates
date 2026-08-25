import axios, { type InternalAxiosRequestConfig } from 'axios'

import router from '@/router'
import { useAuthStore } from '@/stores/auth'

const client = axios.create({
  baseURL: '/api/v1',
})

/**
 * Attaches the bearer token from the in-memory auth store to every outgoing
 * request, when a token is present.
 */
function attachAuthorizationHeader(config: InternalAxiosRequestConfig): InternalAxiosRequestConfig {
  const authStore = useAuthStore()
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }
  return config
}

client.interceptors.request.use(attachAuthorizationHeader)

client.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error?.response?.status
    if (status === 401) {
      const authStore = useAuthStore()
      authStore.clear()
      router.push('/login')
    }
    return Promise.reject(error)
  },
)

export default client
