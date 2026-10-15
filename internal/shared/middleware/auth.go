package middleware

import (
	"errors"
	"net/http"
	"strings"

	sharedjwt "github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/internal/shared/jwt"
	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// UserClaimsKey is the context key for authenticated user claims.
const UserClaimsKey = "user_claims"

// AuthRequired verifies the presence and validity of the Bearer JWT token in the Authorization header.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.AbortWithError(c, http.StatusUnauthorized, "Authorization header is required", nil)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.AbortWithError(c, http.StatusUnauthorized, "Authorization header must use Bearer scheme", nil)
			return
		}

		claims, err := sharedjwt.ParseToken(parts[1])
		if err != nil {
			msg := "Invalid token"
			if errors.Is(err, sharedjwt.ErrExpiredToken) {
				msg = "Token has expired"
			}
			response.AbortWithError(c, http.StatusUnauthorized, msg, nil)
			return
		}

		c.Set(UserClaimsKey, claims)
		c.Next()
	}
}

// RequireRoles restricts endpoint access to authenticated users having one of the specified roles.
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := CurrentUser(c)
		if !ok {
			response.AbortWithError(c, http.StatusUnauthorized, "Authentication required", nil)
			return
		}

		for _, role := range allowedRoles {
			if strings.EqualFold(claims.Role, role) {
				c.Next()
				return
			}
		}

		response.AbortWithError(c, http.StatusForbidden, "You do not have permission to access this resource", nil)
	}
}

// CurrentUser retrieves the authenticated user claims from the Gin request context.
func CurrentUser(c *gin.Context) (sharedjwt.Claims, bool) {
	val, ok := c.Get(UserClaimsKey)
	if !ok {
		return sharedjwt.Claims{}, false
	}
	claims, ok := val.(sharedjwt.Claims)
	return claims, ok
}
