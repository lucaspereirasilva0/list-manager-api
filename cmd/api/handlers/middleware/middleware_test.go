package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseWriter(t *testing.T) {
	tests := []struct {
		name           string
		op             func(*responseWriter, *httptest.ResponseRecorder)
		expectedStatus int
		expectedHeader bool
		expectedBody   []byte
		writeMultiple  bool
	}{
		{
			name: "Given_NewResponseWriter_When_Wrapped_Then_InitialStateCorrect",
			op: func(rw *responseWriter, rr *httptest.ResponseRecorder) {
				// No operation, just check initial state
			},
			expectedStatus: 0,
			expectedHeader: false,
			expectedBody:   []byte{},
		},
		{
			name: "Given_StatusSet_When_StatusQueried_Then_ReturnsCorrectStatus",
			op: func(rw *responseWriter, rr *httptest.ResponseRecorder) {
				rw.WriteHeader(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
			expectedHeader: true,
			expectedBody:   []byte{},
		},
		{
			name: "Given_WriteHeaderCalled_When_CalledAgain_Then_StatusNotOverwritten",
			op: func(rw *responseWriter, rr *httptest.ResponseRecorder) {
				rw.WriteHeader(http.StatusTeapot)
				rw.WriteHeader(http.StatusAccepted) // Should be ignored
			},
			expectedStatus: http.StatusTeapot,
			expectedHeader: true,
			expectedBody:   []byte{},
		},
		{
			name: "Given_DataWritten_When_WriteCalled_Then_DataCapturedAndWritten",
			op: func(rw *responseWriter, rr *httptest.ResponseRecorder) {
				_, _ = rw.Write([]byte("Hello, world!"))
			},
			expectedStatus: 0,
			expectedHeader: false,
			expectedBody:   []byte("Hello, world!"),
		},
		{
			name: "Given_MultipleWrites_When_WriteCalled_Then_DataAppended",
			op: func(rw *responseWriter, rr *httptest.ResponseRecorder) {
				_, _ = rw.Write([]byte("Hello,"))
				_, _ = rw.Write([]byte(" world!"))
			},
			expectedStatus: 0,
			expectedHeader: false,
			expectedBody:   []byte("Hello, world!"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			wrapped := wrapResponseWriter(rr)

			tt.op(wrapped, rr)

			assert.Equal(t, tt.expectedStatus, wrapped.Status(), "Status mismatch")
			assert.Equal(t, tt.expectedHeader, wrapped.wroteHeader, "WroteHeader mismatch")
			assert.True(t, bytes.Equal(tt.expectedBody, wrapped.body.Bytes()), "Body mismatch")

			// For tests that don't involve writing, the underlying recorder body should be empty
			if len(tt.expectedBody) == 0 && rr.Body.Len() > 0 {
				assert.Fail(t, "Underlying ResponseWriter wrote unexpected data")
			} else if len(tt.expectedBody) > 0 {
				assert.True(t, bytes.Equal(tt.expectedBody, rr.Body.Bytes()), "Underlying ResponseWriter body mismatch")
			}
		})
	}
}
