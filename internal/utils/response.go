package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard envelope for all API responses.
type Response struct {
	Success   bool      `json:"success"`
	Data      any       `json:"data,omitempty"`
	Error     *APIError `json:"error,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
}

// APIError carries structured error information.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Meta holds pagination metadata.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse is the envelope used for paginated list endpoints.
type PaginatedResponse struct {
	Success   bool   `json:"success"`
	Data      any    `json:"data"`
	Meta      *Meta  `json:"meta"`
	RequestID string `json:"request_id,omitempty"`
}

// GetRequestID reads the request_id key set by the RequestID middleware.
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}

// OK sends a 200 response with the supplied data.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success:   true,
		Data:      data,
		RequestID: GetRequestID(c),
	})
}

// Created sends a 201 response with the supplied data.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Success:   true,
		Data:      data,
		RequestID: GetRequestID(c),
	})
}

// NoContent sends a 204 response with no body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a 400 response. An optional details value can be provided.
func BadRequest(c *gin.Context, message string, details ...any) {
	var det any
	if len(details) > 0 {
		det = details[0]
	}
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &APIError{
			Code:    "BAD_REQUEST",
			Message: message,
			Details: det,
		},
		RequestID: GetRequestID(c),
	})
}

// Unauthorized sends a 401 response.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error: &APIError{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
		RequestID: GetRequestID(c),
	})
}

// Forbidden sends a 403 response.
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Error: &APIError{
			Code:    "FORBIDDEN",
			Message: message,
		},
		RequestID: GetRequestID(c),
	})
}

// NotFound sends a 404 response.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error: &APIError{
			Code:    "NOT_FOUND",
			Message: message,
		},
		RequestID: GetRequestID(c),
	})
}

// Conflict sends a 409 response.
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Error: &APIError{
			Code:    "CONFLICT",
			Message: message,
		},
		RequestID: GetRequestID(c),
	})
}

// InternalError sends a 500 response.
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &APIError{
			Code:    "INTERNAL_ERROR",
			Message: message,
		},
		RequestID: GetRequestID(c),
	})
}

// TooManyRequests sends a 429 response.
func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Success: false,
		Error: &APIError{
			Code:    "TOO_MANY_REQUESTS",
			Message: message,
		},
		RequestID: GetRequestID(c),
	})
}

// Paginated sends a 200 response with data and pagination metadata.
func Paginated(c *gin.Context, data any, meta *Meta) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success:   true,
		Data:      data,
		Meta:      meta,
		RequestID: GetRequestID(c),
	})
}
