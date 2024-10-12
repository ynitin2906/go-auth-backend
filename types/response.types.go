package types

import (
	"net/http"
)

// ErrorResponse represents a structured error response for API endpoints.
type ErrorResponse struct {
	Error   bool   `json:"error"`   // Indicates if there was an error
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Error message
}

// APIResponse represents a structured response for all API endpoints.
type APIResponse struct {
	Error   bool        `json:"error"`          // Indicates if there was an error
	Code    int         `json:"code"`           // HTTP status code
	Message string      `json:"message"`        // Success or error message
	Data    interface{} `json:"data,omitempty"` // Optional data (for success responses)
}

// CreateSuccessResponse generates a success response.
func CreateSuccessResponse(message string, statusCode int, data interface{}) APIResponse {
	return APIResponse{
		Error:   false,
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
}

// CreateErrorResponse generates an error response.
func CreateErrorResponse(message string, statusCode int, data interface{}) APIResponse {
	return APIResponse{
		Error:   true,
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
}

// Error represents a structured custom error for the API.
type Error struct {
	Code int    `json:"code"`  // HTTP status code
	Err  string `json:"error"` // Error message
}

// Error implements the error interface for the custom Error struct.
func (e Error) Error() string {
	return e.Err
}

// NewError creates a new custom error.
func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

// Predefined errors for common scenarios
func ErrUnAuthorized() Error {
	return NewError(http.StatusUnauthorized, "unauthorized request")
}

func ErrResourceNotFound(resource string) Error {
	return NewError(http.StatusNotFound, resource+" resource not found")
}

func ErrBadRequest(msg ...string) Error {
	errMsg := "Invalid request"
	if len(msg) > 0 {
		errMsg = msg[0]
	}
	return NewError(http.StatusBadRequest, errMsg)
}

func ErrInvalidID() Error {
	return NewError(http.StatusBadRequest, "invalid ID provided")
}

// getValidationError extracts validation errors from the validator and converts them to ErrorResponse format.
// func GetValidationError(errs error) []ErrorResponse {
// 	validationErrors := []ErrorResponse{}

// 	for _, err := range errs.(validator.ValidationErrors) {
// 		validationErrors = append(validationErrors, ErrorResponse{
// 			Error:       true,
// 			Code:        http.StatusBadRequest,
// 			Message:     "Validation failed",
// 			FailedField: err.Field(),
// 			Tag:         err.Tag(),
// 			Value:       err.Value(),
// 		})
// 	}

// 	return validationErrors
// }
