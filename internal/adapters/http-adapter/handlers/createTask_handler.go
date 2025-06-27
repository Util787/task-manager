package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/gin-gonic/gin"
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
