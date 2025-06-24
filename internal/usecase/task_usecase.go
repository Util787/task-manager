package usecase

import (
	"time"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/google/uuid"
)

type TaskUsecase struct {
	taskRepo TaskRepository
}

type TaskRepository interface {
	CreateTask(task *domain.Task)
	GetTaskStateByID(id uuid.UUID) (domain.TaskStatus, time.Duration, error)
	GetTaskResultByID(id uuid.UUID) (string, error)
	DeleteTask(id uuid.UUID) error
}

func NewTaskUsecase(taskRepo TaskRepository) *TaskUsecase {
	return &TaskUsecase{taskRepo: taskRepo}
}

func (t *TaskUsecase) CreateTask(task *domain.Task) error {
	t.taskRepo.CreateTask(task)
	return nil
}

func (t *TaskUsecase) GetTaskStateByID(id uuid.UUID) (domain.TaskStatus, time.Duration, error) {
	return t.taskRepo.GetTaskStateByID(id)
}

func (t *TaskUsecase) GetTaskResultByID(id uuid.UUID) (string, error) {
	return t.taskRepo.GetTaskResultByID(id)
}

func (t *TaskUsecase) DeleteTask(id uuid.UUID) error {
	return t.taskRepo.DeleteTask(id)
}