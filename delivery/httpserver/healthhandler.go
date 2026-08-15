package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"telegram-service-platform/logger"
	"time"

	"go.uber.org/zap"
)

type healthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	checks := make(map[string]string)
	overall := "healthy"

	if err := s.postgresPing(ctx); err != nil {
		checks["postgres"] = "unhealthy: " + err.Error()
		overall = "unhealthy"
		logger.Logger.Warn("health check postgres failed", zap.Error(err))
	} else {
		checks["postgres"] = "healthy"
	}

	if err := s.redisPing(ctx); err != nil {
		checks["redis"] = "unhealthy: " + err.Error()
		overall = "unhealthy"
		logger.Logger.Warn("health check redis failed", zap.Error(err))
	} else {
		checks["redis"] = "healthy"
	}

	w.Header().Set("Content-Type", "application/json")
	if overall != "healthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	_ = json.NewEncoder(w).Encode(healthResponse{
		Status:    overall,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	})
}
