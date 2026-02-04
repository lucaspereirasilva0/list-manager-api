package repository

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Cause   error
	Message string
	HTTP    int
}

func (e Error) Error() string {
	return fmt.Sprintf("message: %s, cause: %s", e.Message, e.Cause)
}

func (e Error) Unwrap() error {
	return e.Cause
}

func NewItemNotFoundError() error {
	return Error{
		Message: "item not found",
		HTTP:    http.StatusNotFound,
	}
}

func NewInvalidHexIDError() error {
	return Error{
		Message: "invalid hexadecimal representation of an ObjectID",
		HTTP:    http.StatusUnprocessableEntity,
	}
}

func NewGenericRepositoryError(cause error) error {
	return Error{
		Cause:   cause,
		Message: "generic repository error",
		HTTP:    http.StatusInternalServerError,
	}
}

func HandleError(err error) error {
	var (
		errRepository Error
	)
	switch {
	case errors.As(err, &errRepository):
		return errRepository
	}
	return NewGenericRepositoryError(err)
}
