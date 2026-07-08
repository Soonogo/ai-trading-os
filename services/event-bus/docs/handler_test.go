package docs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/docs/openapi.yaml", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "openapi: 3.0.3") {
		t.Fatalf("response missing openapi version: %s", w.Body.String())
	}

	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/swagger", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for swagger, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "AITOS API Docs") {
		t.Fatalf("swagger html missing title: %s", w.Body.String())
	}
}
