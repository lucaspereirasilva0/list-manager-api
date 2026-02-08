package handlers

import (
	"errors"
	"net/http"

	"github.com/lucaspereirasilva0/list-manager-api/internal/service"
)

type ErrorAPI struct {
	Cause   string `json:"cause"`
	Message string `json:"message"`
	HTTP    int    `json:"http"`
}

var (
	ErrIDRequired = errors.New("id is required")
)

func (e ErrorAPI) Error() string {
	if e.Cause == "" {
		return e.Message
	}
	return e.Cause
}

func NewDecodeRequestError(err error) ErrorAPI {
	return ErrorAPI{
		Cause:   err.Error(),
		Message: "failed to decode request",
		HTTP:    http.StatusBadRequest,
	}
}

func NewInternalServerError(err error) ErrorAPI {
	return ErrorAPI{
		Cause:   err.Error(),
		Message: "internal server error",
		HTTP:    http.StatusInternalServerError,
	}
}

func HandleError(w http.ResponseWriter, err error) ErrorAPI {
	var (
		errService service.ErrorService
		errAPI     ErrorAPI
	)

	switch {
	case errors.As(err, &errAPI):
		return errAPI
	case errors.As(err, &errService):
		return ErrorAPI{
			Cause:   errService.Cause.Error(),
			Message: errService.Message,
			HTTP:    errService.HTTP,
		}
	}

	return NewInternalServerError(err)
}
