import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as usersApi from '@/api/users'
import * as organisationsApi from '@/api/organisations'
import { useAuthStore } from '@/stores/auth'
import UserListView from './UserListView.vue'

vi.mock('@/api/users')
vi.mock('@/api/organisations')

const users = [{ id: '1', email: 'a@b.com', role: 'student' as const, organisationId: 'o1' }]
const organisations = [{ id: 'o1', name: 'Org 1', description: '' }]

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/users', component: UserListView },
      { path: '/users/new', component: { template: '<div />' } },
      { path: '/users/:id/edit', component: { template: '<div />' } },
    ],
  })
}

describe('UserListView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    vi.mocked(usersApi.listUsers).mockResolvedValue(users)
    vi.mocked(organisationsApi.listOrganisations).mockResolvedValue(organisations)
  })

  it('renders each user with their resolved organisation name and role label', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    await router.push('/users')
    await router.isReady()
    const wrapper = mount(UserListView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    expect(wrapper.text()).toContain('a@b.com')
    expect(wrapper.text()).toContain('Org 1')
    expect(wrapper.text()).toContain('Student')
  })

  it('shows the delete action for an Admin', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    await router.push('/users')
    await router.isReady()
    const wrapper = mount(UserListView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    expect(wrapper.findAll('button').some((b) => b.text() === 'Delete')).toBe(true)
  })

  it('hides the delete action for Logistics', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'logistics'
    const router = makeRouter()
    await router.push('/users')
    await router.isReady()
    const wrapper = mount(UserListView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    expect(wrapper.findAll('button').some((b) => b.text() === 'Delete')).toBe(false)
  })

  it('navigates to the create route when the create button is clicked', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    await router.push('/users')
    await router.isReady()
    const wrapper = mount(UserListView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    const createButton = wrapper.findAll('button').find((b) => b.text() === 'New user')
    await createButton?.trigger('click')
    await flushPromises()

    expect(router.currentRoute.value.path).toBe('/users/new')
  })

  it('has no detectable accessibility violations', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const router = makeRouter()
    await router.push('/users')
    await router.isReady()
    const wrapper = mount(UserListView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
