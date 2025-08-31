package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           string         `gorm:"primaryKey;type:char(36)" json:"id"`
	UserName     string         `gorm:"size:15;not null;unique" json:"user_name"`
	Email        string         `gorm:"size:50;not null;unique" json:"email"`
	PasswordHash string         `gorm:"size:100;not null" json:"-"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Profile *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}
