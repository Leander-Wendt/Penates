import client from './client'
import type { User, UserCreateInput, UserUpdateInput } from '@/types'

/**
 * Fetches all users.
 */
export async function listUsers(): Promise<User[]> {
  const { data } = await client.get<User[]>('/users')
  return data
}

/**
 * Fetches a single user by id.
 */
export async function getUser(id: string): Promise<User> {
  const { data } = await client.get<User>(`/users/${id}`)
  return data
}

/**
 * Creates a new user. The plaintext password is only ever sent here, on
 * creation; it is never received back from the API.
 */
export async function createUser(input: UserCreateInput): Promise<User> {
  const { data } = await client.post<User>('/users', input)
  return data
}

/**
 * Updates an existing user's non-password fields.
 */
export async function updateUser(id: string, input: UserUpdateInput): Promise<User> {
  const { data } = await client.put<User>(`/users/${id}`, input)
  return data
}

/**
 * Deletes a user. Admin only; the backend returns 403 for Logistics.
 */
export async function deleteUser(id: string): Promise<void> {
  await client.delete(`/users/${id}`)
}
