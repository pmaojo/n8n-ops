package client

import (
	"errors"
	"fmt"
	"net/http"
)

// Predefined error types for better error handling
var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrBadRequest    = errors.New("bad request")
	ErrServerError   = errors.New("internal server error") 
	ErrTimeout       = errors.New("request timeout")
	ErrConnection    = errors.New("connection failed")
	ErrInvalidConfig = errors.New("invalid configuration")
)

// APIError represents a structured API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error %d: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

// NewAPIError creates an APIError from HTTP status code
func NewAPIError(statusCode int, message string) *APIError {
	baseError := &APIError{
		Code:    statusCode,
		Message: message,
	}
	
	// Map common HTTP status codes to specific errors
	switch statusCode {
	case http.StatusNotFound:
		return &APIError{Code: statusCode, Message: "resource not found"}
	case http.StatusUnauthorized:
		return &APIError{Code: statusCode, Message: "unauthorized access"}
	case http.StatusBadRequest:
		return &APIError{Code: statusCode, Message: "bad request"}
	case http.StatusInternalServerError:
		return &APIError{Code: statusCode, Message: "internal server error"}
	default:
		return baseError
	}
}

// IsAPIError checks if error is an APIError with specific status code
func IsAPIError(err error, statusCode int) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == statusCode
	}
	return false
}

// IsNotFound checks if error represents a "not found" condition
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) || IsAPIError(err, http.StatusNotFound)
}

// IsUnauthorized checks if error represents an "unauthorized" condition  
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized) || IsAPIError(err, http.StatusUnauthorized)
}

// IsBadRequest checks if error represents a "bad request" condition
func IsBadRequest(err error) bool {
	return errors.Is(err, ErrBadRequest) || IsAPIError(err, http.StatusBadRequest)
}

// IsServerError checks if error represents a server error condition
func IsServerError(err error) bool {
	return errors.Is(err, ErrServerError) || IsAPIError(err, http.StatusInternalServerError)
}