import { describe, expect, it, vi } from 'vitest'

import client from './client'
import { login } from './auth'

vi.mock('./client', () => ({
  default: {
    post: vi.fn(),
  },
}))

describe('auth api module', () => {
  it('posts credentials to /auth/login and returns the response body', async () => {
    const responseBody = { token: 'jwt-token', expiresAt: '2026-01-01T00:00:00Z' }
    ;(client.post as ReturnType<typeof vi.fn>).mockResolvedValue({ data: responseBody })

    const result = await login({ email: 'a@b.com', password: 'secret' })

    expect(client.post).toHaveBeenCalledWith('/auth/login', {
      email: 'a@b.com',
      password: 'secret',
    })
    expect(result).toEqual(responseBody)
  })
})
