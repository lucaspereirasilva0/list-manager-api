package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lucaspereirasilva0/list-manager-api/internal/repository"
)

type ErrorService struct {
	Cause   error
	Message string
	HTTP    int
}

func (e ErrorService) Error() string {
	return fmt.Sprintf("message: %s, cause: %s", e.Message, e.Cause)
}

func NewErrorService(cause error, message string, http int) error {
	return ErrorService{
		Cause:   cause,
		Message: message,
		HTTP:    http,
	}
}

func NewErrorEmptyItem() error {
	return NewErrorService(nil, "item is empty", http.StatusBadRequest)
}

func handleError(err error) error {
	var (
		errService ErrorService
	)
	switch {
	case errors.Is(err, repository.ErrItemNotFound):
		return NewErrorService(err, repository.ErrItemNotFound.Error(), http.StatusNotFound)
	case errors.As(err, &errService):
		return err
	}

	return NewErrorService(err, "internal server error", http.StatusInternalServerError)
}
