package dto

type CreateContactInput struct {
	Name    string `json:"name" validate:"required,min=2"`
	Email   string `json:"email" validate:"required,email"`
	Subject string `json:"subject" validate:"max=500"`
	Message string `json:"message" validate:"required,min=10"`
}

type ContactOutput struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
