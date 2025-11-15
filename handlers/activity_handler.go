package handlers

import (
	"net/http"
	"time"
	"todo-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GET /logs (admin only)
func GetActivityLogs(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	// Hanya admin yang boleh akses
	if claims.Name != "dewi" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Hanya admin yang boleh mengakses log aktivitas",
		})
	}

	// Struct untuk menampung hasil join
	type ActivityWithUser struct {
		ID       uint      `json:"id"`
		Username string    `json:"username"`
		Action   string    `json:"action"`
		Time     time.Time `json:"time"`
	}

	var result []ActivityWithUser

	// Query join ke tabel users
	err := models.DB.Table("activities").
		Select("activities.id, users.username, activities.action, activities.time").
		Joins("LEFT JOIN users ON users.id = activities.user_id").
		Order("activities.id ASC").
		Scan(&result).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil log aktivitas",
		})
	}

	return c.JSON(http.StatusOK, result)
}
