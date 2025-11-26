package port

type JWTService interface {
	GenerateAccessToken(userID string, role string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string) (string, string, error) // Returns userID and role
	ValidateRefreshToken(token string) (string, error)   // Returns userID
}



