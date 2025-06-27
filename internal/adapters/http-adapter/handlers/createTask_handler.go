package handlers

import (
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
		newErrorResponse(c, log, http.StatusInternalServerError, "failed to create task", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("task created successfully with id %s", taskID)})
}
