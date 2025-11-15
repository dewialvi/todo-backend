package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// UploadFile handler untuk upload file (Echo v4)
func UploadFile(c echo.Context) error {
// Ambil file dari form-data
file, err := c.FormFile("file")
if err != nil {
return c.JSON(http.StatusBadRequest, map[string]string{
"message": "file not found",
})
}

// Pastikan folder uploads ada
uploadDir := "./uploads"
if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "failed to create uploads directory",
		})
	}
}

// Path tujuan file
dst := filepath.Join(uploadDir, file.Filename)

// Buka file yang diupload
src, err := file.Open()
if err != nil {
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"message": "failed to open uploaded file",
	})
}
defer src.Close()

// Buat file tujuan
out, err := os.Create(dst)
if err != nil {
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"message": "failed to create destination file",
	})
}
defer out.Close()

// Salin isi file
if _, err = io.Copy(out, src); err != nil {
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"message": "failed to save file",
	})
}

return c.JSON(http.StatusOK, map[string]string{
	"message": "File uploaded successfully",
	"path":    "/uploads/" + file.Filename,
})


}
