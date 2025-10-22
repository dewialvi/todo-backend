package models

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Todo struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
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

	// Auto migrate tabel Todo
	if err := DB.AutoMigrate(&Todo{}); err != nil {
		log.Fatalf("❌ Gagal migrasi database: %v", err)
	}

	fmt.Println("✅ Koneksi ke PostgreSQL berhasil & migrasi selesai.")
}
