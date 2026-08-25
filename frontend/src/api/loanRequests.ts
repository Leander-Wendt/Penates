import client from './client'
import type { LoanRequest, LoanRequestCreateInput, LoanRequestStatus } from '@/types'

/**
 * Fetches loan requests visible to the current user (the backend
 * role-scopes the result: students only see their own).
 */
export async function listLoanRequests(): Promise<LoanRequest[]> {
  const { data } = await client.get<LoanRequest[]>('/loan-requests')
  return data
}

/**
 * Fetches a single loan request by id.
 */
export async function getLoanRequest(id: string): Promise<LoanRequest> {
  const { data } = await client.get<LoanRequest>(`/loan-requests/${id}`)
  return data
}

/**
 * Creates a new loan request for the current user.
 */
export async function createLoanRequest(input: LoanRequestCreateInput): Promise<LoanRequest> {
  const { data } = await client.post<LoanRequest>('/loan-requests', input)
  return data
}

/**
 * Updates the status of a loan request. Admin/Logistics only.
 */
export async function updateLoanRequestStatus(
  id: string,
  status: LoanRequestStatus,
): Promise<LoanRequest> {
  const { data } = await client.patch<LoanRequest>(`/loan-requests/${id}/status`, { status })
  return data
}
