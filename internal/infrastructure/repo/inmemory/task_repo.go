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
	tasks map[uuid.UUID]*domain.Task
	mu    sync.RWMutex // in context of this task its better to use rwmutex than sync.map/mutex
}

func NewTaskRepository(log *slog.Logger) *TaskRepository {
	return &TaskRepository{
		tasks: make(map[uuid.UUID]*domain.Task),
	}
}

func (r *TaskRepository) CreateTask(task *domain.Task) uuid.UUID {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	id := uuid.New()
	r.tasks[id] = task
	return id
}

func (r *TaskRepository) GetTaskStateByID(id uuid.UUID) (domain.TaskState, error) {
	const op = "TaskRepository.GetTaskStateByID"
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return domain.TaskState{}, fmt.Errorf("%s: %w", op, domain.ErrTaskNotFound)
	}

	return task.TaskState, nil
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
