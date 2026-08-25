import { beforeEach, describe, expect, it, vi } from 'vitest'
import axios from 'axios'

vi.mock('axios', () => {
  const instance = {
    interceptors: {
      request: { use: vi.fn() },
      response: { use: vi.fn() },
    },
  }
  return {
    default: {
      create: vi.fn(() => instance),
    },
  }
})

const authStoreInstance = {
  token: 'stored-token',
  clear: vi.fn(),
}

vi.mock('@/stores/auth', () => ({
  useAuthStore: vi.fn(() => authStoreInstance),
}))

vi.mock('@/router', () => ({
  default: {
    push: vi.fn(),
  },
}))

describe('api client', () => {
  beforeEach(() => {
    vi.resetModules()
    authStoreInstance.clear.mockClear()
  })

  it('creates an axios instance with the /api/v1 base URL', async () => {
    await import('./client')
    expect(axios.create).toHaveBeenCalledWith(expect.objectContaining({ baseURL: '/api/v1' }))
  })

  it('registers a request interceptor that attaches the bearer token', async () => {
    const { default: client } = await import('./client')
    expect(client.interceptors.request.use).toHaveBeenCalled()
  })

  it('registers a response interceptor for 401/403 handling', async () => {
    const { default: client } = await import('./client')
    expect(client.interceptors.response.use).toHaveBeenCalled()
  })

  it('attaches the Authorization header when a token is present', async () => {
    const { default: client } = await import('./client')
    const requestInterceptor = (client.interceptors.request.use as ReturnType<typeof vi.fn>).mock
      .calls[0][0]
    const config = requestInterceptor({ headers: {} })
    expect(config.headers.Authorization).toBe('Bearer stored-token')
  })

  it('clears the auth store and redirects to /login on 401', async () => {
    const { useAuthStore } = await import('@/stores/auth')
    const router = (await import('@/router')).default
    const { default: client } = await import('./client')
    const errorInterceptor = (client.interceptors.response.use as ReturnType<typeof vi.fn>).mock
      .calls[0][1]

    const authStore = useAuthStore()
    await expect(errorInterceptor({ response: { status: 401 } })).rejects.toBeDefined()

    expect(authStore.clear).toHaveBeenCalled()
    expect(router.push).toHaveBeenCalledWith('/login')
  })

  it('does not clear the auth store on 403', async () => {
    const { useAuthStore } = await import('@/stores/auth')
    const { default: client } = await import('./client')
    const errorInterceptor = (client.interceptors.response.use as ReturnType<typeof vi.fn>).mock
      .calls[0][1]

    const authStore = useAuthStore()
    await expect(errorInterceptor({ response: { status: 403 } })).rejects.toBeDefined()

    expect(authStore.clear).not.toHaveBeenCalled()
  })
})
