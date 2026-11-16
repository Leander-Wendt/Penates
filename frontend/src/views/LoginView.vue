<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import { useAuth } from '@/composables/useAuth'

/**
 * Login page: authenticates via email/password and redirects to the item
 * list (or the originally requested route) on success. Also supports SSO.
 */
const { t } = useI18n()
const { login, loginWithSSO, setToken } = useAuth()
const router = useRouter()
const route = useRoute()

const form = reactive({ email: '', password: '' })
const errorMessage = ref('')
const isSubmitting = ref(false)
const ssoAvailable = ref(false)

// Check if SSO is available on mount and handle OAuth callback
onMounted(async () => {
  // Check for OAuth callback with token in URL hash
  const hash = window.location.hash.substring(1) // Remove leading #
  const params = new URLSearchParams(hash)
  const token = params.get('token')
  const expiresAt = params.get('expiresAt')

  if (token && expiresAt) {
    // Store token and redirect to items page
    setToken(token, expiresAt)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/items'
    router.push(redirect)
    return // Don't check SSO availability if we just logged in
  }

  // Clear the hash from the URL if there was one (but no valid token)
  if (window.location.hash) {
    window.history.replaceState({}, '', window.location.pathname + window.location.search)
  }

  // Check if SSO is available
  try {
    const apiBase = import.meta.env.VITE_API_PROXY_TARGET || ''
    const response = await fetch(`${apiBase}/api/v1/auth/oauth/config`)
    if (response.ok) {
      const config = await response.json()
      ssoAvailable.value = config.enabled
    }
  } catch {
    // SSO not available or error fetching config
    ssoAvailable.value = false
  }
})

async function handleSubmit(): Promise<void> {
  errorMessage.value = ''
  isSubmitting.value = true
  try {
    await login({ email: form.email, password: form.password })
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/items'
    router.push(redirect)
  } catch {
    errorMessage.value = t('errors.loginFailed')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="flex min-h-[70vh] items-center justify-center">
    <div
      class="mx-auto flex w-full max-w-sm flex-col gap-6 border-3 border-brut-black bg-brut-white p-6 shadow-brut"
    >
      <h1>{{ t('login.title') }}</h1>
      <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
        <BaseInput
          v-model="form.email"
          type="email"
          :label="t('login.emailLabel')"
          required
          autocomplete="email"
        />
        <BaseInput
          v-model="form.password"
          type="password"
          :label="t('login.passwordLabel')"
          required
          autocomplete="current-password"
        />
        <p v-if="errorMessage" role="alert" aria-live="assertive" class="font-bold text-brut-red">
          {{ errorMessage }}
        </p>
        <BaseButton type="submit" :disabled="isSubmitting">{{ t('login.submit') }}</BaseButton>
      </form>

      <div
        class="flex items-center gap-2 before:content-[''] before:flex-1 before:border-b-3 before:border-brut-black after:content-[''] after:flex-1 after:border-b-3 after:border-brut-black"
      >
        <span class="text-sm text-brut-gray">{{ t('login.or') }}</span>
      </div>

      <BaseButton v-if="ssoAvailable" type="button" variant="secondary" @click="loginWithSSO">
        {{ t('login.ssoButton') }}
      </BaseButton>
      <p v-else class="text-sm text-center text-brut-gray">
        {{ t('login.ssoUnavailable') }}
      </p>
    </div>
  </div>
</template>
