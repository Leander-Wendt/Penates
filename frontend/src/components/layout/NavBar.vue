<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

import { useAuth } from '@/composables/useAuth'
import LanguageSwitcher from './LanguageSwitcher.vue'

/**
 * Primary application navigation. Shows Admin/Logistics-only links only
 * when the current user holds one of those roles, and exposes the
 * language switcher and logout action.
 */
const { isAuthenticated, hasRole, logout } = useAuth()
const { t } = useI18n()
const router = useRouter()

function handleLogout(): void {
  logout()
  router.push('/login')
}
</script>

<template>
  <nav
    v-if="isAuthenticated"
    aria-label="Main"
    class="flex flex-wrap items-center justify-between gap-4 border-b-3 border-brut-black bg-brut-white px-4 py-3"
  >
    <span class="text-lg font-extrabold uppercase">{{ $t('app.name') }}</span>
    <ul class="flex flex-wrap items-center gap-4 font-bold uppercase">
      <li>
        <RouterLink to="/items">{{ t('nav.items') }}</RouterLink>
      </li>
      <li>
        <RouterLink to="/loan-requests">{{ t('nav.loanRequests') }}</RouterLink>
      </li>
      <li v-if="hasRole(['admin', 'logistics'])">
        <RouterLink to="/users">{{ t('nav.users') }}</RouterLink>
      </li>
      <li v-if="hasRole(['admin', 'logistics'])">
        <RouterLink to="/organisations">{{ t('nav.organisations') }}</RouterLink>
      </li>
    </ul>
    <div class="flex items-center gap-4">
      <LanguageSwitcher />
      <button
        type="button"
        class="border-3 border-brut-black bg-brut-white px-3 py-1 font-bold uppercase"
        @click="handleLogout"
      >
        {{ t('nav.logout') }}
      </button>
    </div>
  </nav>
</template>
