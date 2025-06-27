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

type GetTaskStateResponse struct {
	State     domain.TaskState `json:"state"`
	CreatedAt time.Time        `json:"created_at"`
}

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

	c.JSON(http.StatusOK, GetTaskStateResponse{
		State:     state,
		CreatedAt: createdAt,
	})
}
