package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	Profile      *UserProfile   `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	ID           string         `gorm:"primaryKey;type:char(36)" json:"id"`
	UserName     string         `gorm:"size:15;not null;unique" json:"user_name"`
	Email        string         `gorm:"size:50;not null;unique" json:"email"`
	PasswordHash string         `gorm:"size:100;not null" json:"-"`
}

type UserProfile struct {
	UserID string `gorm:"type:char(36);not null" json:"user_id"`
	Name   string `gorm:"size:100" json:"name,omitempty"`
	ID     uint   `gorm:"primaryKey" json:"id"`
}

func (u *User) TableName() string {
	return "users"
}

func (up *UserProfile) TableName() string {
	return "user_profiles"
}
