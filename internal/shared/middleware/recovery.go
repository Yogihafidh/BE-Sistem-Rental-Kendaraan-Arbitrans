package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// Recovery returns a Gin middleware that recovers from panics, logs the stack trace with slog,
// and returns a standardized 500 Internal Server Error response.
func Recovery(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				errMsg := fmt.Sprintf("%v", r)

				if l != nil {
					l.Error("panic recovered in HTTP handler",
						slog.String("error", errMsg),
						slog.String("path", c.Request.URL.Path),
						slog.String("method", c.Request.Method),
						slog.String("stack", stack),
					)
				}

				response.AbortWithError(c, http.StatusInternalServerError, "An unexpected server error occurred", nil)
			}
		}()

		c.Next()
	}
}
