package repository

import (
	"errors"
	"fmt"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

func NewRepositoryError(cause error) error {
	return fmt.Errorf("repository error: %w", cause)
}
