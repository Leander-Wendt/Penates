import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as loanRequestsApi from '@/api/loanRequests'
import { useAuthStore } from '@/stores/auth'
import LoanRequestListView from './LoanRequestListView.vue'

vi.mock('@/api/loanRequests')

const request = {
  id: '1',
  items: [
    {
      inventoryNumber: '000001',
      name: 'Drill',
      description: '',
      imagePath: null,
      category: 'Tools',
      location: 'A1',
      amount: 1,
      electricalAppliance: true,
      note: '',
    },
  ],
  requestingUserId: 'u1',
  dateOfLending: '2026-09-01',
  dateOfReturn: '2026-09-10',
  locationOfItems: 'Room 1',
  status: 'pending' as const,
}

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/loan-requests', component: LoanRequestListView },
      { path: '/loan-requests/new', component: { template: '<div />' } },
      { path: '/loan-requests/:id', component: { template: '<div />' } },
    ],
  })
}

describe('LoanRequestListView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    vi.mocked(loanRequestsApi.listLoanRequests).mockResolvedValue([request])
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
  })

  it('renders a row per loan request with its status conveyed as text', async () => {
    const router = makeRouter()
    await router.push('/loan-requests')
    await router.isReady()
    const wrapper = mount(LoanRequestListView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    expect(wrapper.text()).toContain('Drill')
    expect(wrapper.text()).toContain('Room 1')
    expect(wrapper.text()).toContain('Pending')
  })

  it('navigates to the detail page when viewing a request', async () => {
    const router = makeRouter()
    await router.push('/loan-requests')
    await router.isReady()
    const wrapper = mount(LoanRequestListView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    const viewButton = wrapper.findAll('button').find((b) => b.text() === 'View')
    await viewButton?.trigger('click')
    await flushPromises()

    expect(router.currentRoute.value.path).toBe('/loan-requests/1')
  })

  it('has no detectable accessibility violations', async () => {
    const router = makeRouter()
    await router.push('/loan-requests')
    await router.isReady()
    const wrapper = mount(LoanRequestListView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
