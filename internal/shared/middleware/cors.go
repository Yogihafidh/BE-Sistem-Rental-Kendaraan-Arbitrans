package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/configs"
	"github.com/gin-gonic/gin"
)

// CORS returns a Gin middleware handling Cross-Origin Resource Sharing securely.
// It explicitly enforces that wildcard '*' cannot be combined with credentials (AllowCredentials: true),
// preventing cross-origin credential theft and reflection vulnerabilities.
func CORS(cfg configs.CORSConfig) gin.HandlerFunc {
	origins := cfg.AllowedOrigins
	methods := strings.Join(cfg.AllowedMethods, ", ")
	headers := strings.Join(cfg.AllowedHeaders, ", ")
	exposeHeaders := strings.Join(cfg.ExposeHeaders, ", ")
	maxAge := strconv.Itoa(int(cfg.MaxAge.Seconds()))

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowOrigin := ""
		hasWildcard := false

		for _, o := range origins {
			if o == "*" {
				hasWildcard = true
			} else if origin != "" && strings.EqualFold(o, origin) {
				allowOrigin = origin
				break
			}
		}

		// When origin is not explicitly in whitelist, check if wildcard is permitted
		if allowOrigin == "" && hasWildcard {
			// Per W3C CORS specification: Wildcard '*' cannot be used with credentials.
			// When AllowCredentials is true, we must NOT reflect untrusted origins.
			if !cfg.AllowCredentials {
				allowOrigin = "*"
			}
		}

		if allowOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
			c.Header("Access-Control-Allow-Methods", methods)
			c.Header("Access-Control-Allow-Headers", headers)
			if exposeHeaders != "" {
				c.Header("Access-Control-Expose-Headers", exposeHeaders)
			}
			if cfg.AllowCredentials && allowOrigin != "*" {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			if maxAge != "" && maxAge != "0" {
				c.Header("Access-Control-Max-Age", maxAge)
			}
		}

		// Handle preflight OPTIONS request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
