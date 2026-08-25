import { defineStore } from 'pinia'
import { ref } from 'vue'

import * as organisationsApi from '@/api/organisations'
import type { Organisation, OrganisationInput } from '@/types'

/**
 * Holds the organisation list and CRUD operations against the
 * organisations API.
 */
export const useOrganisationsStore = defineStore('organisations', () => {
  const organisations = ref<Organisation[]>([])
  const isLoading = ref(false)

  /**
   * Loads all organisations from the API into the store.
   */
  async function fetchAll(): Promise<void> {
    isLoading.value = true
    try {
      organisations.value = await organisationsApi.listOrganisations()
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Creates a new organisation and appends it to the store.
   */
  async function create(input: OrganisationInput): Promise<Organisation> {
    const created = await organisationsApi.createOrganisation(input)
    organisations.value.push(created)
    return created
  }

  /**
   * Updates an organisation and replaces it in the store.
   */
  async function update(id: string, input: OrganisationInput): Promise<Organisation> {
    const updated = await organisationsApi.updateOrganisation(id, input)
    const index = organisations.value.findIndex((org) => org.id === id)
    if (index !== -1) {
      organisations.value[index] = updated
    }
    return updated
  }

  /**
   * Deletes an organisation and removes it from the store.
   */
  async function remove(id: string): Promise<void> {
    await organisationsApi.deleteOrganisation(id)
    organisations.value = organisations.value.filter((org) => org.id !== id)
  }

  return { organisations, isLoading, fetchAll, create, update, remove }
})
