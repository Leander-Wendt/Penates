import { beforeEach, describe, expect, it, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

import * as usersApi from '@/api/users'

vi.mock('@/api/users')

describe('users store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('fetches and stores the user list', async () => {
    const data = [{ id: '1', email: 'a@b.com', role: 'admin' as const, organisationId: 'o1' }]
    vi.mocked(usersApi.listUsers).mockResolvedValue(data)

    const { useUsersStore } = await import('./users')
    const store = useUsersStore()
    await store.fetchAll()

    expect(store.users).toEqual(data)
  })

  it('creates a user and appends it to the list', async () => {
    const created = { id: '2', email: 'new@b.com', role: 'student' as const, organisationId: 'o1' }
    vi.mocked(usersApi.createUser).mockResolvedValue(created)

    const { useUsersStore } = await import('./users')
    const store = useUsersStore()
    await store.create({
      email: 'new@b.com',
      password: 'secret',
      role: 'student',
      organisationId: 'o1',
    })

    expect(store.users).toContainEqual(created)
  })

  it('updates a user in place', async () => {
    const original = { id: '1', email: 'a@b.com', role: 'admin' as const, organisationId: 'o1' }
    const updated = { id: '1', email: 'a@b.com', role: 'logistics' as const, organisationId: 'o1' }
    vi.mocked(usersApi.listUsers).mockResolvedValue([original])
    vi.mocked(usersApi.updateUser).mockResolvedValue(updated)

    const { useUsersStore } = await import('./users')
    const store = useUsersStore()
    await store.fetchAll()
    await store.update('1', { role: 'logistics' })

    expect(store.users).toContainEqual(updated)
  })

  it('removes a user from the list on delete', async () => {
    const user = { id: '1', email: 'a@b.com', role: 'admin' as const, organisationId: 'o1' }
    vi.mocked(usersApi.listUsers).mockResolvedValue([user])
    vi.mocked(usersApi.deleteUser).mockResolvedValue(undefined)

    const { useUsersStore } = await import('./users')
    const store = useUsersStore()
    await store.fetchAll()
    await store.remove('1')

    expect(store.users).toHaveLength(0)
  })

  it('propagates a 403 error from delete without mutating the list', async () => {
    const user = { id: '1', email: 'a@b.com', role: 'admin' as const, organisationId: 'o1' }
    vi.mocked(usersApi.listUsers).mockResolvedValue([user])
    vi.mocked(usersApi.deleteUser).mockRejectedValue({ response: { status: 403 } })

    const { useUsersStore } = await import('./users')
    const store = useUsersStore()
    await store.fetchAll()
    await expect(store.remove('1')).rejects.toBeDefined()

    expect(store.users).toHaveLength(1)
  })
})
