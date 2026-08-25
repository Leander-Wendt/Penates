import { describe, expect, it, vi } from 'vitest'

import client from './client'
import {
  listItems,
  getItem,
  createItem,
  updateItem,
  deleteItem,
  uploadItemImage,
  deleteItemImage,
} from './items'

vi.mock('./client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

describe('items api module', () => {
  it('lists items with optional search and category params', async () => {
    const data = [{ inventoryNumber: '000001', name: 'Drill' }]
    ;(client.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await listItems({ search: 'drill', category: 'tools' })

    expect(client.get).toHaveBeenCalledWith('/items', {
      params: { search: 'drill', category: 'tools' },
    })
    expect(result).toEqual(data)
  })

  it('lists items with no params', async () => {
    ;(client.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data: [] })

    await listItems()

    expect(client.get).toHaveBeenCalledWith('/items', { params: {} })
  })

  it('gets a single item by inventory number', async () => {
    const data = { inventoryNumber: '000001', name: 'Drill' }
    ;(client.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await getItem('000001')

    expect(client.get).toHaveBeenCalledWith('/items/000001')
    expect(result).toEqual(data)
  })

  it('creates an item', async () => {
    const input = {
      name: 'Drill',
      description: '',
      category: 'tools',
      location: 'A1',
      amount: 1,
      electricalAppliance: true,
      note: '',
    }
    const data = { inventoryNumber: '000001', ...input }
    ;(client.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await createItem(input as never)

    expect(client.post).toHaveBeenCalledWith('/items', input)
    expect(result).toEqual(data)
  })

  it('updates an item', async () => {
    const input = { name: 'Drill v2' }
    const data = { inventoryNumber: '000001', name: 'Drill v2' }
    ;(client.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await updateItem('000001', input as never)

    expect(client.put).toHaveBeenCalledWith('/items/000001', input)
    expect(result).toEqual(data)
  })

  it('deletes an item', async () => {
    ;(client.delete as ReturnType<typeof vi.fn>).mockResolvedValue({})

    await deleteItem('000001')

    expect(client.delete).toHaveBeenCalledWith('/items/000001')
  })

  it('uploads an item image as multipart/form-data', async () => {
    const data = { inventoryNumber: '000001', imagePath: '/uploads/000001.png' }
    ;(client.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data })
    const file = new File(['x'], 'photo.png', { type: 'image/png' })

    const result = await uploadItemImage('000001', file)

    expect(client.post).toHaveBeenCalledWith('/items/000001/image', expect.any(FormData), {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    expect(result).toEqual(data)
  })

  it('deletes an item image', async () => {
    ;(client.delete as ReturnType<typeof vi.fn>).mockResolvedValue({})

    await deleteItemImage('000001')

    expect(client.delete).toHaveBeenCalledWith('/items/000001/image')
  })
})
