import { describe, expect, it, vi } from 'vitest'

import client from './client'
import { listUsers, getUser, createUser, updateUser, deleteUser } from './users'

vi.mock('./client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

describe('users api module', () => {
  it('lists users', async () => {
    const data = [{ id: '1', email: 'a@b.com', role: 'admin', organisationId: 'o1' }]
    ;(client.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await listUsers()

    expect(client.get).toHaveBeenCalledWith('/users')
    expect(result).toEqual(data)
  })

  it('gets a single user by id', async () => {
    const data = { id: '1', email: 'a@b.com', role: 'admin', organisationId: 'o1' }
    ;(client.get as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await getUser('1')

    expect(client.get).toHaveBeenCalledWith('/users/1')
    expect(result).toEqual(data)
  })

  it('creates a user including password', async () => {
    const input = { email: 'a@b.com', password: 'secret', role: 'student', organisationId: 'o1' }
    const data = {
      id: '1',
      email: input.email,
      role: input.role,
      organisationId: input.organisationId,
    }
    ;(client.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await createUser(input)

    expect(client.post).toHaveBeenCalledWith('/users', input)
    expect(result).toEqual(data)
  })

  it('updates a user without a password field', async () => {
    const input = { role: 'logistics' as const }
    const data = { id: '1', email: 'a@b.com', role: 'logistics', organisationId: 'o1' }
    ;(client.put as ReturnType<typeof vi.fn>).mockResolvedValue({ data })

    const result = await updateUser('1', input)

    expect(client.put).toHaveBeenCalledWith('/users/1', input)
    expect(result).toEqual(data)
  })

  it('deletes a user', async () => {
    ;(client.delete as ReturnType<typeof vi.fn>).mockResolvedValue({})

    await deleteUser('1')

    expect(client.delete).toHaveBeenCalledWith('/users/1')
  })
})
