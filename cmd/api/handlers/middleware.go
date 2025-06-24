package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	body        *bytes.Buffer // Add a buffer to capture the response body
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, body: new(bytes.Buffer)}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(buf []byte) (int, error) {
	rw.body.Write(buf)
	return rw.ResponseWriter.Write(buf)
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		start := time.Now()
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					wrapped := wrapResponseWriter(w)
					logger.Error("request not completed",
						zap.Int("status", wrapped.Status()),
						zap.String("method", r.Method),
						zap.String("path", r.URL.EscapedPath()),
						zap.Any("response", wrapped.body.String()),
						zap.Duration("duration", time.Since(start)),
						zap.Any("error", err),
					)
				}
			}()

			wrapped := wrapResponseWriter(w)
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

// ErrorHandlingMiddleware handles errors that occur during request processing.
func ErrorHandlingMiddleware(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrappedWriter := wrapResponseWriter(w)
		err := handler(wrappedWriter, r)

		defer func() {
			if recoverableError := recover(); recoverableError != nil {
				// If a panic occurs, use the recovered value as the error
				errFromPanic, ok := recoverableError.(error)
				if !ok {
					errFromPanic = fmt.Errorf("%v", recoverableError)
				}
				if !wrappedWriter.wroteHeader {
					errAPI := handleError(wrappedWriter, errFromPanic)
					http.Error(wrappedWriter, errAPI.Message, errAPI.HTTP)
				}
				return
			}

			// If the handler returned an error (and no panic occurred), handle it here
			if err != nil && !wrappedWriter.wroteHeader {
				errAPI := handleError(wrappedWriter, err)
				http.Error(wrappedWriter, errAPI.Message, errAPI.HTTP)
			}
		}()
	}
}
