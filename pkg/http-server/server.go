package http_server

import (
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + cfg.Port,
			Handler:           handler,
			MaxHeaderBytes:    1 << 20, // 1 MB
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		}}
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
