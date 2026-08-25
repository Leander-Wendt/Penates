import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { login as apiLogin } from '@/api/auth'
import type { LoginCredentials, Role } from '@/types'

interface DecodedTokenClaims {
  sub?: string
  role?: Role
}

function decodeTokenClaims(token: string): DecodedTokenClaims {
  try {
    const payload = token.split('.')[1]
    const normalized = payload.replace(/-/g, '+').replace(/_/g, '/')
    const json = decodeURIComponent(
      atob(normalized)
        .split('')
        .map((char) => '%' + char.charCodeAt(0).toString(16).padStart(2, '0'))
        .join(''),
    )
    return JSON.parse(json) as DecodedTokenClaims
  } catch {
    return {}
  }
}

/**
 * In-memory authentication store. Holds the current JWT and its derived
 * claims for the lifetime of the page only; nothing is persisted to
 * localStorage, sessionStorage or cookies.
 */
export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const expiresAt = ref<string | null>(null)
  const role = ref<Role | null>(null)
  const userId = ref<string | null>(null)

  const isAuthenticated = computed(() => token.value !== null)

  /**
   * Authenticates against the backend and, on success, populates the
   * in-memory session (token, expiry, and claims decoded from the JWT).
   */
  async function login(credentials: LoginCredentials): Promise<void> {
    const response = await apiLogin(credentials)
    const claims = decodeTokenClaims(response.token)
    token.value = response.token
    expiresAt.value = response.expiresAt
    role.value = claims.role ?? null
    userId.value = claims.sub ?? null
  }

  /**
   * Clears the in-memory session. Called on explicit logout and whenever
   * the API client observes a 401 (invalid/expired session).
   */
  function clear(): void {
    token.value = null
    expiresAt.value = null
    role.value = null
    userId.value = null
  }

  return { token, expiresAt, role, userId, isAuthenticated, login, clear }
})
