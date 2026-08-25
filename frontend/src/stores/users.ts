import { defineStore } from 'pinia'
import { ref } from 'vue'

import * as usersApi from '@/api/users'
import type { User, UserCreateInput, UserUpdateInput } from '@/types'

/**
 * Holds the user list and CRUD operations against the users API.
 */
export const useUsersStore = defineStore('users', () => {
  const users = ref<User[]>([])
  const isLoading = ref(false)

  /**
   * Loads all users from the API into the store.
   */
  async function fetchAll(): Promise<void> {
    isLoading.value = true
    try {
      users.value = await usersApi.listUsers()
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Creates a new user and appends it to the store.
   */
  async function create(input: UserCreateInput): Promise<User> {
    const created = await usersApi.createUser(input)
    users.value.push(created)
    return created
  }

  /**
   * Updates a user and replaces it in the store.
   */
  async function update(id: string, input: UserUpdateInput): Promise<User> {
    const updated = await usersApi.updateUser(id, input)
    const index = users.value.findIndex((user) => user.id === id)
    if (index !== -1) {
      users.value[index] = updated
    }
    return updated
  }

  /**
   * Deletes a user and removes it from the store. Rejects (and leaves the
   * store unchanged) if the API returns 403, which happens when a
   * Logistics user attempts a delete reserved for Admins.
   */
  async function remove(id: string): Promise<void> {
    await usersApi.deleteUser(id)
    users.value = users.value.filter((user) => user.id !== id)
  }

  return { users, isLoading, fetchAll, create, update, remove }
})
