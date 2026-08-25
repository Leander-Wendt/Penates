import client from './client'
import type { Item, ItemInput, ItemSearchParams } from '@/types'

/**
 * Fetches items, optionally filtered by free-text search and category.
 */
export async function listItems(params: ItemSearchParams = {}): Promise<Item[]> {
  const { data } = await client.get<Item[]>('/items', { params })
  return data
}

/**
 * Fetches a single item by inventory number.
 */
export async function getItem(inventoryNumber: string): Promise<Item> {
  const { data } = await client.get<Item>(`/items/${inventoryNumber}`)
  return data
}

/**
 * Creates a new item.
 */
export async function createItem(input: ItemInput): Promise<Item> {
  const { data } = await client.post<Item>('/items', input)
  return data
}

/**
 * Updates an existing item.
 */
export async function updateItem(
  inventoryNumber: string,
  input: Partial<ItemInput>,
): Promise<Item> {
  const { data } = await client.put<Item>(`/items/${inventoryNumber}`, input)
  return data
}

/**
 * Deletes an item.
 */
export async function deleteItem(inventoryNumber: string): Promise<void> {
  await client.delete(`/items/${inventoryNumber}`)
}

/**
 * Uploads a replacement image for an item as multipart/form-data.
 */
export async function uploadItemImage(inventoryNumber: string, file: File): Promise<Item> {
  const formData = new FormData()
  formData.append('image', file)
  const { data } = await client.post<Item>(`/items/${inventoryNumber}/image`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return data
}

/**
 * Removes an item's image.
 */
export async function deleteItemImage(inventoryNumber: string): Promise<void> {
  await client.delete(`/items/${inventoryNumber}/image`)
}
