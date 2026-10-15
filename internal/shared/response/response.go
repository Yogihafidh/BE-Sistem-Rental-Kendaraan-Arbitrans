package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents the standard JSON API response structure.
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Meta holds pagination metadata for list endpoints.
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"per_page,omitempty"`
	TotalRows  int64 `json:"total_rows,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// Success sends a success response with status code, message, and payload data.
func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a paginated success response with metadata.
func SuccessWithMeta(c *gin.Context, statusCode int, message string, data any, meta *Meta) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends an error response with status code, message, and error details.
func Error(c *gin.Context, statusCode int, message string, errDetails any) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errDetails,
	})
}

// AbortWithError sends an error response and aborts the Gin middleware chain.
func AbortWithError(c *gin.Context, statusCode int, message string, errDetails any) {
	c.AbortWithStatusJSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errDetails,
	})
}

// BadRequest sends a 400 Bad Request error response.
func BadRequest(c *gin.Context, message string, errDetails any) {
	if message == "" {
		message = "Bad request"
	}
	Error(c, http.StatusBadRequest, message, errDetails)
}

// Unauthorized sends a 401 Unauthorized error response.
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden sends a 403 Forbidden error response.
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Access forbidden"
	}
	Error(c, http.StatusForbidden, message, nil)
}

// NotFound sends a 404 Not Found error response.
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	Error(c, http.StatusNotFound, message, nil)
}

// InternalServerError sends a 500 Internal Server Error response.
func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "An unexpected error occurred"
	}
	Error(c, http.StatusInternalServerError, message, nil)
}
