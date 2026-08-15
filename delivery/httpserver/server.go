package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"telegram-service-platform/logger"
)

type Server struct {
	server       *http.Server
	postgresPing func(ctx context.Context) error
	redisPing    func(ctx context.Context) error
	port         int
}

func New(port int, postgresPing, redisPing func(ctx context.Context) error) *Server {
	mux := http.NewServeMux()

	s := &Server{
		postgresPing: postgresPing,
		redisPing:    redisPing,
		port:         port,
	}

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/ready", s.readyHandler)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	logger.Logger.Info("http server starting", zap.Int("port", s.port))
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Error("http server error", zap.Error(err))
		}
	}()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	logger.Logger.Info("http server shutting down")
	return s.server.Shutdown(shutdownCtx)
}
