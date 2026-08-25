import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { axe } from 'vitest-axe'

import i18n from '@/locales'
import * as loanRequestsApi from '@/api/loanRequests'
import { useAuthStore } from '@/stores/auth'
import LoanRequestDetailView from './LoanRequestDetailView.vue'

vi.mock('@/api/loanRequests')

const pendingRequest = {
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

describe('LoanRequestDetailView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    vi.mocked(loanRequestsApi.listLoanRequests).mockResolvedValue([pendingRequest])
  })

  it('renders the request details and status', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
    const wrapper = mount(LoanRequestDetailView, {
      props: { id: '1' },
      global: { plugins: [i18n] },
    })
    await flushPromises()

    expect(wrapper.text()).toContain('Drill')
    expect(wrapper.text()).toContain('Room 1')
    expect(wrapper.text()).toContain('Pending')
  })

  it('shows approve/reject actions for admin on a pending request', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const wrapper = mount(LoanRequestDetailView, {
      props: { id: '1' },
      global: { plugins: [i18n] },
    })
    await flushPromises()

    expect(wrapper.findAll('button').some((b) => b.text() === 'Approve')).toBe(true)
    expect(wrapper.findAll('button').some((b) => b.text() === 'Reject')).toBe(true)
  })

  it('hides status-change actions for a student', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'student'
    const wrapper = mount(LoanRequestDetailView, {
      props: { id: '1' },
      global: { plugins: [i18n] },
    })
    await flushPromises()

    expect(wrapper.findAll('button')).toHaveLength(0)
  })

  it('calls updateLoanRequestStatus when an action is clicked', async () => {
    vi.mocked(loanRequestsApi.updateLoanRequestStatus).mockResolvedValue({
      ...pendingRequest,
      status: 'approved',
    })
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const wrapper = mount(LoanRequestDetailView, {
      props: { id: '1' },
      global: { plugins: [i18n] },
    })
    await flushPromises()

    const approveButton = wrapper.findAll('button').find((b) => b.text() === 'Approve')
    await approveButton?.trigger('click')
    await flushPromises()

    expect(loanRequestsApi.updateLoanRequestStatus).toHaveBeenCalledWith('1', 'approved')
  })

  it('has no detectable accessibility violations', async () => {
    const authStore = useAuthStore()
    authStore.token = 'jwt'
    authStore.role = 'admin'
    const wrapper = mount(LoanRequestDetailView, {
      props: { id: '1' },
      global: { plugins: [i18n] },
    })
    await flushPromises()
    const results = await axe(wrapper.element, { rules: { region: { enabled: false } } })
    expect(results).toHaveNoViolations()
  })
})
