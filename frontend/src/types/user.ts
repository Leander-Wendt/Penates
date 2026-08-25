import type { Role } from './auth'

export interface User {
  id: string
  email: string
  role: Role
  organisationId: string
}

export interface UserCreateInput {
  email: string
  password: string
  role: Role
  organisationId: string
}

export type UserUpdateInput = Partial<Omit<UserCreateInput, 'password'>>
