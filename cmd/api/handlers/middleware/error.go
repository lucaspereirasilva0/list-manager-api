package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/handlers"
)

// ErrorHandlingMiddleware handles errors that occur during request processing.
func ErrorHandlingMiddleware(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		wrappedWriter := wrapResponseWriter(w)

		defer func() {
			if recoverableError := recover(); recoverableError != nil {
				// If a panic occurs, use the recovered value as the error
				errFromPanic, ok := recoverableError.(error)
				if !ok {
					errFromPanic = fmt.Errorf("%v", recoverableError)
				}
				if !wrappedWriter.wroteHeader {
					errAPI := handlers.HandleError(wrappedWriter, errFromPanic)
					wrappedWriter.Header().Set("Content-Type", "application/json")
					wrappedWriter.WriteHeader(errAPI.HTTP)
					_ = json.NewEncoder(wrappedWriter).Encode(errAPI)
				}
				return
			}

			// If the handler returned an error (and no panic occurred), handle it here
			if err != nil && !wrappedWriter.wroteHeader {
				errAPI := handlers.HandleError(wrappedWriter, err)
				wrappedWriter.Header().Set("Content-Type", "application/json")
				wrappedWriter.WriteHeader(errAPI.HTTP)
				_ = json.NewEncoder(wrappedWriter).Encode(errAPI)
			}
		}()

		err = handler(wrappedWriter, r)
	}
}
