package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
