package handlers

import (
	"log/slog"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/google/uuid"
)

type Handlers struct {
	log         *slog.Logger
	taskUsecase TaskUsecase
}

type TaskUsecase interface {
	CreateTask(task *domain.Task) (uuid.UUID, error)
	GetTaskStateByID(id uuid.UUID) (domain.TaskState, error)
	GetTaskResultByID(id uuid.UUID) (string, error)
	DeleteTask(id uuid.UUID) error
}

func New(log *slog.Logger, taskUsecase TaskUsecase) *Handlers {
	return &Handlers{log: log, taskUsecase: taskUsecase}
}
