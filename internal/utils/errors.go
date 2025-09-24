package utils

import "fmt"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func WrapAppError(code int, format string, a ...interface{}) *AppError {
	return &AppError{Code: code, Message: fmt.Sprintf(format, a...)}
}
