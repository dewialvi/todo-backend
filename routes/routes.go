package routes

import (
	"todo-backend/handlers"

	"github.com/labstack/echo/v4"
)

// InitRoutes berfungsi untuk mengatur semua endpoint
func InitRoutes(e *echo.Echo) {
	// Endpoint utama CRUD To-Do
	e.GET("/todos", handlers.GetTodos)           // Menampilkan semua to-do
	e.POST("/todos", handlers.CreateTodo)        // Menambahkan to-do baru
	e.PUT("/todos/:id", handlers.UpdateTodoStatus) // Mengubah status to-do
	e.DELETE("/todos/:id", handlers.DeleteTodo)  // Menghapus to-do
}
