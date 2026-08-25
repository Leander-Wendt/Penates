import { describe, expect, it, vi } from 'vitest'

import client from './client'
import {
  listOrganisations,
  getOrganisation,
  createOrganisation,
  updateOrganisation,
  deleteOrganisation,
} from './organisations'

vi.mock('./client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

describe('organisations api module', () => {
  it('lists organisations', async () => {
    const data = [{ id: '1', name: 'Org', description: 'desc' }]
    ;(client.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await listOrganisations()

    expect(client.get).toHaveBeenCalledWith('/organisations')
    expect(result).toEqual(data)
  })

  it('gets a single organisation by id', async () => {
    const data = { id: '1', name: 'Org', description: 'desc' }
    ;(client.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await getOrganisation('1')

    expect(client.get).toHaveBeenCalledWith('/organisations/1')
    expect(result).toEqual(data)
  })

  it('creates an organisation', async () => {
    const input = { name: 'Org', description: 'desc' }
    const data = { id: '1', ...input }
    ;(client.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await createOrganisation(input)

    expect(client.post).toHaveBeenCalledWith('/organisations', input)
    expect(result).toEqual(data)
  })

  it('updates an organisation', async () => {
    const input = { name: 'New name', description: 'desc' }
    const data = { id: '1', ...input }
    ;(client.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await updateOrganisation('1', input)

    expect(client.put).toHaveBeenCalledWith('/organisations/1', input)
    expect(result).toEqual(data)
  })

  it('deletes an organisation', async () => {
    ;(client.delete as ReturnType<typeof vi.fn>).mockResolvedValue({})

    await deleteOrganisation('1')

    expect(client.delete).toHaveBeenCalledWith('/organisations/1')
  })
})
