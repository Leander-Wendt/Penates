import { createI18n } from 'vue-i18n'

import en from './en.json'
import de from './de.json'

export type SupportedLocale = 'en' | 'de'

export const SUPPORTED_LOCALES: SupportedLocale[] = ['en', 'de']
export const LOCALE_STORAGE_KEY = 'penates.locale'

/**
 * Determines the initial locale: a previously persisted choice takes
 * priority, then the browser's language (German if it starts with "de"),
 * falling back to English.
 */
export function detectInitialLocale(): SupportedLocale {
  const stored = localStorage.getItem(LOCALE_STORAGE_KEY)
  if (stored === 'en' || stored === 'de') {
    return stored
  }
  return navigator.language.toLowerCase().startsWith('de') ? 'de' : 'en'
}

const i18n = createI18n({
  legacy: false,
  locale: detectInitialLocale(),
  fallbackLocale: 'en',
  messages: { en, de },
})

export default i18n
