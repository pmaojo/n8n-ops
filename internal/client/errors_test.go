package client

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewAPIError(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{http.StatusNotFound, "resource not found"},
		{http.StatusUnauthorized, "unauthorized access"},
		{http.StatusBadRequest, "bad request"},
		{http.StatusInternalServerError, "internal server error"},
		{http.StatusTeapot, "custom"},
	}
	for _, tt := range tests {
		err := NewAPIError(tt.code, "custom")
		if err.Code != tt.code {
			t.Errorf("expected code %d got %d", tt.code, err.Code)
		}
		if err.Message != tt.expected {
			t.Errorf("expected message %q got %q", tt.expected, err.Message)
		}
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{Code: 418, Message: "short"}
	if got := err.Error(); got == "" {
		t.Fatal("expected non-empty error string")
	}
}

func TestIsAPIError(t *testing.T) {
	err := &APIError{Code: http.StatusNotFound}
	if !IsAPIError(err, http.StatusNotFound) {
		t.Error("expected true for matching code")
	}
	if IsAPIError(err, http.StatusBadRequest) {
		t.Error("expected false for different code")
	}
}

func TestIsNotFound(t *testing.T) {
	if !IsNotFound(ErrNotFound) {
		t.Error("expected true for ErrNotFound")
	}
	if !IsNotFound(&APIError{Code: http.StatusNotFound}) {
		t.Error("expected true for APIError 404")
	}
	if IsNotFound(errors.New("other")) {
		t.Error("expected false for unrelated error")
	}
}

func TestIsUnauthorized(t *testing.T) {
	if !IsUnauthorized(ErrUnauthorized) {
		t.Error("expected true for ErrUnauthorized")
	}
	if !IsUnauthorized(&APIError{Code: http.StatusUnauthorized}) {
		t.Error("expected true for APIError 401")
	}
	if IsUnauthorized(errors.New("other")) {
		t.Error("expected false for unrelated error")
	}
}

func TestIsBadRequest(t *testing.T) {
	if !IsBadRequest(ErrBadRequest) {
		t.Error("expected true for ErrBadRequest")
	}
	if !IsBadRequest(&APIError{Code: http.StatusBadRequest}) {
		t.Error("expected true for APIError 400")
	}
	if IsBadRequest(errors.New("other")) {
		t.Error("expected false for unrelated error")
	}
}

func TestIsServerError(t *testing.T) {
	if !IsServerError(ErrServerError) {
		t.Error("expected true for ErrServerError")
	}
	if !IsServerError(&APIError{Code: http.StatusInternalServerError}) {
		t.Error("expected true for APIError 500")
	}
	if IsServerError(errors.New("other")) {
		t.Error("expected false for unrelated error")
	}
}
