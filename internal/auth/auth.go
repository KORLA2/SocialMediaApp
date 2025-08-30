package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateToken(jwt.Claims) (string, error)
	ValidateToken(string) (*jwt.Token, error)
}
