import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import ForbiddenView from './ForbiddenView.vue'

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/forbidden', component: ForbiddenView },
      { path: '/items', component: { template: '<div />' } },
    ],
  })
}

describe('ForbiddenView', () => {
  it('shows the not-permitted message and a link back to items', async () => {
    const router = makeRouter()
    await router.push('/forbidden')
    await router.isReady()
    const wrapper = mount(ForbiddenView, { global: { plugins: [i18n, router] } })

    expect(wrapper.text()).toContain('Not permitted')
    expect(wrapper.get('a').attributes('href')).toBe('/items')
  })

  it('has no detectable accessibility violations', async () => {
    const router = makeRouter()
    await router.push('/forbidden')
    await router.isReady()
    const wrapper = mount(ForbiddenView, { global: { plugins: [i18n, router] } })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
