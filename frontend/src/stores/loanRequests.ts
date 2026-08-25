import { defineStore } from 'pinia'
import { ref } from 'vue'

import * as loanRequestsApi from '@/api/loanRequests'
import type { LoanRequest, LoanRequestCreateInput, LoanRequestStatus } from '@/types'

/**
 * Holds the loan request list (already role-scoped by the backend) and
 * operations against the loan requests API.
 */
export const useLoanRequestsStore = defineStore('loanRequests', () => {
  const loanRequests = ref<LoanRequest[]>([])
  const isLoading = ref(false)

  /**
   * Loads loan requests visible to the current user.
   */
  async function fetchAll(): Promise<void> {
    isLoading.value = true
    try {
      loanRequests.value = await loanRequestsApi.listLoanRequests()
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Creates a new loan request and appends it to the store.
   */
  async function create(input: LoanRequestCreateInput): Promise<LoanRequest> {
    const created = await loanRequestsApi.createLoanRequest(input)
    loanRequests.value.push(created)
    return created
  }

  /**
   * Updates a loan request's status and replaces it in the store.
   */
  async function updateStatus(id: string, status: LoanRequestStatus): Promise<LoanRequest> {
    const updated = await loanRequestsApi.updateLoanRequestStatus(id, status)
    const index = loanRequests.value.findIndex((request) => request.id === id)
    if (index !== -1) {
      loanRequests.value[index] = updated
    }
    return updated
  }

  return { loanRequests, isLoading, fetchAll, create, updateStatus }
})
