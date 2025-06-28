package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type createTaskResponse struct {
	Message string `json:"message" example:"task created successfully with id 6bcd175e-cba9-4ba6-b6ef-f3ac37864118"`
}

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a new task with the specified title and description
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body createTaskRequest true "Task title and description"
// @Success 201 {object} createTaskResponse "task created successfully with id {task_id}"
// @Failure 400 {object} errorResponse "invalid request body"
// @Failure 500 {object} errorResponse "failed to create task"
// @Router /tasks [post]
func (h *Handlers) createTask(c *gin.Context) {
	op, _ := c.Get("op")
	log := h.log.With(
		slog.Any("op", op),
	)

	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "invalid request body", err)
		return
	}

	taskID, err := h.taskUsecase.CreateTask(&domain.Task{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, domain.ErrTitleEmpty) || errors.Is(err, domain.ErrTitleTooLong) || errors.Is(err, domain.ErrDescriptionTooLong) {
			newErrorResponse(c, log, http.StatusBadRequest, "invalid request body: "+err.Error(), err)
			return
		}
		newErrorResponse(c, log, http.StatusInternalServerError, "failed to create task", err)
		return
	}

	c.JSON(http.StatusCreated, createTaskResponse{
		Message: fmt.Sprintf("task created successfully with id %s", taskID),
	})
}

type deleteTaskResponse struct {
	Message string `json:"message" example:"task deleted successfully"`
}

// DeleteTask godoc
// @Summary Delete task by ID
// @Description Deletes a task with the specified ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} deleteTaskResponse "task deleted successfully"
// @Failure 400 {object} errorResponse "invalid task ID"
// @Failure 404 {object} errorResponse "task not found"
// @Failure 500 {object} errorResponse "failed to delete task"
// @Router /tasks/{id} [delete]
func (h *Handlers) deleteTask(c *gin.Context) {
	op, _ := c.Get("op")
	log := h.log.With(
		slog.Any("op", op),
	)

	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "invalid task id", err)
		return
	}

	err = h.taskUsecase.DeleteTask(uuid)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			newErrorResponse(c, log, http.StatusNotFound, "task not found", err)
			return
		}
		newErrorResponse(c, log, http.StatusInternalServerError, "failed to delete task", err)
		return
	}

	c.JSON(http.StatusOK, deleteTaskResponse{
		Message: "task deleted successfully",
	})
}

type getTaskResultResponse struct {
	Message string `json:"message" example:"task result: completed"`
}

// GetTaskResultByID godoc
// @Summary Get task result by ID
// @Description Returns the result of task execution
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} getTaskResultResponse "task result: {task_result}"
// @Failure 400 {object} errorResponse "invalid task ID"
// @Failure 404 {object} errorResponse "task not found"
// @Failure 500 {object} errorResponse "failed to get task result"
// @Router /tasks/{id}/result [get]
func (h *Handlers) getTaskResultByID(c *gin.Context) {
	op, _ := c.Get("op")
	log := h.log.With(
		slog.Any("op", op),
	)

	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "invalid task id", err)
		return
	}

	result, err := h.taskUsecase.GetTaskResultByID(uuid)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			newErrorResponse(c, log, http.StatusNotFound, "task not found", err)
			return
		}
		newErrorResponse(c, log, http.StatusInternalServerError, "failed to get task result", err)
		return
	}

	c.JSON(http.StatusOK, getTaskResultResponse{
		Message: fmt.Sprintf("task result: %s", result),
	})
}

type getTaskStateResponse struct {
	State     domain.TaskState `json:"state"`
	CreatedAt time.Time        `json:"created_at" example:"2025-06-28T01:31:19.1864825+03:00"`
}

// GetTaskStateByID godoc
// @Summary Get task state by ID
// @Description Returns the current state (status, work duration in nanoseconds) of the task and its creation time
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID" format(uuid)
// @Success 200 {object} getTaskStateResponse "task state: {state}, created at: {created_at}"
// @Failure 400 {object} errorResponse "invalid task ID"
// @Failure 404 {object} errorResponse "task not found"
// @Failure 500 {object} errorResponse "failed to get task state"
// @Router /tasks/{id}/state [get]
func (h *Handlers) getTaskStateByID(c *gin.Context) {
	op, _ := c.Get("op")
	log := h.log.With(
		slog.Any("op", op),
	)

	id := c.Param("id")

	uuid, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, log, http.StatusBadRequest, "invalid task id", err)
		return
	}

	state, createdAt, err := h.taskUsecase.GetTaskStateByID(uuid)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			newErrorResponse(c, log, http.StatusNotFound, "task not found", err)
			return
		}
		newErrorResponse(c, log, http.StatusInternalServerError, "failed to get task state", err)
		return
	}

	c.JSON(http.StatusOK, getTaskStateResponse{
		State:     state,
		CreatedAt: createdAt,
	})
}
