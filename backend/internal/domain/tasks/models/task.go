package models

import (
	"time"
)

type (
	TaskStatus string
	TaskType   string
)

const (
	StatusPending     TaskStatus = "Pending"
	StatusRunning     TaskStatus = "Running"
	StatusProgressing TaskStatus = "Progressing"
	StatusCompleted   TaskStatus = "Completed"
	StatusFailed      TaskStatus = "Failed"

	TypeBagit    TaskType = "Bagit"
	TypeIdentify TaskType = "Identify"
	TypeValidate TaskType = "Validate"
)

type Task struct {
	ID          uint       `gorm:"primaryKey"       json:"id"`
	FileName    string     `gorm:"not null"         json:"filename"`
	Status      TaskStatus `gorm:"default:Pending"  json:"status"`
	Type        TaskType   `gorm:"default:Identify" json:"type"`
	StartedAt   *time.Time `                        json:"started_at"`
	CompletedAt *time.Time `                        json:"completed_at"`
	Output      string     `gorm:"type:text"        json:"output"` // Ausgabe des Programms
	Error       string     `gorm:"type:text"        json:"error"`  // Fehlerausgabe
	CreatedAt   time.Time  `gorm:"autoCreateTime"   json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"   json:"updated_at"`
}
