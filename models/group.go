package models

import "time"

type Group struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	OwnerID   uint   `json:"owner_id"` // creator of the group
	Owner     User   `gorm:"foreignKey:OwnerID"`
	CreatedAt time.Time
}

type GroupUser struct {
	ID        uint  `json:"id"`
	UserID    uint  `json:"user_id"`
	User      User  `gorm:"foreignKey:UserID"`
	GroupID   uint  `json:"group_id"`
	Group     Group `gorm:"foreignKey:GroupID"`
	IsAdmin   bool  `json:"is_admin"` // true if this user can create/update/delete
	CreatedAt time.Time
}
