package model

import "time"

type User struct {
	ID           string     `gorm:"primaryKey;type:char(36)" json:"id"`
	Email        string     `gorm:"size:50;not null;unique"   json:"email"`
	Phone        string     `gorm:"size:15;not null;unique"   json:"phone"`
	PasswordHash string     `gorm:"column:password_hash;size:100;not null" json:"-"`
	IsDeleted    bool       `gorm:"default:false"             json:"is_deleted"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"            json:"created_at"`
	DeletedAt    *time.Time `gorm:"index"                     json:"deleted_at,omitempty"`
}
