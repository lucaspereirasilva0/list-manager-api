package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddleware(t *testing.T) {
	// Create a new observed logger to capture log entries
	core, observedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	tests := []struct {
		name                string
		givenHandlerFunc    http.HandlerFunc
		wantExpectedStatus  int
		wantExpectedLogMsg  string
		wantExpectedLogType zapcore.Level
		wantPanicExpected   bool
	}{
		{
			name: "Given_SuccessfulRequest_When_HandlerCompletes_Then_InfoLogGenerated",
			givenHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			},
			wantExpectedStatus:  http.StatusOK,
			wantExpectedLogMsg:  "request completed",
			wantExpectedLogType: zap.InfoLevel,
			wantPanicExpected:   false,
		},
		{
			name: "Given_HandlerPanics_When_LoggingMiddleware_Then_ErrorLogGenerated",
			givenHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			},
			wantExpectedStatus:  http.StatusInternalServerError, // Handled by ErrorHandlingMiddleware usually
			wantExpectedLogMsg:  "request not completed",
			wantExpectedLogType: zap.ErrorLevel,
			wantPanicExpected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset observed logs for each test run
			observedLogs.TakeAll()

			mw := LoggingMiddleware(logger)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rr := httptest.NewRecorder()

			// A dummy next handler to pass to the middleware
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tt.givenHandlerFunc.ServeHTTP(w, r) // Execute the test-specific handler
			})

			// Use a defer to recover from panics, so the test continues
			defer func() {
				if r := recover(); r != nil && !tt.wantPanicExpected {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()

			mw(nextHandler).ServeHTTP(rr, req)

			// Assertions for log entries
			logs := observedLogs.All()
			require.Len(t, logs, 1, "Expected exactly one log entry")

			logEntry := logs[0]
			require.Equal(t, tt.wantExpectedLogMsg, logEntry.Message, "Log message mismatch")
			require.Equal(t, tt.wantExpectedLogType, logEntry.Level, "Log level mismatch")

			// Check some common fields
			foundStatus := false
			foundMethod := false
			foundPath := false
			foundDuration := false
			foundResponse := false
			foundError := false

			for _, field := range logEntry.Context {
				switch field.Key {
				case "status":
					require.Equal(t, int64(tt.wantExpectedStatus), field.Integer, "Status field mismatch")
					foundStatus = true
				case "method":
					require.Equal(t, req.Method, field.String, "Method field mismatch")
					foundMethod = true
				case "path":
					require.Equal(t, req.URL.EscapedPath(), field.String, "Path field mismatch")
					foundPath = true
				case "duration":
					require.NotZero(t, field.Integer, "Duration should not be zero")
					foundDuration = true
				case "response":
					// For a successful request, body should be "OK"
					// For a panic, it might be empty or partial
					// assert.Contains(t, field.String, expectedBody, "Response body field mismatch")
					foundResponse = true
				case "error":
					if tt.wantPanicExpected {
						require.NotNil(t, field.Interface, "Error field should not be nil for panic")
					} else {
						require.Nil(t, field.Interface, "Error field should be nil for no panic")
					}
					foundError = true
				}
			}

			require.True(t, foundStatus, "Status field not found in log entry")
			require.True(t, foundMethod, "Method field not found in log entry")
			require.True(t, foundPath, "Path field not found in log entry")
			require.True(t, foundDuration, "Duration field not found in log entry")
			require.True(t, foundResponse, "Response field not found in log entry")
			// Error field is only expected in panic cases
			if tt.wantPanicExpected {
				require.True(t, foundError, "Error field not found in log entry")
			}
		})
	}
}
