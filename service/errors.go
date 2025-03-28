package service

import (
	"errors"
	"fmt"
)

var ErrConflict = errors.New("Conflict")
var ErrUnprocessable = errors.New("Invalid")
var ErrUnauthorized = errors.New("Authorized")

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

func NewValidationErrors() *validationErrors {
	return &validationErrors{make(map[string]ValidationError)}
}

func NewValidationErrorsSingle(field, message string, sentinel error) *validationErrors {
	err := NewValidationError(message, sentinel)
	return &validationErrors{map[string]ValidationError{field: err}}
}

func (ve *validationErrors) add(field string, err ValidationError) {
	ve.errorMap[field] = err
}

func (v *validationErrors) Errors(errorMap map[string]ValidationError) ValidationErrors {
	return &validationErrors{errorMap}
}

func (e *validationErrors) Error() string {
	return fmt.Sprintf("Validation Errors: %s", e.errorMap)
}

func (e *validationErrors) FieldToError() map[string]ValidationError {
	return e.errorMap
}

func (e *validationErrors) FieldToMessage() map[string]string {
	messageMap := make(map[string]string)
	for field, err := range e.errorMap {
		messageMap[field] = err.Message()
	}
	return messageMap
}
