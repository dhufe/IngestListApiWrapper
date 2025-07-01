package models

import (
	"time"

	"gorm.io/gorm"
)

type Jobs struct {
	Id       uint      `gorm:"primary key;autoIncrement" json:"id"`
	FilePath *string   `                                 json:"filePath"`
	Status   *string   `                                 json:"status"`
	Created  time.Time `                                 json:"created"`
	Result   *string   `                                 json:"result"`
}

type Job struct {
	FilePath *string   `                                 json:"filePath"`
	Status   *string   `                                 json:"status"`
	Created  time.Time `                                 json:"created"`
	Result   *string   `                                 json:"result"`
}

func MigrateJobs(db *gorm.DB) error {
	err := db.AutoMigrate(&Jobs{})
	return err
}
