<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { LOCALE_STORAGE_KEY, SUPPORTED_LOCALES, type SupportedLocale } from '@/locales'

/**
 * Toggles the active vue-i18n locale between the supported languages and
 * persists the choice to localStorage.
 */
const { locale, t } = useI18n()

function selectLocale(next: SupportedLocale): void {
  locale.value = next
  localStorage.setItem(LOCALE_STORAGE_KEY, next)
}
</script>

<template>
  <div role="group" :aria-label="t('nav.language')" class="flex gap-1">
    <button
      v-for="code in SUPPORTED_LOCALES"
      :key="code"
      type="button"
      class="border-3 border-brut-black px-2 py-1 text-sm font-bold uppercase"
      :class="locale === code ? 'bg-brut-yellow text-brut-black' : 'bg-brut-white text-brut-black'"
      :aria-pressed="locale === code"
      @click="selectLocale(code)"
    >
      {{ code }}
    </button>
  </div>
</template>
