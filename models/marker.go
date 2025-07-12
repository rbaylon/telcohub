package models

import "time"

type Marker struct {
	ID          uint     `gorm:"primaryKey"`
	Title       string   `gorm:"not null"`
	Description string   `gorm:"type:text"`
	Latitude    float64  `gorm:"not null"`
	Longitude   float64  `gorm:"not null"`
	UserID      uint     `gorm:"not null"`
	User        User     `gorm:"foreignKey:UserID"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	GroupID     uint     `json:"group_id"`
	Group       Group    `gorm:"foreignKey:GroupID"`
	CreatedAt   time.Time
}
