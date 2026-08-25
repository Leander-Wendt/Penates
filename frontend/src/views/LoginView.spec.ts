import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as authApi from '@/api/auth'
import LoginView from './LoginView.vue'

vi.mock('@/api/auth')

function makeRouter() {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/login', component: LoginView },
      { path: '/items', component: { template: '<div>Items</div>' } },
    ],
  })
  return router
}

describe('LoginView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('renders email and password fields and a submit button', async () => {
    const router = makeRouter()
    await router.push('/login')
    await router.isReady()
    const wrapper = mount(LoginView, { global: { plugins: [i18n, router] } })

    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('redirects to /items on successful login', async () => {
    vi.mocked(authApi.login).mockResolvedValue({ token: 'jwt', expiresAt: '2026-01-01T00:00:00Z' })
    const router = makeRouter()
    await router.push('/login')
    await router.isReady()
    const wrapper = mount(LoginView, { global: { plugins: [i18n, router] } })

    await wrapper.get('input[type="email"]').setValue('a@b.com')
    await wrapper.get('input[type="password"]').setValue('secret')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(router.currentRoute.value.path).toBe('/items')
  })

  it('shows an inline error and stays on the page when credentials are rejected', async () => {
    vi.mocked(authApi.login).mockRejectedValue({ response: { status: 401 } })
    const router = makeRouter()
    await router.push('/login')
    await router.isReady()
    const wrapper = mount(LoginView, { global: { plugins: [i18n, router] } })

    await wrapper.get('input[type="email"]').setValue('a@b.com')
    await wrapper.get('input[type="password"]').setValue('wrong')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(wrapper.get('[role="alert"]').text()).toContain('Invalid email or password')
    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('has no detectable accessibility violations', async () => {
    const router = makeRouter()
    await router.push('/login')
    await router.isReady()
    const wrapper = mount(LoginView, { global: { plugins: [i18n, router] } })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
