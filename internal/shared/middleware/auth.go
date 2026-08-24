package middleware

import (
	"net/http"
	"strings"

	sharedjwt "github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/internal/shared/jwt"
	"github.com/gin-gonic/gin"
)

const UserClaimsKey = "user_claims"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authorization header is required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authorization header must use bearer token"})
			return
		}

		claims, err := sharedjwt.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		c.Set(UserClaimsKey, claims)
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (sharedjwt.Claims, bool) {
	value, ok := c.Get(UserClaimsKey)
	if !ok {
		return sharedjwt.Claims{}, false
	}

	claims, ok := value.(sharedjwt.Claims)
	return claims, ok
}
