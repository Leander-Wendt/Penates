<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import { useAuth } from '@/composables/useAuth'

/**
 * Login page: authenticates via email/password and redirects to the item
 * list (or the originally requested route) on success.
 */
const { t } = useI18n()
const { login } = useAuth()
const router = useRouter()
const route = useRoute()

const form = reactive({ email: '', password: '' })
const errorMessage = ref('')
const isSubmitting = ref(false)

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
  <div
    class="mx-auto flex max-w-sm flex-col gap-6 border-3 border-brut-black bg-brut-white p-6 shadow-brut"
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
  </div>
</template>
