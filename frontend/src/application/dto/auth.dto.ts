import { RegisterData, LoginCredentials, AuthTokens } from '@/domain/value-objects/auth.vo'
import { User } from '@/domain/entities/user.entity'

export interface RegisterRequest extends RegisterData {}
export interface LoginRequest extends LoginCredentials {}
export interface AuthResponse extends AuthTokens {
  user: User
}

