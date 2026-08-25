import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as itemsApi from '@/api/items'
import * as loanRequestsApi from '@/api/loanRequests'
import { useAuthStore } from '@/stores/auth'
import LoanRequestFormView from './LoanRequestFormView.vue'

vi.mock('@/api/items')
vi.mock('@/api/loanRequests')

const item = {
  inventoryNumber: '000001',
  name: 'Drill',
  description: '',
  imagePath: null,
  category: 'Tools',
  location: 'A1',
  amount: 1,
  electricalAppliance: true,
  note: '',
}

function makeRouter() {
  return createRouter({
    history: createWebHistory(),
    routes: [
      { path: '/loan-requests', component: { template: '<div />' } },
      { path: '/loan-requests/new', component: LoanRequestFormView },
      { path: '/loan-requests/:id', component: { template: '<div />' } },
    ],
  })
}

describe('LoanRequestFormView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    vi.mocked(itemsApi.listItems).mockResolvedValue([item])
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
  })

  it('creates a loan request and redirects to its detail page', async () => {
    vi.mocked(loanRequestsApi.createLoanRequest).mockResolvedValue({
      id: '1',
      items: [item],
      requestingUserId: 'u1',
      dateOfLending: '2026-09-01',
      dateOfReturn: '2026-09-10',
      locationOfItems: 'Room 1',
      status: 'pending',
    })
    const router = makeRouter()
    await router.push('/loan-requests/new')
    await router.isReady()
    const wrapper = mount(LoanRequestFormView, { global: { plugins: [i18n, router] } })
    await flushPromises()

    await wrapper.get('input[type="checkbox"]').setValue(true)
    const dateInputs = wrapper.findAll('input[type="date"]')
    await dateInputs[0].setValue('2026-09-01')
    await dateInputs[1].setValue('2026-09-10')
    await wrapper.get('input[type="text"]').setValue('Room 1')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(loanRequestsApi.createLoanRequest).toHaveBeenCalledWith({
      itemInventoryNumbers: ['000001'],
      dateOfLending: '2026-09-01',
      dateOfReturn: '2026-09-10',
      locationOfItems: 'Room 1',
    })
    expect(router.currentRoute.value.path).toBe('/loan-requests/1')
  })

  it('has no detectable accessibility violations', async () => {
    const router = makeRouter()
    await router.push('/loan-requests/new')
    await router.isReady()
    const wrapper = mount(LoanRequestFormView, { global: { plugins: [i18n, router] } })
    await flushPromises()
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
