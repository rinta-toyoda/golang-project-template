package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect は DSN から GORM の DB を返します
func Connect(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
