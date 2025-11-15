package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"todo-backend/handlers"
	"todo-backend/models"
	"todo-backend/routes"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Struktur claims JWT lokal (jika dibutuhkan login manual)
type jwtCustomClaims struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

// Handler login sederhana (opsional)
func login(c echo.Context) error {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&creds); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request"})
	}

	if creds.Username != "dewi" || creds.Password != "12345" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid username or password"})
	}

	claims := &jwtCustomClaims{
		Name:  creds.Username,
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("rahasia_superkuat"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t})
}

func main() {
	// Inisialisasi database
	models.InitDB()

	e := echo.New()

	// Konfigurasi JWT global
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JwtCustomClaims)
		},
		SigningKey: []byte("rahasia_superkuat"),
	}

	// Middleware utama
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Endpoint publik (tanpa JWT)
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)
	e.POST("/upload", handlers.UploadFile)
	e.Static("/uploads", "./uploads") // akses file di http://localhost:8080/uploads/nama_file

	// Grup endpoint ToDo (perlu JWT)
	api := e.Group("/todos")
	api.Use(echojwt.WithConfig(config))
	routes.InitRoutes(api)

	// Grup endpoint admin (perlu JWT)
	admin := e.Group("/admin")
	admin.Use(echojwt.WithConfig(config))
	routes.InitAdminRoutes(admin)

	// Endpoint profil user dari token JWT
	e.GET("/me", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*handlers.JwtCustomClaims)
		return c.JSON(http.StatusOK, echo.Map{
			"id":       claims.ID,
			"username": claims.Name,
		})
	}, echojwt.WithConfig(config))

	// Endpoint reminder — tampilkan todo yang akan jatuh tempo 1 jam lagi
	e.GET("/reminder", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*handlers.JwtCustomClaims)
		userID := claims.ID

		now := time.Now()
		limit := now.Add(1 * time.Hour)
		var todos []models.Todo
		if err := models.DB.Where(
			"user_id = ? AND completed = ? AND deadline BETWEEN ? AND ?",
			userID, false, now, limit,
		).Find(&todos).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Gagal mengambil data reminder",
			})
		}
		return c.JSON(http.StatusOK, todos)
	}, echojwt.WithConfig(config))

	// Reminder otomatis di background setiap 1 menit
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			var todos []models.Todo
			now := time.Now()
			reminderThreshold := now.Add(10 * time.Minute)
			models.DB.Where("deadline <= ? AND completed = ?", reminderThreshold, false).Find(&todos)
			for _, t := range todos {
				fmt.Printf("[Reminder] Todo '%s' akan lewat deadline pada %s!\n",
					t.Title, t.Deadline.Format("15:04"))
			}
		}
	}()

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
}
