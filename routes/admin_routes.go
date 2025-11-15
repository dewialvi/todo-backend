package routes

import (
	"todo-backend/handlers"

	"github.com/labstack/echo/v4"
)

// InitAdminRoutes untuk endpoint admin (misalnya /logs)
func InitAdminRoutes(g *echo.Group) {
	g.GET("/logs", handlers.GetActivityLogs)
}
