import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as usersApi from '@/api/users'
import * as organisationsApi from '@/api/organisations'
import UserFormView from './UserFormView.vue'

vi.mock('@/api/users')
vi.mock('@/api/organisations')

const organisations = [{ id: 'o1', name: 'Org 1', description: '' }]
const existingUser = { id: '1', email: 'a@b.com', role: 'student' as const, organisationId: 'o1' }

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/users', component: { template: '<div />' } },
      { path: '/users/new', component: UserFormView },
      { path: '/users/:id/edit', component: UserFormView },
    ],
  })
}

describe('UserFormView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    vi.mocked(organisationsApi.listOrganisations).mockResolvedValue(organisations)
  })

  it('creates a user and redirects to the list on submit', async () => {
    vi.mocked(usersApi.createUser).mockResolvedValue({
      id: '2',
      email: 'new@b.com',
      role: 'student',
      organisationId: 'o1',
    })
    const router = makeRouter()
    await router.push('/users/new')
    await router.isReady()
    const wrapper = mount(UserFormView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    await wrapper.get('input[type="email"]').setValue('new@b.com')
    await wrapper.get('input[type="password"]').setValue('secret123')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(usersApi.createUser).toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/users')
  })

  it('loads and pre-fills an existing user in edit mode', async () => {
    vi.mocked(usersApi.listUsers).mockResolvedValue([existingUser])
    const router = makeRouter()
    await router.push('/users/1/edit')
    await router.isReady()
    const wrapper = mount(UserFormView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    expect(wrapper.text()).toContain('Edit user')
    const emailInput = wrapper.get('input[type="email"]').element as HTMLInputElement
    expect(emailInput.value).toBe('a@b.com')
  })

  it('has no detectable accessibility violations', async () => {
    const router = makeRouter()
    await router.push('/users/new')
    await router.isReady()
    const wrapper = mount(UserFormView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
