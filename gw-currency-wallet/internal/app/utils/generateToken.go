package utils

import "github.com/go-chi/jwtauth/v5"

func GenerateToken() *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte("secret"), nil)
}
