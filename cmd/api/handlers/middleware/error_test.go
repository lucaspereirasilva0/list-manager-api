package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucaspereirasilva0/list-manager-api/cmd/api/handlers"
	"github.com/lucaspereirasilva0/list-manager-api/internal/service"

	"github.com/stretchr/testify/require"
)

func TestErrorHandlingMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		givenHandlerFunc func(http.ResponseWriter, *http.Request) error
		wantStatus       int
		wantBody         string
	}{
		{
			name: "Given_NoError_When_HandlerCompletes_Then_StatusOkAndSuccessBody",
			givenHandlerFunc: func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("Success"))
				return nil
			},
			wantStatus: http.StatusOK,
			wantBody:   "Success",
		},
		{
			name: "Given_APIError_When_HandlerReturnsAPIError_Then_APIErrorHandled",
			givenHandlerFunc: func(w http.ResponseWriter, r *http.Request) error {
				return handlers.NewInternalServerError(errors.New("test error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal server error",
		},
		{
			name: "Given_ServiceError_When_HandlerReturnsServiceError_Then_APIErrorHandled",
			givenHandlerFunc: func(w http.ResponseWriter, r *http.Request) error {
				return service.NewErrorService(errors.New("item not found"), "item not found", service.ServiceSource, http.StatusNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   "item not found",
		},
		{
			name: "Given_Panic_When_HandlerPanics_Then_InternalServerError",
			givenHandlerFunc: func(w http.ResponseWriter, r *http.Request) error {
				panic("something went wrong!")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal server error",
		},
		{
			name: "Given_AlreadyWrittenHeader_When_HandlerErrors_Then_OriginalStatusPreserved",
			givenHandlerFunc: func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte("Partial content"))
				return errors.New("some error after writing")
			},
			wantStatus: http.StatusAccepted,
			wantBody:   "Partial content",
		},
		{
			name: "Given_NonErrorPanic_When_HandlerPanicsWithNonError_Then_InternalServerError",
			givenHandlerFunc: func(w http.ResponseWriter, r *http.Request) error {
				panic(123)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rr := httptest.NewRecorder()

			mw := ErrorHandlingMiddleware(tt.givenHandlerFunc)
			mw.ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Result().StatusCode, "Expected status mismatch")

			if tt.wantStatus == http.StatusInternalServerError && (tt.name == "Given_Panic_When_HandlerPanics_Then_InternalServerError" || tt.name == "Given_NonErrorPanic_When_HandlerPanicsWithNonError_Then_InternalServerError") {
				require.Contains(t, rr.Body.String(), "internal server error", "Expected body for internal server error mismatch")
			} else {
				require.Contains(t, rr.Body.String(), tt.wantBody, "Expected body mismatch")
			}
		})
	}
}
