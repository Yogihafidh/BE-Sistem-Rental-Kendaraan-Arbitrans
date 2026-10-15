package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, http.StatusOK, "data retrieved", map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success true, got %v", resp.Success)
	}
	if resp.Message != "data retrieved" {
		t.Errorf("expected message 'data retrieved', got '%s'", resp.Message)
	}
}

func TestErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	BadRequest(c, "invalid parameters", "field 'name' is required")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Success {
		t.Errorf("expected success false, got %v", resp.Success)
	}
	if resp.Message != "invalid parameters" {
		t.Errorf("expected message 'invalid parameters', got '%s'", resp.Message)
	}
}

func TestSuccessWithMeta(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	meta := &Meta{
		Page:       1,
		PerPage:    10,
		TotalRows:  50,
		TotalPages: 5,
	}

	SuccessWithMeta(c, http.StatusOK, "list retrieved", []string{"item1"}, meta)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Meta == nil {
		t.Fatal("expected meta to be non-nil")
	}
	if resp.Meta.TotalPages != 5 {
		t.Errorf("expected TotalPages 5, got %d", resp.Meta.TotalPages)
	}
}
