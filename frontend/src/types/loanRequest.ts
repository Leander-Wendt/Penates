import type { Item } from './item'

export type LoanRequestStatus = 'pending' | 'approved' | 'rejected' | 'returned' | 'cancelled'

export interface LoanRequest {
  id: string
  items: Item[]
  requestingUserId: string
  dateOfLending: string
  dateOfReturn: string
  locationOfItems: string
  status: LoanRequestStatus
}

export interface LoanRequestCreateInput {
  itemInventoryNumbers: string[]
  dateOfLending: string
  dateOfReturn: string
  locationOfItems: string
}

export interface LoanRequestStatusUpdate {
  status: LoanRequestStatus
}
