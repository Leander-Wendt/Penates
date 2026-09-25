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
 * Authentication store. Persists JWT and its derived claims to localStorage
 * so the user stays logged in across page refreshes.
 */
const STORAGE_KEY = 'penates-auth'

function isTokenExpired(expiresAt: string | null): boolean {
  if (!expiresAt) {
    return true
  }
  try {
    const expiryDate = new Date(expiresAt)
    const now = new Date()
    // Add a small buffer to account for clock skew
    return expiryDate <= new Date(now.getTime() + 60000)
  } catch {
    return true
  }
}

function loadFromStorage(): { token: string | null; expiresAt: string | null } {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      const parsed = JSON.parse(stored)
      // Only return stored data if token hasn't expired
      if (parsed.expiresAt && !isTokenExpired(parsed.expiresAt)) {
        return parsed
      }
    }
  } catch {
    // Ignore storage errors
  }
  return { token: null, expiresAt: null }
}

function saveToStorage(token: string | null, expiresAt: string | null): void {
  try {
    if (token && expiresAt) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({ token, expiresAt }))
    } else {
      localStorage.removeItem(STORAGE_KEY)
    }
  } catch {
    // Ignore storage errors
  }
}

export const useAuthStore = defineStore('auth', () => {
  const { token: storedToken, expiresAt: storedExpiresAt } = loadFromStorage()
  const token = ref<string | null>(storedToken)
  const expiresAt = ref<string | null>(storedExpiresAt)
  const role = ref<Role | null>(null)
  const userId = ref<string | null>(null)

  // Initialize role and userId from the stored token if present
  if (storedToken) {
    const claims = decodeTokenClaims(storedToken)
    role.value = claims.role ?? null
    userId.value = claims.sub ?? null
  }

  const isAuthenticated = computed(() => token.value !== null)

  /**
   * Authenticates against the backend and, on success, populates the
   * session and persists it to localStorage. On failure, clears any
   * potentially stale auth data.
   */
  async function login(credentials: LoginCredentials): Promise<void> {
    try {
      const response = await apiLogin(credentials)
      const claims = decodeTokenClaims(response.token)
      token.value = response.token
      expiresAt.value = response.expiresAt
      role.value = claims.role ?? null
      userId.value = claims.sub ?? null
      saveToStorage(response.token, response.expiresAt)
    } catch (error) {
      // Clear stale auth data on login failure
      clear()
      throw error
    }
  }

  /**
   * Clears the session from both memory and localStorage. Called on explicit
   * logout and whenever the API client observes a 401 (invalid/expired session).
   */
  function clear(): void {
    token.value = null
    expiresAt.value = null
    role.value = null
    userId.value = null
    saveToStorage(null, null)
  }

  return { token, expiresAt, role, userId, isAuthenticated, login, clear }
})
