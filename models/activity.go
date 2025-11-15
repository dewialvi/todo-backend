package models

import (
	"time"
)

// Activity mencatat aktivitas user seperti create, update, delete
type Activity struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	Time      time.Time `json:"time"`
}

// CreateActivityLog membuat log baru di tabel activities
func CreateActivityLog(userID uint, action string) {
	activity := Activity{
		UserID: userID,
		Action: action,
		Time:   time.Now(),
	}
	DB.Create(&activity)
}
