package app

import (
	"log/slog"

	http_adapter "github.com/Util787/task-manager/internal/adapters/http-adapter"
	"github.com/Util787/task-manager/internal/config"
	"github.com/Util787/task-manager/internal/infrastructure/repo/inmemory"
	"github.com/Util787/task-manager/internal/usecase"
)

type App struct {
	HttpAdapter *http_adapter.HttpAdapter
}

func New(cfg config.Config, logger *slog.Logger) *App {
	taskRepo := inmemory.NewTaskRepository(logger)
	taskUsecase := usecase.NewTaskUsecase(taskRepo)
	httpAdapter := http_adapter.New(cfg, logger, taskUsecase)

	return &App{
		HttpAdapter: httpAdapter,
	}
}