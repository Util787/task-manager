package http_adapter

import (
	"context"
	"log/slog"

	"github.com/Util787/task-manager/internal/adapters/http-adapter/handlers"
	"github.com/Util787/task-manager/internal/config"
	"github.com/Util787/task-manager/internal/usecase"
	http_server "github.com/Util787/task-manager/pkg/http-server"
)

type HttpAdapter struct {
	server *http_server.Server
}

func New(cfg config.Config, logger *slog.Logger, svc *usecase.TaskUsecase) *HttpAdapter {
	handler := handlers.New(logger, svc)
	router := handler.InitRoutes(cfg.Env)
	s := http_server.New(cfg.HttpServerCfg, router)

	return &HttpAdapter{
		server: s,
	}
}

func (a HttpAdapter) Start() error {
	if err := a.server.Run(); err != nil {
		return err
	}
	return nil
}

func (a HttpAdapter) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

// use for debug only
func (a HttpAdapter) GetInfo() map[string]string {
	return a.server.GetInfo()
}
