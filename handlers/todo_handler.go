package handlers

import (
	"net/http"
	"strconv"
	"todo-backend/models"

	"github.com/labstack/echo/v4"
)

// ---------------------------
// 1️⃣ GET /todos
// ---------------------------
func GetTodos(c echo.Context) error {
	var todos []models.Todo
	models.DB.Find(&todos)
	return c.JSON(http.StatusOK, todos)
}

// ---------------------------
// 2️⃣ POST /todos
// ---------------------------
func CreateTodo(c echo.Context) error {
	var newTodo models.Todo

	if err := c.Bind(&newTodo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format data tidak valid"})
	}

	if newTodo.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title tidak boleh kosong"})
	}

	newTodo.Completed = false
	models.DB.Create(&newTodo)

	return c.JSON(http.StatusCreated, newTodo)
}

// ---------------------------
// 3️⃣ PUT /todos/:id
// ---------------------------
func UpdateTodoStatus(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
	}

	var todo models.Todo
	if err := models.DB.First(&todo, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Todo tidak ditemukan"})
	}

	var updateReq struct {
		Completed bool `json:"completed"`
	}
	if err := c.Bind(&updateReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format data tidak valid"})
	}

	todo.Completed = updateReq.Completed
	models.DB.Save(&todo)

	return c.JSON(http.StatusOK, todo)
}

// ---------------------------
// 4️⃣ DELETE /todos/:id
// ---------------------------
func DeleteTodo(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID tidak valid"})
	}

	var todo models.Todo
	if err := models.DB.First(&todo, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Todo tidak ditemukan"})
	}

	models.DB.Delete(&todo)
	return c.JSON(http.StatusOK, map[string]string{"message": "Todo berhasil dihapus"})
}
