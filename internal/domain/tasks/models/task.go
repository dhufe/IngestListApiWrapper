package models

import (
	"time"
)

type TaskStatus string

const (
	StatusPending     TaskStatus = "Pending"
	StatusRunning     TaskStatus = "Running"
	StatusProgressing TaskStatus = "Progressing"
	StatusCompleted   TaskStatus = "Completed"
	StatusFailed      TaskStatus = "Failed"
)

type Task struct {
	ID          uint       `gorm:"primaryKey"`
	FileName    string     `gorm:"not null"`
	Status      TaskStatus `gorm:"default:pending"`
	StartedAt   *time.Time
	CompletedAt *time.Time
	Output      string    `gorm:"type:text"` // Ausgabe des Programms
	Error       string    `gorm:"type:text"` // Fehlerausgabe
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
