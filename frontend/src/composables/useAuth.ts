import { storeToRefs } from 'pinia'

import { useAuthStore } from '@/stores/auth'
import type { Role } from '@/types'

/**
 * Exposes the current authentication state and a role-check helper, backed
 * by the in-memory `auth` Pinia store.
 */
export function useAuth() {
  const authStore = useAuthStore()
  const { isAuthenticated, role, userId, token } = storeToRefs(authStore)

  /**
   * Returns true when the currently authenticated user holds one of the
   * given roles.
   */
  function hasRole(roles: Role[]): boolean {
    if (!role.value) {
      return false
    }
    return roles.includes(role.value)
  }

  /**
   * Redirects to the SSO provider for authentication.
   * This will navigate away from the app to the OIDC provider.
   */
  function loginWithSSO(): void {
    // Get the SSO redirect URL from the backend config
    const apiBase = import.meta.env.VITE_API_PROXY_TARGET || ''
    const redirectUrl = `${apiBase}/api/v1/auth/oauth/redirect`
    window.location.href = redirectUrl
  }

  return {
    isAuthenticated,
    role,
    userId,
    token,
    hasRole,
    login: authStore.login,
    loginWithSSO,
    setToken: authStore.setToken,
    logout: authStore.clear,
  }
}
