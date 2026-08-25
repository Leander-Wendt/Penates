import { beforeEach, describe, expect, it, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

import * as itemsApi from '@/api/items'

vi.mock('@/api/items')

const baseItem = {
  inventoryNumber: '000001',
  name: 'Drill',
  description: '',
  imagePath: null,
  category: 'tools',
  location: 'A1',
  amount: 1,
  electricalAppliance: true,
  note: '',
}

describe('items store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('fetches items using the current search and category filters', async () => {
    vi.mocked(itemsApi.listItems).mockResolvedValue([baseItem])

    const { useItemsStore } = await import('./items')
    const store = useItemsStore()
    store.searchTerm = 'drill'
    store.categoryFilter = 'tools'
    await store.fetchAll()

    expect(itemsApi.listItems).toHaveBeenCalledWith({ search: 'drill', category: 'tools' })
    expect(store.items).toEqual([baseItem])
  })

  it('creates an item and appends it to the list', async () => {
    vi.mocked(itemsApi.createItem).mockResolvedValue(baseItem)

    const { useItemsStore } = await import('./items')
    const store = useItemsStore()
    await store.create({
      name: 'Drill',
      description: '',
      category: 'tools',
      location: 'A1',
      amount: 1,
      electricalAppliance: true,
      note: '',
    })

    expect(store.items).toContainEqual(baseItem)
  })

  it('updates an item in place', async () => {
    const updated = { ...baseItem, name: 'Drill v2' }
    vi.mocked(itemsApi.listItems).mockResolvedValue([baseItem])
    vi.mocked(itemsApi.updateItem).mockResolvedValue(updated)

    const { useItemsStore } = await import('./items')
    const store = useItemsStore()
    await store.fetchAll()
    await store.update('000001', { name: 'Drill v2' })

    expect(store.items).toContainEqual(updated)
  })

  it('removes an item from the list on delete', async () => {
    vi.mocked(itemsApi.listItems).mockResolvedValue([baseItem])
    vi.mocked(itemsApi.deleteItem).mockResolvedValue(undefined)

    const { useItemsStore } = await import('./items')
    const store = useItemsStore()
    await store.fetchAll()
    await store.remove('000001')

    expect(store.items).toHaveLength(0)
  })
})
