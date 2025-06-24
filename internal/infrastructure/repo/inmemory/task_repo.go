package inmemory

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/google/uuid"
)

type TaskRepository struct {
	log   *slog.Logger
	tasks map[uuid.UUID]*domain.Task
	mu    sync.RWMutex // in context of this task its better to use rwmutex than sync.map/mutex
}

func NewTaskRepository(log *slog.Logger) *TaskRepository {
	return &TaskRepository{
		log:   log,
		tasks: make(map[uuid.UUID]*domain.Task),
	}
}

func (r *TaskRepository) CreateTask(task *domain.Task) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	r.tasks[uuid.New()] = task
}

func (r *TaskRepository) GetTaskStateByID(id uuid.UUID) (domain.TaskStatus, time.Duration, error) {
	const op = "TaskRepository.GetTaskStateByID"
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return "", 0, fmt.Errorf("%s: %w", op, domain.ErrTaskNotFound)
	}
	return task.Status, task.WorkDuration, nil
}

func (r *TaskRepository) GetTaskResultByID(id uuid.UUID) (string, error) {
	const op = "TaskRepository.GetTaskResultByID"
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return "", fmt.Errorf("%s: %w", op, domain.ErrTaskNotFound)
	}
	return task.Result, nil
}

func (r *TaskRepository) DeleteTask(id uuid.UUID) error {
	const op = "TaskRepository.DeleteTask"
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return fmt.Errorf("%s: %w", op, domain.ErrTaskNotFound)
	}

	delete(r.tasks, id)
	return nil
}
