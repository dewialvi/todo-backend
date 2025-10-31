package main

import (
	"net/http"
	"os"
	"time"

	"todo-backend/models"
	"todo-backend/routes"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Struktur claims JWT
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

// Handler login untuk menghasilkan token JWT
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Autentikasi sederhana (sementara hardcoded)
	if username != "dewi" || password != "1234" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "invalid username or password",
		})
	}

	// Membuat claims token
	claims := &jwtCustomClaims{
		Name:  username,
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // token berlaku 3 hari
		},
	}

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to generate token",
		})
	}

	// Return token ke client
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func main() {
	// Inisialisasi database
	models.InitDB()

	e := echo.New()

	// Middleware utama
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Endpoint login — tanpa autentikasi
	e.POST("/login", login)

	// Grup endpoint /todos yang dilindungi JWT
	api := e.Group("/todos")

	// Konfigurasi JWT middleware
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}

	// Pasang middleware JWT
	api.Use(echojwt.WithConfig(config))

	// Inisialisasi rute CRUD todo
	routes.InitRoutes(api)

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

	e.Logger.Fatal(e.Start(":" + port))
}
