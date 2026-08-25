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

  return {
    isAuthenticated,
    role,
    userId,
    token,
    hasRole,
    login: authStore.login,
    logout: authStore.clear,
  }
}
