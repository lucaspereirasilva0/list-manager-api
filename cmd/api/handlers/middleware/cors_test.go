package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name                  string
		givenAllowedOrigins   []string
		givenRequestOrigin    string
		givenRequestMethod    string
		wantStatusCode        int
		wantAllowOriginHeader string
		wantVaryHeader        bool
	}{
		{
			name:                  "Given_AllOriginsAllowed_When_GETRequest_Then_AllowAllOriginHeaderAndOK",
			givenAllowedOrigins:   []string{"*"},
			givenRequestOrigin:    "https://any-domain.com",
			givenRequestMethod:    http.MethodGet,
			wantStatusCode:        http.StatusOK,
			wantAllowOriginHeader: "*",
			wantVaryHeader:        true,
		},
		{
			name:                  "Given_AllOriginsAllowed_When_OPTIONSRequest_Then_AllowAllOriginHeaderAndOK",
			givenAllowedOrigins:   []string{"*"},
			givenRequestOrigin:    "https://another-domain.net",
			givenRequestMethod:    http.MethodOptions,
			wantStatusCode:        http.StatusOK,
			wantAllowOriginHeader: "*",
			wantVaryHeader:        true,
		},
		{
			name:                  "Given_SpecificOriginAllowed_When_GETRequestWithAllowedOrigin_Then_SpecificOriginHeaderAndOK",
			givenAllowedOrigins:   []string{"https://allowed.com"},
			givenRequestOrigin:    "https://allowed.com",
			givenRequestMethod:    http.MethodGet,
			wantStatusCode:        http.StatusOK,
			wantAllowOriginHeader: "https://allowed.com",
			wantVaryHeader:        true,
		},
		{
			name:                  "Given_SpecificOriginAllowed_When_OPTIONSRequestWithAllowedOrigin_Then_SpecificOriginHeaderAndOK",
			givenAllowedOrigins:   []string{"https://allowed.com"},
			givenRequestOrigin:    "https://allowed.com",
			givenRequestMethod:    http.MethodOptions,
			wantStatusCode:        http.StatusOK,
			wantAllowOriginHeader: "https://allowed.com",
			wantVaryHeader:        true,
		},
		{
			name:                  "Given_SpecificOriginAllowed_When_GETRequestWithDisallowedOrigin_Then_NoOriginHeaderAndOK",
			givenAllowedOrigins:   []string{"https://allowed.com"},
			givenRequestOrigin:    "https://disallowed.com",
			givenRequestMethod:    http.MethodGet,
			wantStatusCode:        http.StatusOK,
			wantAllowOriginHeader: "", // Should not set Access-Control-Allow-Origin
			wantVaryHeader:        true,
		},
		{
			name:                  "Given_SpecificOriginAllowed_When_OPTIONSRequestWithDisallowedOrigin_Then_NoOriginHeaderAndOK",
			givenAllowedOrigins:   []string{"https://allowed.com"},
			givenRequestOrigin:    "https://disallowed.com",
			givenRequestMethod:    http.MethodOptions,
			wantStatusCode:        http.StatusOK,
			wantAllowOriginHeader: "", // Should not set Access-Control-Allow-Origin
			wantVaryHeader:        true,
		},
		{
			name:                  "Given_MultipleOriginsAllowed_When_GETRequestWithOneAllowedOrigin_Then_SpecificOriginHeaderAndOK",
			givenAllowedOrigins:   []string{"https://one.com", "https://two.com"},
			givenRequestOrigin:    "https://one.com",
			givenRequestMethod:    http.MethodGet,
			wantStatusCode:        http.StatusOK,
			wantAllowOriginHeader: "https://one.com",
			wantVaryHeader:        true,
		},
		{
			name:                  "Given_MultipleOriginsAllowed_When_OPTIONSRequestWithOneAllowedOrigin_Then_SpecificOriginHeaderAndOK",
			givenAllowedOrigins:   []string{"https://one.com", "https://two.com"},
			givenRequestOrigin:    "https://two.com",
			givenRequestMethod:    http.MethodOptions,
			wantStatusCode:        http.StatusOK,
			wantAllowOriginHeader: "https://two.com",
			wantVaryHeader:        true,
		},
		{
			name:                  "Given_NoOriginHeaderInRequest_When_GETRequest_Then_AllowOriginHeaderNotSetAndOK",
			givenAllowedOrigins:   []string{"https://allowed.com"},
			givenRequestOrigin:    "", // No Origin header
			givenRequestMethod:    http.MethodGet,
			wantStatusCode:        http.StatusOK,
			wantAllowOriginHeader: "", // Should not set Access-Control-Allow-Origin if origin is not in allowed list
			wantVaryHeader:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			middleware := CORSMiddleware(tt.givenAllowedOrigins)
			handlerToTest := middleware(nextHandler)

			req := httptest.NewRequest(tt.givenRequestMethod, "http://example.com/foo", nil)
			if tt.givenRequestOrigin != "" {
				req.Header.Set("Origin", tt.givenRequestOrigin)
			}
			rr := httptest.NewRecorder()

			// Act
			handlerToTest.ServeHTTP(rr, req)

			// Assert
			require.Equal(t, tt.wantStatusCode, rr.Code, "Expected status code mismatch")

			if tt.wantAllowOriginHeader != "" {
				require.Equal(t, tt.wantAllowOriginHeader, rr.Header().Get("Access-Control-Allow-Origin"), "Access-Control-Allow-Origin header mismatch")
			} else {
				require.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"), "Access-Control-Allow-Origin header should be empty")
			}

			require.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"), "Access-Control-Allow-Methods header mismatch")
			require.Equal(t, "Content-Type, Authorization, ngrok-skip-browser-warning", rr.Header().Get("Access-Control-Allow-Headers"), "Access-Control-Allow-Headers header mismatch")
			require.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"), "Access-Control-Allow-Credentials header mismatch")
			if tt.wantVaryHeader {
				require.Equal(t, "Origin", rr.Header().Get("Vary"), "Vary header mismatch")
			}
		})
	}
}
