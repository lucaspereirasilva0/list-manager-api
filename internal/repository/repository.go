package repository

import (
	"context"
)

// ItemRepository defines the interface for item persistence operations
type ItemRepository interface {
	// Create inserts a new item in the repository
	Create(ctx context.Context, item Item) (Item, error)

	// Update modifies an existing item in the repository
	Update(ctx context.Context, item Item) (Item, error)

	// Delete removes an item from the repository
	Delete(ctx context.Context, id string) error

	// GetByID retrieves an item by its ID
	GetByID(ctx context.Context, id string) (Item, error)

	// List retrieves all items from the repository
	List(ctx context.Context) ([]Item, error)

	// BulkUpdateActive updates the active field for all items in the repository
	BulkUpdateActive(ctx context.Context, active bool) (matchedCount int64, modifiedCount int64, err error)
}
