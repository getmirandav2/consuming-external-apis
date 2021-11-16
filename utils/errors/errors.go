package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type APIError interface {
	Status() int
	Message() string
	Error() string
}

type apiError struct {
	ASStatus  int    `json:"status"`
	ASMessage string `json:"message"`
	AnErr     string `json:"error,omitempty"`
}

func (e *apiError) Status() int {
	return e.ASStatus
}

func (e *apiError) Message() string {
	return e.ASMessage
}

func (e *apiError) Error() string {
	return e.AnErr
}

func NewAPIErrorFromBytes(body []byte) (APIError, error) {
	var result apiError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("invalid JSON body")
	}
	return &result, nil
}

func NewAPIError(code int, msg string) APIError {
	return &apiError{
		ASStatus:  code,
		ASMessage: msg,
	}
}

func NewInternalServerError(msg string) APIError {
	return &apiError{
		ASStatus:  http.StatusInternalServerError,
		ASMessage: msg,
	}
}

func NewNotFoundError(msg string) APIError {
	return &apiError{
		ASStatus:  http.StatusNotFound,
		ASMessage: msg,
	}
}

func NewBadRequestError(msg string) APIError {
	return &apiError{
		ASStatus:  http.StatusBadRequest,
		ASMessage: msg,
	}
}
