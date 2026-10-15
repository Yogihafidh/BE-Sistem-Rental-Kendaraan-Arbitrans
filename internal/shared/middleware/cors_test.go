package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Yogihafidh/BE-Sistem-Rental-Kendaraan-Arbitrans/configs"
	"github.com/gin-gonic/gin"
)

func TestCORS_ExplicitOriginWithCredentials(t *testing.T) {
	cfg := configs.CORSConfig{
		AllowedOrigins:   []string{"https://arbitrans.com", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           time.Hour,
	}

	r := gin.New()
	r.Use(CORS(cfg))
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://arbitrans.com")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "https://arbitrans.com" {
		t.Errorf("expected origin 'https://arbitrans.com', got '%s'", w.Header().Get("Access-Control-Allow-Origin"))
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Errorf("expected credentials 'true', got '%s'", w.Header().Get("Access-Control-Allow-Credentials"))
	}
}

func TestCORS_WildcardWithoutCredentials(t *testing.T) {
	cfg := configs.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: false,
	}

	r := gin.New()
	r.Use(CORS(cfg))
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://any-domain.com")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected origin '*', got '%s'", w.Header().Get("Access-Control-Allow-Origin"))
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "" {
		t.Errorf("credentials must not be set with wildcard, got '%s'", w.Header().Get("Access-Control-Allow-Credentials"))
	}
}

func TestCORS_WildcardWithCredentials_RejectsUntrustedOrigin(t *testing.T) {
	// Security test: Wildcard + Credentials must NOT reflect untrusted origins
	cfg := configs.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	}

	r := gin.New()
	r.Use(CORS(cfg))
	r.GET("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://evil-attacker.com")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Origin header must NOT reflect evil-attacker.com with credentials!
	if w.Header().Get("Access-Control-Allow-Origin") == "https://evil-attacker.com" {
		t.Fatal("vulnerability: reflected untrusted origin with AllowCredentials!")
	}
}

func TestCORS_PreflightOptions(t *testing.T) {
	cfg := configs.CORSConfig{
		AllowedOrigins: []string{"https://arbitrans.com"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Authorization"},
	}

	r := gin.New()
	r.Use(CORS(cfg))
	r.POST("/test", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://arbitrans.com")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204 for OPTIONS preflight, got %d", w.Code)
	}
}
