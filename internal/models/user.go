package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramID   int64 `gorm:"uniqueIndex"`
	Username     string
	FirstName    string
	LastName     string
	IsSubscribed bool      `gorm:"default:false"`
	LastChecked  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Course struct {
	gorm.Model
	Title       string
	Description string
	Price       float64
}

type Access struct {
	gorm.Model
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	CourseID  uint
	Course    Course `gorm:"foreignKey:CourseID"`
	ExpiresAt time.Time
}
