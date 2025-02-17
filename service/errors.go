package service

import (
	"errors"
	"fmt"
)

var ErrConflict = errors.New("Conflict")
var ErrUnprocessable = errors.New("Invalid")

type ValidationError interface {
	Is(target error) bool
	Message() string
	error
}

type validationError struct {
	VMessage  string
	ErrorType error
}

func (e *validationError) Error() string {
	return fmt.Sprintf("Validation Error: %s", e.VMessage)
}

func (e *validationError) Message() string {
	return e.VMessage
}

func (e *validationError) Is(target error) bool {
	return errors.Is(e.ErrorType, target)
}

func NewValidationError(message string, sentinel error) *validationError {
	return &validationError{message, sentinel}
}

type ValidationErrors interface {
	FieldToMessage() map[string]string
	FieldToError() map[string]ValidationError
	error
}
type validationErrors struct {
	errorMap map[string]ValidationError
}

func NewValidationErrors(errors map[string]ValidationError) *validationErrors {
	return &validationErrors{errors}
}

func (e *validationErrors) Error() string {
	return fmt.Sprintf("Validation Errors: %s", e.FieldToError)
}

func (e *validationErrors) FieldToError() map[string]ValidationError {
	return e.errorMap
}

func (e *validationErrors) FieldToMessage() map[string]string {
	messages := make(map[string]string, 3)
	for field, err := range e.errorMap {
		messages[field] = err.Message()
	}

	return messages
}
