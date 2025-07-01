package repository

import (
	"errors"
	"fmt"
)

var (
	ErrItemNotFound = errors.New("item not found")
	ErrNotFound     = ErrItemNotFound
	ErrInvalidHexID = errors.New("invalid hexadecimal representation of an ObjectID")
)

func NewRepositoryError(cause error) error {
	return fmt.Errorf("repository error: %w", cause)
}
