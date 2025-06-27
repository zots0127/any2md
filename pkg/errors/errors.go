package errors

import "fmt"

type ConversionError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

func (e *ConversionError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewValidationError(message string) *ConversionError {
	return &ConversionError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Details: make(map[string]interface{}),
	}
}

func NewParsingError(message string, details map[string]interface{}) *ConversionError {
	return &ConversionError{
		Code:    "PARSING_ERROR",
		Message: message,
		Details: details,
	}
}

func NewInternalError(message string) *ConversionError {
	return &ConversionError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Details: make(map[string]interface{}),
	}
}