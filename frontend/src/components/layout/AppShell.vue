<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import NavBar from './NavBar.vue'

/**
 * Application-wide shell: a skip-to-content link for keyboard/screen
 * reader users, the primary navigation, and the routed page content.
 */
const { t } = useI18n()
</script>

<template>
  <a href="#main-content" class="skip-link">{{ t('app.skipToContent') }}</a>
  <NavBar />
  <main id="main-content" tabindex="-1" class="mx-auto max-w-5xl px-4 py-6">
    <RouterView />
  </main>
</template>
