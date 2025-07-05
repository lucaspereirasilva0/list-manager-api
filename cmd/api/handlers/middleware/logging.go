package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		start := time.Now()
		fn := func(w http.ResponseWriter, r *http.Request) {
			wrapped := wrapResponseWriter(w)

			defer func() {
				if err := recover(); err != nil {
					// Set status to 500 for panic cases
					wrapped.WriteHeader(http.StatusInternalServerError)
					logger.Error("request not completed",
						zap.Int("status", wrapped.Status()),
						zap.String("method", r.Method),
						zap.String("path", r.URL.EscapedPath()),
						zap.Any("response", wrapped.body.String()),
						zap.Duration("duration", time.Since(start)),
						zap.Error(fmt.Errorf("%v", err)),
					)
				}
			}()

			next.ServeHTTP(wrapped, r)
			logger.Info("request completed",
				zap.Int("status", wrapped.status),
				zap.String("method", r.Method),
				zap.String("path", r.URL.EscapedPath()),
				zap.Any("response", wrapped.body.String()),
				zap.Duration("duration", time.Since(start)),
			)
		}

		return http.HandlerFunc(fn)
	}
}
