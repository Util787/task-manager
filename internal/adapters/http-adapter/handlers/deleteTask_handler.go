package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
