import client from './client'
import type { Organisation, OrganisationInput } from '@/types'

/**
 * Fetches all organisations.
 */
export async function listOrganisations(): Promise<Organisation[]> {
  const { data } = await client.get<Organisation[]>('/organisations')
  return data
}

/**
 * Fetches a single organisation by id.
 */
export async function getOrganisation(id: string): Promise<Organisation> {
  const { data } = await client.get<Organisation>(`/organisations/${id}`)
  return data
}

/**
 * Creates a new organisation.
 */
export async function createOrganisation(input: OrganisationInput): Promise<Organisation> {
  const { data } = await client.post<Organisation>('/organisations', input)
  return data
}

/**
 * Updates an existing organisation.
 */
export async function updateOrganisation(
  id: string,
  input: OrganisationInput,
): Promise<Organisation> {
  const { data } = await client.put<Organisation>(`/organisations/${id}`, input)
  return data
}

/**
 * Deletes an organisation.
 */
export async function deleteOrganisation(id: string): Promise<void> {
  await client.delete(`/organisations/${id}`)
}
