import { api } from './client'

export interface AuthTokens {
  access_token: string
  refresh_token: string
  expires_in: number
}

export interface LoginResponse {
  data: AuthTokens
}

export interface User {
  id: string
  email: string
  username: string
}

export async function login(email: string, password: string): Promise<AuthTokens> {
  const res = await api.post<LoginResponse>('/auth/login', { email, password })
  return res.data
}

export async function register(email: string, username: string, password: string): Promise<unknown> {
  return api.post('/auth/register', { email, username, password })
}

export async function getProfile(): Promise<User> {
  return api.get<User>('/users/me')
}
