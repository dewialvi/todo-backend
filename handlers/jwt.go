package handlers

import "github.com/golang-jwt/jwt/v5"

// JwtCustomClaims untuk semua JWT di proyek
type JwtCustomClaims struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}
