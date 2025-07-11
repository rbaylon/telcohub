package models

import "time"

type User struct {
	ID        uint     `gorm:"primaryKey"`
	Username  string   `gorm:"uniqueIndex;not null"`
	Password  string   `gorm:"not null"`     // bcrypt-hashed
	Role      string   `gorm:"default:user"` // "user" or "admin"
	Markers   []Marker `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
}
