package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Util787/task-manager/internal/app"
	"github.com/Util787/task-manager/internal/config"
	"github.com/Util787/task-manager/pkg/logger/handlers/slogpretty"
	"github.com/Util787/task-manager/pkg/logger/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title           Task Manager API
// @version         1.0
// @description     Task Manager API

// @host      localhost:8080
// @BasePath  /api

func main() {

	config, err := config.Load()
	if err != nil {
		panic("Failed to load config" + err.Error())
	}

	log := setupLogger(config.Env)

	app := app.New(*config, log)

	go func() {
		err := app.HttpAdapter.Start()
		if err != nil {
			log.Error("Server was interrupted", sl.Err(err))
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Shutting down the server")
	if err := app.HttpAdapter.Shutdown(context.Background()); err != nil {
		log.Error("Failed to shut down the server", sl.Err(err))
	}

	log.Info("Gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
