package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	Id       uint      `gorm:"primary key;autoIncrement" json:"id"`
	FilePath *string   `                                 json:"filePath"`
	Status   *string   `                                 json:"status"`
	Created  time.Time `                                 json:"created"`
}

func MigrateJobs(db *gorm.DB) error {
	err := db.AutoMigrate(&Job{})
	return err
}
