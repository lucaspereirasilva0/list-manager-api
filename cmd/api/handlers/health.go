package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	dbmongo "github.com/lucaspereirasilva0/list-manager-api/internal/database/mongodb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type healthHandler struct {
	mongoClient dbmongo.PingClientOperations
	logger      *zap.Logger
}

type HealthHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request) error
}

func NewHealthHandler(mongoClient dbmongo.PingClientOperations, logger *zap.Logger) HealthHandler {
	return &healthHandler{
		mongoClient: mongoClient,
		logger:      logger,
	}
}

func (h *healthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	dbStatus, dbCheck := h.checkDatabase(ctx)

	response := h.buildResponse(timestamp, dbStatus, dbCheck)
	h.logHealthCheck(r.Method, r.URL.Path, response)

	return h.writeResponse(w, r, response)
}

func (h *healthHandler) checkDatabase(ctx context.Context) (ComponentStatus, Check) {
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := h.mongoClient.Ping(pingCtx)
	if err != nil {
		return ComponentStatusDisconnected, Check{
			Status: ComponentStatusFailed,
			Error:  err.Error(),
		}
	}

	return ComponentStatusConnected, Check{
		Status: ComponentStatusPassed,
	}
}

func (h *healthHandler) buildResponse(timestamp string, dbStatus ComponentStatus, dbCheck Check) HealthCheckResponse {
	var status HealthStatus
	if dbStatus == ComponentStatusDisconnected {
		status = HealthStatusDegraded
	} else {
		status = HealthStatusUp
	}

	return HealthCheckResponse{
		Status:    status,
		Server:    ComponentStatusUp,
		Database:  dbStatus,
		Timestamp: timestamp,
		Checks: map[string]Check{
			"database": dbCheck,
		},
	}
}

func (h *healthHandler) logHealthCheck(method string, path string, response HealthCheckResponse) {
	level := zapcore.InfoLevel
	message := "health check passed"

	if response.Status == HealthStatusDegraded {
		level = zapcore.WarnLevel
		message = "health check degraded"
	}

	h.logWithLevel(level, message, method, path, response)
}

func (h *healthHandler) logWithLevel(level zapcore.Level, message string, method string, path string, response HealthCheckResponse) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.String("status", string(response.Status)),
		zap.String("database", string(response.Database)),
	}

	if response.Status == HealthStatusDegraded {
		fields = append(fields, zap.String("error", response.Checks["database"].Error))
	}

	h.logger.Log(level, message, fields...)
}

func (h *healthHandler) writeResponse(w http.ResponseWriter, r *http.Request, response HealthCheckResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		h.logger.Error("failed to encode health check response",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Error(encodeErr),
		)
		return encodeErr
	}

	return nil
}
