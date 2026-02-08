package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/handlers"
	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/handlers/middleware"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockPingClient struct {
	mock.Mock
}

func (m *mockPingClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name               string
		givenPingErr       error
		wantHTTPStatus     int
		wantOverallStatus  string
		wantServerStatus   string
		wantDatabaseStatus string
		wantCheckStatus    string
		wantCheckError     string
	}{
		{
			name:               "Given_MongoDBConnected_When_HealthCheck_Then_ReturnsUpStatus",
			givenPingErr:       nil,
			wantHTTPStatus:     http.StatusOK,
			wantOverallStatus:  "up",
			wantServerStatus:   "up",
			wantDatabaseStatus: "connected",
			wantCheckStatus:    "passed",
		},
		{
			name:               "Given_MongoDBDisconnected_When_HealthCheck_Then_ReturnsDegradedStatus",
			givenPingErr:       errors.New("connection failed"),
			wantHTTPStatus:     http.StatusOK,
			wantOverallStatus:  "degraded",
			wantServerStatus:   "up",
			wantDatabaseStatus: "disconnected",
			wantCheckStatus:    "failed",
			wantCheckError:     "connection failed",
		},
		{
			name:               "Given_ContextTimeout_When_HealthCheck_Then_ReturnsDegradedStatus",
			givenPingErr:       errors.New("context deadline exceeded"),
			wantHTTPStatus:     http.StatusOK,
			wantOverallStatus:  "degraded",
			wantServerStatus:   "up",
			wantDatabaseStatus: "disconnected",
			wantCheckStatus:    "failed",
			wantCheckError:     "context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMongo := new(mockPingClient)
			mockMongo.On("Ping", mock.Anything).Return(tt.givenPingErr)

			logger, _ := zap.NewDevelopment()
			healthHandler := handlers.NewHealthHandler(mockMongo, logger)

			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			rec := httptest.NewRecorder()

			handlerWithMiddleware := middleware.ErrorHandlingMiddleware(healthHandler.HealthCheck)
			handlerWithMiddleware.ServeHTTP(rec, req)

			require.Equal(t, tt.wantHTTPStatus, rec.Code)

			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)

			require.Equal(t, tt.wantOverallStatus, response["status"])
			require.Equal(t, tt.wantServerStatus, response["server"])
			require.Equal(t, tt.wantDatabaseStatus, response["database"])

			checks := response["checks"].(map[string]interface{})
			dbCheck := checks["database"].(map[string]interface{})
			require.Equal(t, tt.wantCheckStatus, dbCheck["status"])

			if tt.wantCheckError != "" {
				require.Equal(t, tt.wantCheckError, dbCheck["error"])
			}
		})
	}
}

func TestHealthCheckResponseFields(t *testing.T) {
	mockMongo := new(mockPingClient)
	mockMongo.On("Ping", mock.Anything).Return(nil)

	logger, _ := zap.NewDevelopment()
	healthHandler := handlers.NewHealthHandler(mockMongo, logger)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handlerWithMiddleware := middleware.ErrorHandlingMiddleware(healthHandler.HealthCheck)
	handlerWithMiddleware.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Contains(t, response, "status")
	require.Contains(t, response, "server")
	require.Contains(t, response, "database")
	require.Contains(t, response, "timestamp")
	require.Contains(t, response, "checks")
}
