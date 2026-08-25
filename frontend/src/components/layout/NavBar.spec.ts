import { describe, expect, it, beforeEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import { useAuthStore } from '@/stores/auth'
import NavBar from './NavBar.vue'

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/', component: { template: '<div />' } },
      { path: '/items', component: { template: '<div />' } },
      { path: '/loan-requests', component: { template: '<div />' } },
      { path: '/users', component: { template: '<div />' } },
      { path: '/organisations', component: { template: '<div />' } },
      { path: '/login', component: { template: '<div />' } },
    ],
  })
}

describe('NavBar', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders nothing when not authenticated', () => {
    const wrapper = mount(NavBar, { global: { plugins: [i18n, makeRouter()] } })
    expect(wrapper.find('nav').exists()).toBe(false)
  })

  it('shows core nav links but hides admin-only links for a student', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
    const router = makeRouter()
    const wrapper = mount(NavBar, { global: { plugins: [i18n, router] } })

    expect(wrapper.text()).toContain('Items')
    expect(wrapper.text()).not.toContain('Users')
    expect(wrapper.text()).not.toContain('Organisations')
  })

  it('shows admin-only links for an admin', () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    const wrapper = mount(NavBar, { global: { plugins: [i18n, router] } })

    expect(wrapper.text()).toContain('Users')
    expect(wrapper.text()).toContain('Organisations')
  })

  it('logs out and redirects to /login when the logout button is clicked', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    await router.push('/items')
    await router.isReady()
    const wrapper = mount(NavBar, { global: { plugins: [i18n, router] } })

    const logoutButton = wrapper.findAll('button').find((b) => b.text() === 'Log out')
    await logoutButton?.trigger('click')
    await flushPromises()

    expect(authStore.isAuthenticated).toBe(false)
    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('has no detectable accessibility violations', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    const wrapper = mount(NavBar, { global: { plugins: [i18n, router] } })
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
