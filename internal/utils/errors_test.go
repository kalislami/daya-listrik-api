package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppError(t *testing.T) {
	const notFound = "not found"
	err := NewAppError(404, notFound)

	assert.NotNil(t, err)
	assert.Equal(t, 404, err.Code)
	assert.Equal(t, notFound, err.Message)
	assert.Equal(t, notFound, err.Error()) // cek implementasi error interface
}

func TestWrapAppError(t *testing.T) {
	err := WrapAppError(500, "failed: %s", "database down")

	assert.NotNil(t, err)
	assert.Equal(t, 500, err.Code)
	assert.Equal(t, "failed: database down", err.Message)
	assert.Equal(t, "failed: database down", err.Error())

	// cek konsistensi kalau format kosong
	err2 := WrapAppError(400, "bad request")
	assert.Equal(t, 400, err2.Code)
	assert.Equal(t, "bad request", err2.Message)
}

func TestAppErrorImplementsError(t *testing.T) {
	var err error = NewAppError(401, "unauthorized")
	assert.Equal(t, "unauthorized", err.Error())
	assert.IsType(t, &AppError{}, err)
}

func TestWrapAppErrorFormatting(t *testing.T) {
	code := 422
	field := "email"
	expectedMsg := fmt.Sprintf("invalid field: %s", field)

	err := WrapAppError(code, "invalid field: %s", field)
	assert.Equal(t, code, err.Code)
	assert.Equal(t, expectedMsg, err.Message)
}
