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
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TaskState   TaskState `json:"task_state"`
	Result      string    `json:"result"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskState struct {
	Status       TaskStatus    `json:"status"`
	WorkDuration time.Duration `json:"work_duration"`
}

var (
	ErrTaskNotFound = errors.New("task not found")
)
