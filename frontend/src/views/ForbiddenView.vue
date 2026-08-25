<script setup lang="ts">
import { useI18n } from 'vue-i18n'

/**
 * Shown when a valid session lacks the role required for a route or an
 * API call responded 403. The session is left intact.
 */
const { t } = useI18n()
</script>

<template>
  <div
    class="mx-auto flex max-w-md flex-col items-start gap-4 border-3 border-brut-black bg-brut-white p-6 shadow-brut"
  >
    <h1>{{ t('forbidden.title') }}</h1>
    <p>{{ t('forbidden.message') }}</p>
    <RouterLink to="/items" class="font-bold">{{ t('forbidden.backToItems') }}</RouterLink>
  </div>
</template>
