package handlers

import (
	"net/http"
	"time"
	"todo-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Register user baru
func Register(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	user := models.User{Username: req.Username}
	if err := user.HashPassword(req.Password); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to hash password"})
	}

	if err := models.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "username already exists"})
	}

	// Log aktivitas
	models.CreateActivityLog(user.ID, "register_user")

	return c.JSON(http.StatusOK, map[string]string{"message": "user registered"})
}

// Login user
func Login(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid username or password"})
	}

	if !user.CheckPassword(req.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid username or password"})
	}

	claims := &JwtCustomClaims{
	ID:    user.ID,
	Name:  user.Username,
	Admin: user.Username == "dewi", // set true kalau admin
	RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
	},
}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("rahasia_superkuat"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to generate token"})
	}

	// Log aktivitas login
	models.CreateActivityLog(user.ID, "login_user")

	return c.JSON(http.StatusOK, map[string]string{"token": t})
}
