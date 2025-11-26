package port

type TokenService interface {
	GenerateToken() (string, error)
	GenerateInvitationToken() (string, error)
}



