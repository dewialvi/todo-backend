package main

import (
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

	// Jalankan server
	e.Logger.Fatal(e.Start(":8080"))
}
