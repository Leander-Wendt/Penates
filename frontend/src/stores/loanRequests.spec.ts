import { beforeEach, describe, expect, it, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

import * as loanRequestsApi from '@/api/loanRequests'

vi.mock('@/api/loanRequests')

const baseRequest = {
  id: '1',
  items: [],
  requestingUserId: 'u1',
  dateOfLending: '2026-09-01',
  dateOfReturn: '2026-09-10',
  locationOfItems: 'Room 1',
  status: 'pending' as const,
}

describe('loanRequests store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('fetches and stores the loan request list', async () => {
    vi.mocked(loanRequestsApi.listLoanRequests).mockResolvedValue([baseRequest])

    const { useLoanRequestsStore } = await import('./loanRequests')
    const store = useLoanRequestsStore()
    await store.fetchAll()

    expect(store.loanRequests).toEqual([baseRequest])
  })

  it('creates a loan request and appends it to the list', async () => {
    vi.mocked(loanRequestsApi.createLoanRequest).mockResolvedValue(baseRequest)

    const { useLoanRequestsStore } = await import('./loanRequests')
    const store = useLoanRequestsStore()
    await store.create({
      itemInventoryNumbers: ['000001'],
      dateOfLending: '2026-09-01',
      dateOfReturn: '2026-09-10',
      locationOfItems: 'Room 1',
    })

    expect(store.loanRequests).toContainEqual(baseRequest)
  })

  it('updates the status of a loan request in place', async () => {
    const approved = { ...baseRequest, status: 'approved' as const }
    vi.mocked(loanRequestsApi.listLoanRequests).mockResolvedValue([baseRequest])
    vi.mocked(loanRequestsApi.updateLoanRequestStatus).mockResolvedValue(approved)

    const { useLoanRequestsStore } = await import('./loanRequests')
    const store = useLoanRequestsStore()
    await store.fetchAll()
    await store.updateStatus('1', 'approved')

    expect(store.loanRequests).toContainEqual(approved)
  })
})
