import { beforeEach, describe, expect, it, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

import { login as apiLogin } from '@/api/auth'

vi.mock('@/api/auth', () => ({
  login: vi.fn(),
}))

function encodeJwtPayload(payload: Record<string, unknown>): string {
  const header = Buffer.from(JSON.stringify({ alg: 'HS256', typ: 'JWT' })).toString('base64url')
  const body = Buffer.from(JSON.stringify(payload)).toString('base64url')
  return `${header}.${body}.signature`
}

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('starts unauthenticated', async () => {
    const { useAuthStore } = await import('./auth')
    const store = useAuthStore()
    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBeNull()
    expect(store.role).toBeNull()
  })

  it('sets token, expiry, role and userId on successful login', async () => {
    const token = encodeJwtPayload({ sub: 'user-1', role: 'admin' })
    ;(apiLogin as ReturnType<typeof vi.fn>).mockResolvedValue({
      token,
      expiresAt: '2026-01-01T00:00:00Z',
    })

    const { useAuthStore } = await import('./auth')
    const store = useAuthStore()
    await store.login({ email: 'a@b.com', password: 'secret' })

    expect(store.isAuthenticated).toBe(true)
    expect(store.token).toBe(token)
    expect(store.role).toBe('admin')
    expect(store.userId).toBe('user-1')
    expect(store.expiresAt).toBe('2026-01-01T00:00:00Z')
  })

  it('propagates failure and stays unauthenticated on bad credentials', async () => {
    ;(apiLogin as ReturnType<typeof vi.fn>).mockRejectedValue({
      response: { status: 401 },
    })

    const { useAuthStore } = await import('./auth')
    const store = useAuthStore()
    await expect(store.login({ email: 'a@b.com', password: 'wrong' })).rejects.toBeDefined()

    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBeNull()
  })

  it('clears all session state on clear()', async () => {
    const token = encodeJwtPayload({ sub: 'user-1', role: 'student' })
    ;(apiLogin as ReturnType<typeof vi.fn>).mockResolvedValue({
      token,
      expiresAt: '2026-01-01T00:00:00Z',
    })

    const { useAuthStore } = await import('./auth')
    const store = useAuthStore()
    await store.login({ email: 'a@b.com', password: 'secret' })
    store.clear()

    expect(store.isAuthenticated).toBe(false)
    expect(store.token).toBeNull()
    expect(store.role).toBeNull()
    expect(store.userId).toBeNull()
  })
})
