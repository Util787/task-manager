package usecase

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/google/uuid"
)

type TaskUsecase struct {
	taskRepo TaskRepository
}

type TaskRepository interface {
	CreateTask(task *domain.Task) uuid.UUID
	GetTaskStateByID(id uuid.UUID) (domain.TaskState, time.Time, error)
	GetTaskResultByID(id uuid.UUID) (string, error)
	DeleteTask(id uuid.UUID) error
}

func NewTaskUsecase(taskRepo TaskRepository) *TaskUsecase {
	return &TaskUsecase{taskRepo: taskRepo}
}

func (t *TaskUsecase) CreateTask(task *domain.Task) (uuid.UUID, error) {
	const op = "TaskUsecase.CreateTask"

	if err := t.validateTask(task); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	// Initializing task state, in fact the logic would be much more complicated, but I made like this for simplicity
	task.TaskState = domain.TaskState{
		Status:       domain.StatusInProgress,
		WorkDuration: 0,
	}

	id := t.taskRepo.CreateTask(task)
	return id, nil
}

func (t *TaskUsecase) validateTask(task *domain.Task) error {
	if task.Title == "" {
		return domain.ErrTitleEmpty
	}
	if utf8.RuneCountInString(task.Title) > 255 {
		return fmt.Errorf("%w, maximum 255 characters", domain.ErrTitleTooLong)
	}
	if utf8.RuneCountInString(task.Description) > 1000 {
		return fmt.Errorf("%w, maximum 1000 characters", domain.ErrDescriptionTooLong)
	}
	return nil
}

func (t *TaskUsecase) GetTaskStateByID(id uuid.UUID) (domain.TaskState, time.Time, error) {
	const op = "TaskUsecase.GetTaskStateByID"

	state, createdAt, err := t.taskRepo.GetTaskStateByID(id)
	if err != nil {
		return domain.TaskState{}, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}
	return state, createdAt, nil
}

func (t *TaskUsecase) GetTaskResultByID(id uuid.UUID) (string, error) {
	const op = "TaskUsecase.GetTaskResultByID"

	result, err := t.taskRepo.GetTaskResultByID(id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}

func (t *TaskUsecase) DeleteTask(id uuid.UUID) error {
	const op = "TaskUsecase.DeleteTask"

	err := t.taskRepo.DeleteTask(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
