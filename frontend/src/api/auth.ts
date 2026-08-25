import client from './client'
import type { LoginCredentials, LoginResponse } from '@/types'

/**
 * Authenticates a user against the backend and returns the issued JWT and
 * its expiry timestamp.
 */
export async function login(credentials: LoginCredentials): Promise<LoginResponse> {
  const { data } = await client.post<LoginResponse>('/auth/login', credentials)
  return data
}
