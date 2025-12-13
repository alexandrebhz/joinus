export interface CreateContactRequest {
  name: string
  email: string
  subject?: string
  message: string
}

export interface ContactResponse {
  id: string
  name: string
  email: string
  subject: string
  message: string
  created_at: string
}
