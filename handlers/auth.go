package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtKey = []byte("rahasia_superkuat") // ubah sesuai keinginanmu

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTResponse struct {
	Token string `json:"token"`
}

// Login handler
func Login(c echo.Context) error {
	var creds Credentials
	if err := c.Bind(&creds); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	// Ganti ini sesuai sistem user kamu (sementara hardcode)
	if creds.Username != "dewi" || creds.Password != "12345" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	// Buat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // token berlaku 1 jam
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not generate token"})
	}

	return c.JSON(http.StatusOK, JWTResponse{Token: tokenString})
}
