import { describe, expect, it, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { axe } from 'vitest-axe'

import en from '@/locales/en.json'
import de from '@/locales/de.json'
import LanguageSwitcher from './LanguageSwitcher.vue'

function createTestI18n() {
  return createI18n({
    legacy: false,
    locale: 'en',
    fallbackLocale: 'en',
    messages: { en, de },
  })
}

describe('LanguageSwitcher', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('renders a button for each supported locale', () => {
    const i18n = createTestI18n()
    const wrapper = mount(LanguageSwitcher, { global: { plugins: [i18n] } })
    const buttons = wrapper.findAll('button')
    expect(buttons.map((b) => b.text())).toEqual(['en', 'de'])
  })

  it('marks the active locale as pressed', () => {
    const i18n = createTestI18n()
    const wrapper = mount(LanguageSwitcher, { global: { plugins: [i18n] } })
    const [enButton, deButton] = wrapper.findAll('button')
    expect(enButton.attributes('aria-pressed')).toBe('true')
    expect(deButton.attributes('aria-pressed')).toBe('false')
  })

  it('switches the active locale and persists it to localStorage on click', async () => {
    const i18n = createTestI18n()
    const wrapper = mount(LanguageSwitcher, { global: { plugins: [i18n] } })
    const [, deButton] = wrapper.findAll('button')
    await deButton.trigger('click')

    expect(i18n.global.locale.value).toBe('de')
    expect(localStorage.getItem('penates.locale')).toBe('de')
  })

  it('has no detectable accessibility violations', async () => {
    const i18n = createTestI18n()
    const wrapper = mount(LanguageSwitcher, { global: { plugins: [i18n] } })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
