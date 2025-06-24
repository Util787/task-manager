package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusFailed     TaskStatus = "failed"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	Result      string     `json:"result"`
	WorkDuration time.Duration `json:"work_duration"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

var (
	ErrTaskNotFound = errors.New("task not found")
)
