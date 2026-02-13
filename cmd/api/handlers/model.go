package handlers

import "time"

type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Active      bool      `json:"active"`
	Observation *string   `json:"observation,omitempty"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

type HealthCheckResponse struct {
	Status    HealthStatus     `json:"status"`
	Server    ComponentStatus  `json:"server"`
	Database  ComponentStatus  `json:"database"`
	Timestamp string           `json:"timestamp"`
	Checks    map[string]Check `json:"checks"`
}

type Check struct {
	Status ComponentStatus `json:"status"`
	Error  string          `json:"error,omitempty"`
}

type HealthStatus string

const (
	HealthStatusUp       HealthStatus = "up"
	HealthStatusDegraded HealthStatus = "degraded"
	HealthStatusDown     HealthStatus = "down"
)

type BulkActiveRequest struct {
	Active bool `json:"active"`
}

type BulkActiveResponse struct {
	MatchedCount  int64 `json:"matchedCount"`
	ModifiedCount int64 `json:"modifiedCount"`
}

type ComponentStatus string

const (
	ComponentStatusUp           ComponentStatus = "up"
	ComponentStatusConnected    ComponentStatus = "connected"
	ComponentStatusDisconnected ComponentStatus = "disconnected"
	ComponentStatusPassed       ComponentStatus = "passed"
	ComponentStatusFailed       ComponentStatus = "failed"
)
