export type Role = 'admin' | 'logistics' | 'student'

export interface LoginCredentials {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  expiresAt: string
}
