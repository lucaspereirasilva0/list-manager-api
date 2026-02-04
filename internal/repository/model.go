package repository

import "time"

type Item struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Active    bool      `json:"active" bson:"active"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// User represents a user in the repository, mapped to MongoDB collection
type User struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	CreatedBy string `json:"created_by" bson:"created_by"`
}
