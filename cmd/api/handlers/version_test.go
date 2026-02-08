package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name           string
		givenRequest   *http.Request
		givenWriter    *httptest.ResponseRecorder
		wantStatusCode int
		wantResponse   VersionResponse
		wantHeaders    map[string]string
	}{
		{
			name:           "Given_ValidRequest_When_GetVersionCalled_Then_ReturnsCorrectVersion",
			givenRequest:   mockHTTPRequest(),
			givenWriter:    mockHTTPResponseWriter(),
			wantStatusCode: http.StatusOK,
			wantResponse:   mockVersionResponse(),
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			req := tt.givenRequest
			w := tt.givenWriter

			// Act
			GetVersion(w, req)

			// Assert
			// Check status code
			assert.Equal(t, tt.wantStatusCode, w.Code, "Status code should match expected")

			// Check headers
			for headerName, expectedValue := range tt.wantHeaders {
				actualValue := w.Header().Get(headerName)
				assert.Equal(t, expectedValue, actualValue, "Header %s should match expected", headerName)
			}

			// Check JSON response structure
			var response VersionResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Check version value
			assert.Equal(t, tt.wantResponse.Version, response.Version, "Version should match expected")
		})
	}
}

func TestGetVersion_ResponseStructure(t *testing.T) {
	tests := []struct {
		name         string
		givenRequest *http.Request
		givenWriter  *httptest.ResponseRecorder
		wantFields   []string
	}{
		{
			name:         "Given_ValidRequest_When_GetVersionCalled_Then_ReturnsValidJSONStructure",
			givenRequest: mockHTTPRequest(),
			givenWriter:  mockHTTPResponseWriter(),
			wantFields:   []string{"version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			req := tt.givenRequest
			w := tt.givenWriter

			// Act
			GetVersion(w, req)

			// Assert
			// Check that response is valid JSON
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			// Check that all expected fields are present
			for _, field := range tt.wantFields {
				_, exists := response[field]
				assert.True(t, exists, "Response should contain field: %s", field)
			}

			// Check that version field is a string
			versionValue, exists := response["version"]
			assert.True(t, exists, "Response should contain version field")
			assert.IsType(t, "", versionValue, "Version field should be a string")
		})
	}
}

// mockVersionResponse returns the expected version response for testing
func mockVersionResponse() VersionResponse {
	return VersionResponse{
		Version: "1.0.0",
	}
}

// mockHTTPRequest creates a mock HTTP request for testing
func mockHTTPRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/version", nil)
}

// mockHTTPResponseWriter creates a response recorder to capture the response
func mockHTTPResponseWriter() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}
