<script setup lang="ts">
import { useI18n } from 'vue-i18n'

/**
 * Catch-all view for routes that do not match any page.
 */
const { t } = useI18n()
</script>

<template>
  <div
    class="mx-auto flex max-w-md flex-col items-start gap-4 border-3 border-brut-black bg-brut-white p-6 shadow-brut"
  >
    <h1>{{ t('notFound.title') }}</h1>
    <p>{{ t('notFound.message') }}</p>
  </div>
</template>
