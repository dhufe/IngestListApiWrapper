package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewDatabase(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	var db *gorm.DB
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		return db, err
	}

	return db, nil
}
