import { describe, expect, it, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import AppShell from './AppShell.vue'

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [{ path: '/', component: { template: '<p>Home</p>' } }],
  })
}

describe('AppShell', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders a skip-to-content link targeting the main landmark', async () => {
    const router = makeRouter()
    await router.push('/')
    await router.isReady()
    const wrapper = mount(AppShell, { global: { plugins: [i18n, router] } })

    const skipLink = wrapper.get('a.skip-link')
    expect(skipLink.attributes('href')).toBe('#main-content')
    expect(wrapper.get('main').attributes('id')).toBe('main-content')
  })

  it('renders the routed view inside main', async () => {
    const router = makeRouter()
    await router.push('/')
    await router.isReady()
    const wrapper = mount(AppShell, { global: { plugins: [i18n, router] } })

    expect(wrapper.get('main').text()).toContain('Home')
  })

  it('has no detectable accessibility violations', async () => {
    const router = makeRouter()
    await router.push('/')
    await router.isReady()
    const wrapper = mount(AppShell, { global: { plugins: [i18n, router] } })
    const results = await axe(wrapper.element)
    expect(results).toHaveNoViolations()
  })
})
