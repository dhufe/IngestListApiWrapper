package models

import (
	"time"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "Pending"
	StatusRunning   TaskStatus = "Running"
	StatusCompleted TaskStatus = "Completed"
	StatusFailed    TaskStatus = "Failed"
)

type Task struct {
	ID          uint       `gorm:"primaryKey"`
	Title       string     `gorm:"not null"`
	Command     string     `gorm:"not null"`  // Der auszuführende Befehl
	Arguments   string     `gorm:"type:text"` // Argumente für den Befehl
	Status      TaskStatus `gorm:"default:pending"`
	DueDate     *time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	Output      string    `gorm:"type:text"` // Ausgabe des Programms
	Error       string    `gorm:"type:text"` // Fehlerausgabe
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type TaskOutput struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    uint      `gorm:"index"`
	Output    string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
