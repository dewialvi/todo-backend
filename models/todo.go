package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Todo struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	UserID    uint      `json:"user_id"`          // ID user pemilik
	Deadline  time.Time `json:"deadline"`         // Waktu deadline
	FilePath  string	`json:"file_path"`  // path file lampiran
}


// InitDB untuk koneksi ke database PostgreSQL
func InitDB() {
	// Ganti sesuai konfigurasi PostgreSQL kamu
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "200603"        // ubah sesuai password PostgreSQL kamu
	dbname := "todo_db"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal konek ke database PostgreSQL: %v", err)
	}

	// Auto migrate tabel User dan Todo
err = DB.AutoMigrate(&User{}, &Todo{}, &Activity{})

if err != nil {
    fmt.Println("⚠️  Migrasi menghasilkan peringatan, tapi tetap dilanjutkan:", err)
} else {
    fmt.Println("✅ Migrasi database sukses.")
}
}