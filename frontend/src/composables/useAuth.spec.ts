import { describe, expect, it, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

import { useAuth } from './useAuth'
import { useAuthStore } from '@/stores/auth'

describe('useAuth', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('reports isAuthenticated as false with no token', () => {
    const { isAuthenticated } = useAuth()
    expect(isAuthenticated.value).toBe(false)
  })

  it('reports isAuthenticated as true once a token is set', () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'

    const { isAuthenticated } = useAuth()
    expect(isAuthenticated.value).toBe(true)
  })

  it('hasRole returns true when the current role is in the list', () => {
    const authStore = useAuthStore()
    authStore.role = 'logistics'

    const { hasRole } = useAuth()
    expect(hasRole(['admin', 'logistics'])).toBe(true)
  })

  it('hasRole returns false when the current role is not in the list', () => {
    const authStore = useAuthStore()
    authStore.role = 'student'

    const { hasRole } = useAuth()
    expect(hasRole(['admin', 'logistics'])).toBe(false)
  })

  it('hasRole returns false when there is no current role', () => {
    const { hasRole } = useAuth()
    expect(hasRole(['admin'])).toBe(false)
  })
})
