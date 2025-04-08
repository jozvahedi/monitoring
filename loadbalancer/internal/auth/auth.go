package auth

type Authenticator interface {
	Authenticate(username, password string) bool
}

type BasicAuthService struct{}

func NewBasicAuthService() *BasicAuthService {
	return &BasicAuthService{}
}

func (b *BasicAuthService) Authenticate(username, password string) bool {
	// Simple authentication logic (replace with actual logic)
	return username == "admin" && password == "password"
}
