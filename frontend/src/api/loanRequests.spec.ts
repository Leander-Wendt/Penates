import { describe, expect, it, vi } from 'vitest'

import client from './client'
import {
  listLoanRequests,
  getLoanRequest,
  createLoanRequest,
  updateLoanRequestStatus,
} from './loanRequests'

vi.mock('./client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
  },
}))

describe('loanRequests api module', () => {
  it('lists loan requests', async () => {
    const data = [{ id: '1', status: 'pending' }]
    ;(client.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await listLoanRequests()

    expect(client.get).toHaveBeenCalledWith('/loan-requests')
    expect(result).toEqual(data)
  })

  it('gets a single loan request by id', async () => {
    const data = { id: '1', status: 'pending' }
    ;(client.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await getLoanRequest('1')

    expect(client.get).toHaveBeenCalledWith('/loan-requests/1')
    expect(result).toEqual(data)
  })

  it('creates a loan request', async () => {
    const input = {
      itemInventoryNumbers: ['000001'],
      dateOfLending: '2026-09-01',
      dateOfReturn: '2026-09-10',
      locationOfItems: 'Room 1',
    }
    const data = { id: '1', ...input, status: 'pending' }
    ;(client.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await createLoanRequest(input)

    expect(client.post).toHaveBeenCalledWith('/loan-requests', input)
    expect(result).toEqual(data)
  })

  it('patches a loan request status', async () => {
    const data = { id: '1', status: 'approved' }
    ;(client.patch as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await updateLoanRequestStatus('1', 'approved')

    expect(client.patch).toHaveBeenCalledWith('/loan-requests/1/status', { status: 'approved' })
    expect(result).toEqual(data)
  })
})
