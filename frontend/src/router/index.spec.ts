import { describe, expect, it, vi, beforeEach } from 'vitest'
import type { RouteLocationNormalized } from 'vue-router'

import { useAuth } from '@/composables/useAuth'

vi.mock('@/composables/useAuth', () => ({
  useAuth: vi.fn(),
}))

function makeRoute(overrides: Partial<RouteLocationNormalized> = {}): RouteLocationNormalized {
  return {
    fullPath: '/items',
    meta: {},
    ...overrides,
  } as RouteLocationNormalized
}

describe('authGuard', () => {
  beforeEach(() => {
    vi.resetModules()
  })

  it('redirects to login when the route is not public and the user is unauthenticated', async () => {
    ;(useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
      isAuthenticated: { value: false },
      hasRole: vi.fn(() => false),
    })
    const { authGuard } = await import('./index')
    const next = vi.fn()

    authGuard(makeRoute({ meta: {} }), makeRoute(), next)

    expect(next).toHaveBeenCalledWith(expect.objectContaining({ name: 'login' }))
  })

  it('allows a public route for an unauthenticated user', async () => {
    ;(useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
      isAuthenticated: { value: false },
      hasRole: vi.fn(() => false),
    })
    const { authGuard } = await import('./index')
    const next = vi.fn()

    authGuard(makeRoute({ meta: { public: true } }), makeRoute(), next)

    expect(next).toHaveBeenCalledWith()
  })

  it('redirects to forbidden when authenticated but missing the required role', async () => {
    ;(useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
      isAuthenticated: { value: true },
      hasRole: vi.fn(() => false),
    })
    const { authGuard } = await import('./index')
    const next = vi.fn()

    authGuard(makeRoute({ meta: { roles: ['admin'] } }), makeRoute(), next)

    expect(next).toHaveBeenCalledWith(expect.objectContaining({ name: 'forbidden' }))
  })

  it('proceeds when authenticated and holding the required role', async () => {
    ;(useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
      isAuthenticated: { value: true },
      hasRole: vi.fn(() => true),
    })
    const { authGuard } = await import('./index')
    const next = vi.fn()

    authGuard(makeRoute({ meta: { roles: ['admin'] } }), makeRoute(), next)

    expect(next).toHaveBeenCalledWith()
  })

  it('proceeds when authenticated and the route has no role restriction', async () => {
    ;(useAuth as ReturnType<typeof vi.fn>).mockReturnValue({
      isAuthenticated: { value: true },
      hasRole: vi.fn(() => false),
    })
    const { authGuard } = await import('./index')
    const next = vi.fn()

    authGuard(makeRoute({ meta: {} }), makeRoute(), next)

    expect(next).toHaveBeenCalledWith()
  })
})
