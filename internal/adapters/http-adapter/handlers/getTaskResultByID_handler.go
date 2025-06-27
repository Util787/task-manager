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

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("task result: %s", result)})
}
