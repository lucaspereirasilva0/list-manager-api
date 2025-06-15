package domain

import (
	"github.com/google/uuid"
)

// Item represents the main domain entity
type Item struct {
	ID     string
	Name   string
	Active bool
}

// NewItem creates a new instance of Item
func NewItem(name string, active bool) Item {
	return Item{
		ID:     uuid.New().String(),
		Name:   name,
		Active: active,
	}
}

func (i Item) IsActive() bool {
	return i.Active
}

func (i Item) IsEmpty() bool {
	return i.ID == ""
}
