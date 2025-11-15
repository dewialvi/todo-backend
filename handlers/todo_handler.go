package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"todo-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// POST /todos
func CreateTodo(c echo.Context) error {
user := c.Get("user").(*jwt.Token)
claims := user.Claims.(*JwtCustomClaims)

userID := claims.ID


title := c.FormValue("title")
deadlineStr := c.FormValue("deadline")

if title == "" || deadlineStr == "" {
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Title dan deadline harus diisi"})
}

deadline, err := time.Parse(time.RFC3339, deadlineStr)
if err != nil {
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format deadline tidak valid. Gunakan RFC3339"})
}

// Upload file
var filePath string
file, err := c.FormFile("file")
if err == nil && file != nil {
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuka file"})
	}
	defer src.Close()

	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filePath = filepath.Join(uploadDir, filepath.Base(file.Filename))
	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal membuat file di server"})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menyimpan file"})
	}
}

newTodo := models.Todo{
	Title:     title,
	Completed: false,
	UserID:    userID,
	Deadline:  deadline,
	FilePath:  filePath,
}

if err := models.DB.Create(&newTodo).Error; err != nil {
	fmt.Println("❌ Error simpan todo:", err)
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menambahkan todo"})
}

return c.JSON(http.StatusOK, newTodo)


}




// GET /todos
func GetTodos(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	userID := claims.ID

	var todos []models.Todo
	if err := models.DB.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch todos"})
	}

	return c.JSON(http.StatusOK, todos)
}

// PUT /todos/:id
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

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	if todo.UserID != claims.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Anda tidak berhak mengubah todo ini"})
	}

	var updateReq struct {
		Completed bool `json:"completed"`
	}
	if err := c.Bind(&updateReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format data tidak valid"})
	}

	todo.Completed = updateReq.Completed
	if err := models.DB.Save(&todo).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal memperbarui todo"})
	}

	return c.JSON(http.StatusOK, todo)
}

// DELETE /todos/:id
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

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	if todo.UserID != claims.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Anda tidak berhak menghapus todo ini"})
	}

	if err := models.DB.Delete(&todo).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal menghapus todo"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Todo berhasil dihapus"})
}