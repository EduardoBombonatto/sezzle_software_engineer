package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenAPISpecEndpoint(t *testing.T) {
	t.Parallel()

	router := newTestServer()

	t.Run("GET /api/v1/openapi.json returns valid OpenAPI 3.0.3 specification", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/openapi.json", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Fatalf("expected application/json Content-Type, got %q", contentType)
		}

		var spec map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
			t.Fatalf("failed to unmarshal OpenAPI JSON: %v", err)
		}

		if spec["openapi"] != "3.0.3" {
			t.Fatalf("expected openapi version 3.0.3, got %v", spec["openapi"])
		}

		paths, ok := spec["paths"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected paths object in OpenAPI spec")
		}

		expectedPaths := []string{
			"/api/v1/health",
			"/api/v1/calculate/add",
			"/api/v1/calculate/subtract",
			"/api/v1/calculate/multiply",
			"/api/v1/calculate/divide",
			"/api/v1/calculate/power",
			"/api/v1/calculate/sqrt",
			"/api/v1/calculate/percentage",
		}

		for _, path := range expectedPaths {
			if _, exists := paths[path]; !exists {
				t.Fatalf("expected path %q to be documented in OpenAPI spec", path)
			}
		}

		components, ok := spec["components"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected components object in OpenAPI spec")
		}
		schemas, ok := components["schemas"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected schemas object in OpenAPI components")
		}
		for _, schemaName := range []string{
			"BinaryCalculationRequest",
			"UnaryCalculationRequest",
			"CalculationSuccessEnvelope",
			"ErrorEnvelope",
		} {
			if _, exists := schemas[schemaName]; !exists {
				t.Fatalf("expected schema %q in OpenAPI components", schemaName)
			}
		}
	})

	t.Run("POST /api/v1/openapi.json rejected with 405 standardized envelope", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/openapi.json", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assertStandardEnvelope(t, rec, http.StatusMethodNotAllowed, true)
	})
}

func TestOpenAPIDocsUIEndpoint(t *testing.T) {
	t.Parallel()

	router := newTestServer()

	t.Run("GET /api/v1/docs renders interactive Swagger UI HTML", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/docs", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rec.Code)
		}

		contentType := rec.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Fatalf("expected text/html Content-Type, got %q", contentType)
		}

		body := rec.Body.String()
		if !strings.Contains(body, "/api/v1/openapi.json") || !strings.Contains(body, "SwaggerUIBundle") {
			t.Fatalf("expected Swagger UI HTML referencing /api/v1/openapi.json")
		}
	})

	t.Run("POST /api/v1/docs rejected with 405 standardized envelope", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/docs", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assertStandardEnvelope(t, rec, http.StatusMethodNotAllowed, true)
	})
}
