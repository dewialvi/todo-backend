package routes

import (
	"todo-backend/handlers"

	"github.com/labstack/echo/v4"
)

// InitRoutes mengatur semua endpoint untuk To-Do
func InitRoutes(g *echo.Group) {
	g.GET("", handlers.GetTodos)             // GET /todos
	g.POST("", handlers.CreateTodo)          // POST /todos 
	g.PUT("/:id", handlers.UpdateTodoStatus) // PUT /todos/:id
	g.DELETE("/:id", handlers.DeleteTodo)    // DELETE /todos/:id
}
