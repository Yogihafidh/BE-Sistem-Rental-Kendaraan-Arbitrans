package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/configs"
	sharedjwt "github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/internal/shared/jwt"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter(jwtCfg configs.JWTConfig) (*gin.Engine, string) {
	sharedjwt.Init(jwtCfg)

	token, _ := sharedjwt.GenerateToken("test-user-id", "admin@example.com", "admin")

	r := gin.New()
	r.Use(AuthRequired())

	r.GET("/protected", func(c *gin.Context) {
		claims, ok := CurrentUser(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no claims"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": claims.UserID, "role": claims.Role})
	})

	r.GET("/admin-only", RequireRoles("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/operator-only", RequireRoles("operator"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r, token
}

func TestAuthRequired_Success(t *testing.T) {
	cfg := configs.JWTConfig{
		Secret:             "test-secret-key-12345678901234567890",
		ExpirationDuration: time.Hour,
	}
	r, token := setupTestRouter(cfg)

	req, _ := http.NewRequest(http.MethodGet, "/protected", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d (%s)", w.Code, w.Body.String())
	}
}

func TestAuthRequired_MissingHeader(t *testing.T) {
	cfg := configs.JWTConfig{
		Secret:             "test-secret-key-12345678901234567890",
		ExpirationDuration: time.Hour,
	}
	r, _ := setupTestRouter(cfg)

	req, _ := http.NewRequest(http.MethodGet, "/protected", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got: %d", w.Code)
	}
}

func TestRequireRoles_Forbidden(t *testing.T) {
	cfg := configs.JWTConfig{
		Secret:             "test-secret-key-12345678901234567890",
		ExpirationDuration: time.Hour,
	}
	r, token := setupTestRouter(cfg)

	// User role is 'admin', requesting '/operator-only'
	req, _ := http.NewRequest(http.MethodGet, "/operator-only", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got: %d (%s)", w.Code, w.Body.String())
	}
}

func TestRequireRoles_Allowed(t *testing.T) {
	cfg := configs.JWTConfig{
		Secret:             "test-secret-key-12345678901234567890",
		ExpirationDuration: time.Hour,
	}
	r, token := setupTestRouter(cfg)

	// User role is 'admin', requesting '/admin-only'
	req, _ := http.NewRequest(http.MethodGet, "/admin-only", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d (%s)", w.Code, w.Body.String())
	}
}
