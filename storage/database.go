package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func NewDatabase(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/Berlin",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	fmt.Printf("%s\n", dsn)

	var db *gorm.DB
	if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		return db, err
	}

	return db, nil
}
