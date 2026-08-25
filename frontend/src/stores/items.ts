import { defineStore } from 'pinia'
import { ref } from 'vue'

import * as itemsApi from '@/api/items'
import type { Item, ItemInput } from '@/types'

/**
 * Holds the item list, the active search/category filters, and CRUD
 * operations against the items API.
 */
export const useItemsStore = defineStore('items', () => {
  const items = ref<Item[]>([])
  const isLoading = ref(false)
  const searchTerm = ref('')
  const categoryFilter = ref('')

  /**
   * Loads items from the API using the current search term and category
   * filter.
   */
  async function fetchAll(): Promise<void> {
    isLoading.value = true
    try {
      items.value = await itemsApi.listItems({
        search: searchTerm.value,
        category: categoryFilter.value,
      })
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Creates a new item and appends it to the store.
   */
  async function create(input: ItemInput): Promise<Item> {
    const created = await itemsApi.createItem(input)
    items.value.push(created)
    return created
  }

  /**
   * Updates an item and replaces it in the store.
   */
  async function update(inventoryNumber: string, input: Partial<ItemInput>): Promise<Item> {
    const updated = await itemsApi.updateItem(inventoryNumber, input)
    const index = items.value.findIndex((item) => item.inventoryNumber === inventoryNumber)
    if (index !== -1) {
      items.value[index] = updated
    }
    return updated
  }

  /**
   * Deletes an item and removes it from the store.
   */
  async function remove(inventoryNumber: string): Promise<void> {
    await itemsApi.deleteItem(inventoryNumber)
    items.value = items.value.filter((item) => item.inventoryNumber !== inventoryNumber)
  }

  return { items, isLoading, searchTerm, categoryFilter, fetchAll, create, update, remove }
})
