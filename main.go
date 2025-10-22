package main

import (
	"os" // ✅ tambahkan ini

	"todo-backend/models"
	"todo-backend/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Inisialisasi database 
	models.InitDB()

	e := echo.New()

	// Middleware (logging, error recovery, CORS)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Inisialisasi routes
	routes.InitRoutes(e)

	// Jalankan server dengan port dari environment (Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default jika tidak ada environment variable
	}

	e.Logger.Fatal(e.Start(":" + port))
}
