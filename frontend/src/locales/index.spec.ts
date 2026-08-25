import { describe, expect, it, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'

import i18n, { LOCALE_STORAGE_KEY, detectInitialLocale } from './index'
import LoginView from '@/views/LoginView.vue'

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [{ path: '/login', component: LoginView }],
  })
}

describe('i18n locale detection', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('falls back to English when the browser language is neither en nor de', () => {
    Object.defineProperty(navigator, 'language', { value: 'fr-FR', configurable: true })
    expect(detectInitialLocale()).toBe('en')
  })

  it('detects German from the browser language', () => {
    Object.defineProperty(navigator, 'language', { value: 'de-DE', configurable: true })
    expect(detectInitialLocale()).toBe('de')
  })

  it('prefers a persisted locale over the browser language', () => {
    localStorage.setItem(LOCALE_STORAGE_KEY, 'de')
    Object.defineProperty(navigator, 'language', { value: 'en-US', configurable: true })
    expect(detectInitialLocale()).toBe('de')
  })
})

describe('switching locale updates rendered text', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    i18n.global.locale.value = 'en'
  })

  it('changes the login view heading and labels when the locale switches to German', async () => {
    const router = makeRouter()
    await router.push('/login')
    await router.isReady()
    const wrapper = mount(LoginView, { global: { plugins: [i18n, router] } })

    expect(wrapper.text()).toContain('Log in')

    i18n.global.locale.value = 'de'
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Anmelden')
    expect(wrapper.text()).not.toContain('Log in')
  })
})
