package apierror

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Code represents a machine-readable error code.
type Code string

const (
	CodeBadRequest          Code = "BAD_REQUEST"
	CodeUnauthorized        Code = "UNAUTHORIZED"
	CodeForbidden           Code = "FORBIDDEN"
	CodeNotFound            Code = "NOT_FOUND"
	CodeConflict            Code = "CONFLICT"
	CodeRateLimitExceeded   Code = "RATE_LIMIT_EXCEEDED"
	CodeValidationFailed    Code = "VALIDATION_FAILED"
	CodeInternalError       Code = "INTERNAL_ERROR"
	CodeServiceUnavailable  Code = "SERVICE_UNAVAILABLE"
	CodeIdempotencyConflict Code = "IDEMPOTENCY_CONFLICT"
)

// APIError is a structured error response that all API endpoints should return.
type APIError struct {
	Code      Code   `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
	Details   any    `json:"details,omitempty"`

	// HTTPStatus is not serialized — it drives the HTTP status code.
	HTTPStatus int `json:"-"`
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// envelope wraps the error for JSON serialization.
type envelope struct {
	Error *APIError `json:"error"`
}

// Write writes the error as a JSON response.
func (e *APIError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.HTTPStatus)
	//nolint:errcheck // best-effort write
	json.NewEncoder(w).Encode(envelope{Error: e})
}

// --- Constructors ---

// BadRequest creates a 400 Bad Request error.
func BadRequest(message string) *APIError {
	return &APIError{Code: CodeBadRequest, Message: message, HTTPStatus: http.StatusBadRequest}
}

// Unauthorized creates a 401 Unauthorized error.
func Unauthorized(message string) *APIError {
	return &APIError{Code: CodeUnauthorized, Message: message, HTTPStatus: http.StatusUnauthorized}
}

// Forbidden creates a 403 Forbidden error.
func Forbidden(message string) *APIError {
	return &APIError{Code: CodeForbidden, Message: message, HTTPStatus: http.StatusForbidden}
}

// ForbiddenWithDetails creates a 403 Forbidden error with extra details.
func ForbiddenWithDetails(message string, details any) *APIError {
	return &APIError{Code: CodeForbidden, Message: message, HTTPStatus: http.StatusForbidden, Details: details}
}

// NotFound creates a 404 Not Found error.
func NotFound(message string) *APIError {
	return &APIError{Code: CodeNotFound, Message: message, HTTPStatus: http.StatusNotFound}
}

// Conflict creates a 409 Conflict error.
func Conflict(message string) *APIError {
	return &APIError{Code: CodeConflict, Message: message, HTTPStatus: http.StatusConflict}
}

// RateLimited creates a 429 Too Many Requests error.
func RateLimited(retryAfter string) *APIError {
	return &APIError{
		Code:       CodeRateLimitExceeded,
		Message:    fmt.Sprintf("Rate limit exceeded. Retry after %s seconds.", retryAfter),
		HTTPStatus: http.StatusTooManyRequests,
	}
}

// ValidationFailed creates a 422 Unprocessable Entity error with field-level details.
func ValidationFailed(message string, fieldErrors map[string]string) *APIError {
	return &APIError{
		Code:       CodeValidationFailed,
		Message:    message,
		HTTPStatus: http.StatusUnprocessableEntity,
		Details:    fieldErrors,
	}
}

// Internal creates a 500 Internal Server Error.
func Internal(message string) *APIError {
	return &APIError{Code: CodeInternalError, Message: message, HTTPStatus: http.StatusInternalServerError}
}

// ServiceUnavailable creates a 503 Service Unavailable error.
func ServiceUnavailable(message string) *APIError {
	return &APIError{Code: CodeServiceUnavailable, Message: message, HTTPStatus: http.StatusServiceUnavailable}
}

// WithRequestID returns a copy of the error with a request ID attached.
func (e *APIError) WithRequestID(requestID string) *APIError {
	e.RequestID = requestID
	return e
}

// WithTraceID returns a copy of the error with a trace ID attached.
func (e *APIError) WithTraceID(traceID string) *APIError {
	e.TraceID = traceID
	return e
}
