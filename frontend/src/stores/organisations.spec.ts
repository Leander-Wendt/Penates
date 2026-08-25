import { beforeEach, describe, expect, it, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

import * as organisationsApi from '@/api/organisations'

vi.mock('@/api/organisations')

describe('organisations store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('fetches and stores the organisation list', async () => {
    const data = [{ id: '1', name: 'Org', description: 'desc' }]
    vi.mocked(organisationsApi.listOrganisations).mockResolvedValue(data)

    const { useOrganisationsStore } = await import('./organisations')
    const store = useOrganisationsStore()
    await store.fetchAll()

    expect(store.organisations).toEqual(data)
    expect(store.isLoading).toBe(false)
  })

  it('creates an organisation and appends it to the list', async () => {
    const created = { id: '2', name: 'New', description: 'desc' }
    vi.mocked(organisationsApi.createOrganisation).mockResolvedValue(created)

    const { useOrganisationsStore } = await import('./organisations')
    const store = useOrganisationsStore()
    await store.create({ name: 'New', description: 'desc' })

    expect(store.organisations).toContainEqual(created)
  })

  it('updates an organisation in place', async () => {
    const original = { id: '1', name: 'Org', description: 'desc' }
    const updated = { id: '1', name: 'Renamed', description: 'desc' }
    vi.mocked(organisationsApi.listOrganisations).mockResolvedValue([original])
    vi.mocked(organisationsApi.updateOrganisation).mockResolvedValue(updated)

    const { useOrganisationsStore } = await import('./organisations')
    const store = useOrganisationsStore()
    await store.fetchAll()
    await store.update('1', { name: 'Renamed', description: 'desc' })

    expect(store.organisations).toContainEqual(updated)
    expect(store.organisations).not.toContainEqual(original)
  })

  it('removes an organisation from the list on delete', async () => {
    const org = { id: '1', name: 'Org', description: 'desc' }
    vi.mocked(organisationsApi.listOrganisations).mockResolvedValue([org])
    vi.mocked(organisationsApi.deleteOrganisation).mockResolvedValue(undefined)

    const { useOrganisationsStore } = await import('./organisations')
    const store = useOrganisationsStore()
    await store.fetchAll()
    await store.remove('1')

    expect(store.organisations).toHaveLength(0)
  })
})
