package middleware

import (
	"net/http"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockLogger é uma implementação de mock para *zap.Logger
type MockLogger struct {
	mock.Mock
}

// Info implements the Info method of *zap.Logger
func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

// Error implements the Error method of *zap.Logger
func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

// MockResponseWriter is a mock implementation of http.ResponseWriter for testing.
type MockResponseWriter struct {
	HeaderCalled bool
	WrittenBytes []byte
	Status       int
	HeaderMap    http.Header
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		HeaderMap: make(http.Header),
	}
}

func (m *MockResponseWriter) Header() http.Header {
	m.HeaderCalled = true
	return m.HeaderMap
}

func (m *MockResponseWriter) Write(buf []byte) (int, error) {
	m.WrittenBytes = append(m.WrittenBytes, buf...)
	return len(buf), nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Status = statusCode
}
