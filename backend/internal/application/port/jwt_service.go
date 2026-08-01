package port

type JWTService interface {
	GenerateAccessToken(userID string, role string) (string, error)
	GenerateRefreshToken(userID string, version int) (string, error)
	ValidateToken(token string) (string, string, error) // Returns userID and role
	ValidateRefreshToken(token string) (string, int, error) // Returns userID and token version
}



