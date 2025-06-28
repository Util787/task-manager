package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
