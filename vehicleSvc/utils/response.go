package utils

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	StatusCode int    `json:"statusCode"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

func ErrNotFound(entity string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusNotFound,
		Code:       http.StatusNotFound,
		Message:    fmt.Sprintf("%s not found", entity),
	}
}

func ErrBadRequest(msg string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusBadRequest,
		Code:       http.StatusBadRequest,
		Message:    fmt.Sprintf("%s", msg),
	}
}

func ErrConflict(entity string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusConflict,
		Code:       http.StatusConflict,
		Message:    fmt.Sprintf("%s is conflict", entity),
	}
}

func ErrForbidden(entity string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusForbidden,
		Code:       http.StatusForbidden,
		Message:    fmt.Sprintf("%s is forbidden", entity),
	}
}

func ErrInternal(msg string) *CustomError {
	return &CustomError{
		StatusCode: http.StatusInternalServerError,
		Code:       http.StatusInternalServerError,
		Message:    msg,
	}
}
