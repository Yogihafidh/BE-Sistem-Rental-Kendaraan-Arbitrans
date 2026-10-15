package middleware

import (
	"net/http"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// BodySizeLimit limits the maximum allowed request body size in bytes.
// It guards against nil request bodies and prevents memory exhaustion from oversized payloads.
func BodySizeLimit(limitBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > limitBytes {
			response.AbortWithError(c, http.StatusRequestEntityTooLarge, "Request entity too large", nil)
			return
		}
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limitBytes)
		}
		c.Next()
	}
}
