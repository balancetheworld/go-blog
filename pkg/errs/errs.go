package errs

import "net/http"

type ServiceError struct {
	Code       int
	HTTPStatus int
	Message    string
}

func (e *ServiceError) Error() string {
	return e.Message
}

func newServiceError(code, httpStatus int, message string) *ServiceError {
	return &ServiceError{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    message,
	}
}

func NewBadRequest(code int, message string) *ServiceError {
	return newServiceError(code, http.StatusBadRequest, message)
}

func NewUnauthorized(code int, message string) *ServiceError {
	return newServiceError(code, http.StatusUnauthorized, message)
}

func NewForbidden(code int, message string) *ServiceError {
	return newServiceError(code, http.StatusForbidden, message)
}

func NewNotFound(code int, message string) *ServiceError {
	return newServiceError(code, http.StatusNotFound, message)
}

func NewConflict(code int, message string) *ServiceError {
	return newServiceError(code, http.StatusConflict, message)
}

func NewInternalServer(code int, message string) *ServiceError {
	return newServiceError(code, http.StatusInternalServerError, message)
}
