<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import NavBar from './NavBar.vue'
import Footer from './Footer.vue'

/**
 * Application-wide shell: a skip-to-content link for keyboard/screen
 * reader users, the primary navigation, the routed page content, and the footer.
 */
const { t } = useI18n()
</script>

<template>
  <a href="#main-content" class="skip-link">{{ t('app.skipToContent') }}</a>
  <NavBar />
  <main id="main-content" tabindex="-1" class="mx-auto max-w-5xl px-4 py-6">
    <RouterView />
  </main>
  <Footer />
</template>
