package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lucaspereirasilva0/list-manager-api/internal/repository"
)

const (
	RepositorySource = "repository"
	ServiceSource    = "service"

	_errEmptyItem = "item is empty"
)

type ErrorService struct {
	Cause   error
	Message string
	Source  string
	HTTP    int
}

func (e ErrorService) Error() string {
	return fmt.Sprintf("message: %s, cause: %s", e.Message, e.Cause)
}

func (e ErrorService) Unwrap() error {
	return e.Cause
}

func NewErrorService(cause error, message, source string, http int) error {
	return ErrorService{
		Cause:   cause,
		Message: message,
		HTTP:    http,
		Source:  source,
	}
}

func NewErrorEmptyItem() error {
	return ErrorService{
		Message: _errEmptyItem,
		Source:  ServiceSource,
		HTTP:    http.StatusBadRequest,
	}
}

func handleError(err error) error {
	var (
		errService    ErrorService
		errRepository repository.Error
	)
	switch {
	case errors.As(err, &errRepository):
		// For generic repository errors (HTTP 500), use "internal server error"
		// For specific repository errors (like 404), keep the original message
		message := errRepository.Message
		if errRepository.HTTP == http.StatusInternalServerError {
			message = "internal server error"
		}
		return NewErrorService(err, message, RepositorySource, errRepository.HTTP)
	case errors.As(err, &errService):
		return err
	}

	return NewErrorService(err, "internal server error", ServiceSource, http.StatusInternalServerError)
}
