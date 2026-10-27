package api_test

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/eduardobombonatto/sezzle-calculator/backend/internal/api"
	"github.com/eduardobombonatto/sezzle-calculator/backend/internal/calculator"
)

type envelopeResponse struct {
	Erro      bool                   `json:"erro"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
	HTTPCode  int                    `json:"httpcode"`
	Timestamp string                 `json:"timestamp"`
}

func newTestServer() http.Handler {
	calcService := calculator.NewService()
	handler := api.NewHandler(calcService)
	return handler.RegisterRoutes()
}

func assertStandardEnvelope(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedErrorFlag bool) envelopeResponse {
	t.Helper()

	if rec.Code != expectedStatus {
		t.Fatalf("expected HTTP status %d, got %d (body: %s)", expectedStatus, rec.Code, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	var env envelopeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("failed to unmarshal response envelope: %v (body: %s)", err, rec.Body.String())
	}

	if env.Erro != expectedErrorFlag {
		t.Fatalf("expected envelope erro=%v, got %v", expectedErrorFlag, env.Erro)
	}

	if env.HTTPCode != expectedStatus {
		t.Fatalf("expected envelope httpcode=%d, got %d", expectedStatus, env.HTTPCode)
	}

	if strings.TrimSpace(env.Message) == "" {
		t.Fatalf("expected non-empty message in envelope")
	}

	if _, err := time.Parse(time.RFC3339, env.Timestamp); err != nil {
		t.Fatalf("expected valid RFC3339 timestamp, got %q: %v", env.Timestamp, err)
	}

	if expectedErrorFlag && env.Data != nil {
		t.Fatalf("expected data to be null on error, got %v", env.Data)
	}

	if !expectedErrorFlag && env.Data == nil {
		t.Fatalf("expected non-null data on success")
	}

	return env
}

func TestHealthEndpoint(t *testing.T) {
	t.Parallel()

	router := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	env := assertStandardEnvelope(t, rec, http.StatusOK, false)
	if env.Data["status"] != "healthy" {
		t.Fatalf("expected status 'healthy', got %v", env.Data["status"])
	}
}

func TestAddEndpoint(t *testing.T) {
	t.Parallel()

	router := newTestServer()

	testCases := []struct {
		name              string
		method            string
		contentType       string
		body              string
		expectedStatus    int
		expectedErrorFlag bool
		expectedResult    float64
		expectedMsgSubstr string
	}{
		{
			name:              "valid addition of positive numbers",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 10.5, "b": 4.5}`,
			expectedStatus:    http.StatusOK,
			expectedErrorFlag: false,
			expectedResult:    15.0,
		},
		{
			name:              "valid addition with zero operands",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 0, "b": 0}`,
			expectedStatus:    http.StatusOK,
			expectedErrorFlag: false,
			expectedResult:    0.0,
		},
		{
			name:              "valid addition with duplicate negative values",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": -25.5, "b": -25.5}`,
			expectedStatus:    http.StatusOK,
			expectedErrorFlag: false,
			expectedResult:    -51.0,
		},
		{
			name:              "degenerate empty body",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              ``,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Request body must not be empty",
		},
		{
			name:              "degenerate empty JSON object missing both operands",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Both operands 'a' and 'b' are required",
		},
		{
			name:              "degenerate single operand missing b",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 10}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Both operands 'a' and 'b' are required",
		},
		{
			name:              "degenerate single operand missing a",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"b": 10}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Both operands 'a' and 'b' are required",
		},
		{
			name:              "degenerate null operand",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": null, "b": 10}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Both operands 'a' and 'b' are required",
		},
		{
			name:              "invalid data type string operand",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": "ten", "b": 5}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Invalid JSON payload",
		},
		{
			name:              "unknown JSON field rejected",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 1, "b": 2, "c": 3}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Invalid JSON payload",
		},
		{
			name:              "trailing JSON tokens rejected",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 1, "b": 2}{"a": 3, "b": 4}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Request body must contain only a single JSON object",
		},
		{
			name:              "extreme value overflow rejected with 422",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 1.7976931348623157e+308, "b": 1.7976931348623157e+308}`,
			expectedStatus:    http.StatusUnprocessableEntity,
			expectedErrorFlag: true,
			expectedMsgSubstr: "overflow",
		},
		{
			name:              "wrong HTTP method GET rejected with 405",
			method:            http.MethodGet,
			contentType:       "application/json",
			body:              ``,
			expectedStatus:    http.StatusMethodNotAllowed,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Method not allowed",
		},
		{
			name:              "unsupported content type rejected with 415",
			method:            http.MethodPost,
			contentType:       "text/plain",
			body:              `{"a": 1, "b": 2}`,
			expectedStatus:    http.StatusUnsupportedMediaType,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Content-Type must be application/json",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, "/api/v1/calculate/add", bytes.NewBufferString(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			env := assertStandardEnvelope(t, rec, tc.expectedStatus, tc.expectedErrorFlag)

			if tc.expectedErrorFlag {
				if !strings.Contains(strings.ToLower(env.Message), strings.ToLower(tc.expectedMsgSubstr)) {
					t.Fatalf("expected error message containing %q, got %q", tc.expectedMsgSubstr, env.Message)
				}
				return
			}

			gotResult, ok := env.Data["result"].(float64)
			if !ok {
				t.Fatalf("expected numeric result in data, got %v", env.Data["result"])
			}
			if math.Abs(gotResult-tc.expectedResult) > 1e-9 {
				t.Fatalf("expected result %v, got %v", tc.expectedResult, gotResult)
			}
			if env.Data["operation"] != "add" {
				t.Fatalf("expected operation 'add', got %v", env.Data["operation"])
			}
		})
	}
}

func TestSubtractEndpoint(t *testing.T) {
	t.Parallel()

	router := newTestServer()

	testCases := []struct {
		name              string
		method            string
		contentType       string
		body              string
		expectedStatus    int
		expectedErrorFlag bool
		expectedResult    float64
		expectedMsgSubstr string
	}{
		{
			name:              "valid subtraction of positive numbers",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 20.5, "b": 5.25}`,
			expectedStatus:    http.StatusOK,
			expectedErrorFlag: false,
			expectedResult:    15.25,
		},
		{
			name:              "valid subtraction of equal elements yields zero",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 42, "b": 42}`,
			expectedStatus:    http.StatusOK,
			expectedErrorFlag: false,
			expectedResult:    0.0,
		},
		{
			name:              "degenerate single operand missing b",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 20}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Both operands 'a' and 'b' are required",
		},
		{
			name:              "extreme value overflow when subtracting negative MaxFloat64",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 1.7976931348623157e+308, "b": -1.7976931348623157e+308}`,
			expectedStatus:    http.StatusUnprocessableEntity,
			expectedErrorFlag: true,
			expectedMsgSubstr: "overflow",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, "/api/v1/calculate/subtract", bytes.NewBufferString(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			env := assertStandardEnvelope(t, rec, tc.expectedStatus, tc.expectedErrorFlag)

			if tc.expectedErrorFlag {
				if !strings.Contains(strings.ToLower(env.Message), strings.ToLower(tc.expectedMsgSubstr)) {
					t.Fatalf("expected error message containing %q, got %q", tc.expectedMsgSubstr, env.Message)
				}
				return
			}

			gotResult, ok := env.Data["result"].(float64)
			if !ok {
				t.Fatalf("expected numeric result in data, got %v", env.Data["result"])
			}
			if math.Abs(gotResult-tc.expectedResult) > 1e-9 {
				t.Fatalf("expected result %v, got %v", tc.expectedResult, gotResult)
			}
			if env.Data["operation"] != "subtract" {
				t.Fatalf("expected operation 'subtract', got %v", env.Data["operation"])
			}
		})
	}
}

func TestMultiplyEndpoint(t *testing.T) {
	t.Parallel()

	router := newTestServer()

	testCases := []struct {
		name              string
		method            string
		contentType       string
		body              string
		expectedStatus    int
		expectedErrorFlag bool
		expectedResult    float64
		expectedMsgSubstr string
	}{
		{
			name:              "valid multiplication of numbers",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 7.5, "b": 4}`,
			expectedStatus:    http.StatusOK,
			expectedErrorFlag: false,
			expectedResult:    30.0,
		},
		{
			name:              "multiplication of negative by zero returns normalized zero",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": -150.5, "b": 0}`,
			expectedStatus:    http.StatusOK,
			expectedErrorFlag: false,
			expectedResult:    0.0,
		},
		{
			name:              "degenerate single operand missing b",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 7}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Both operands 'a' and 'b' are required",
		},
		{
			name:              "extreme value overflow when multiplying MaxFloat64 by 2",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 1.7976931348623157e+308, "b": 2}`,
			expectedStatus:    http.StatusUnprocessableEntity,
			expectedErrorFlag: true,
			expectedMsgSubstr: "overflow",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, "/api/v1/calculate/multiply", bytes.NewBufferString(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			env := assertStandardEnvelope(t, rec, tc.expectedStatus, tc.expectedErrorFlag)

			if tc.expectedErrorFlag {
				if !strings.Contains(strings.ToLower(env.Message), strings.ToLower(tc.expectedMsgSubstr)) {
					t.Fatalf("expected error message containing %q, got %q", tc.expectedMsgSubstr, env.Message)
				}
				return
			}

			gotResult, ok := env.Data["result"].(float64)
			if !ok {
				t.Fatalf("expected numeric result in data, got %v", env.Data["result"])
			}
			if math.Abs(gotResult-tc.expectedResult) > 1e-9 {
				t.Fatalf("expected result %v, got %v", tc.expectedResult, gotResult)
			}
			if env.Data["operation"] != "multiply" {
				t.Fatalf("expected operation 'multiply', got %v", env.Data["operation"])
			}
		})
	}
}

func TestDivideEndpoint(t *testing.T) {
	t.Parallel()

	router := newTestServer()

	testCases := []struct {
		name              string
		method            string
		contentType       string
		body              string
		expectedStatus    int
		expectedErrorFlag bool
		expectedResult    float64
		expectedMsgSubstr string
	}{
		{
			name:              "valid division of numbers",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 45, "b": 6}`,
			expectedStatus:    http.StatusOK,
			expectedErrorFlag: false,
			expectedResult:    7.5,
		},
		{
			name:              "zero divided by negative returns normalized zero",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 0, "b": -15}`,
			expectedStatus:    http.StatusOK,
			expectedErrorFlag: false,
			expectedResult:    0.0,
		},
		{
			name:              "division by zero rejected with 422 and data null",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 10, "b": 0}`,
			expectedStatus:    http.StatusUnprocessableEntity,
			expectedErrorFlag: true,
			expectedMsgSubstr: "division by zero is not allowed",
		},
		{
			name:              "zero divided by zero rejected with 422",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 0, "b": 0}`,
			expectedStatus:    http.StatusUnprocessableEntity,
			expectedErrorFlag: true,
			expectedMsgSubstr: "division by zero is not allowed",
		},
		{
			name:              "degenerate single operand missing b",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 10}`,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorFlag: true,
			expectedMsgSubstr: "Both operands 'a' and 'b' are required",
		},
		{
			name:              "extreme value overflow when dividing MaxFloat64 by 0.5",
			method:            http.MethodPost,
			contentType:       "application/json",
			body:              `{"a": 1.7976931348623157e+308, "b": 0.5}`,
			expectedStatus:    http.StatusUnprocessableEntity,
			expectedErrorFlag: true,
			expectedMsgSubstr: "overflow",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, "/api/v1/calculate/divide", bytes.NewBufferString(tc.body))
			if tc.contentType != "" {
				req.Header.Set("Content-Type", tc.contentType)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			env := assertStandardEnvelope(t, rec, tc.expectedStatus, tc.expectedErrorFlag)

			if tc.expectedErrorFlag {
				if !strings.Contains(strings.ToLower(env.Message), strings.ToLower(tc.expectedMsgSubstr)) {
					t.Fatalf("expected error message containing %q, got %q", tc.expectedMsgSubstr, env.Message)
				}
				return
			}

			gotResult, ok := env.Data["result"].(float64)
			if !ok {
				t.Fatalf("expected numeric result in data, got %v", env.Data["result"])
			}
			if math.Abs(gotResult-tc.expectedResult) > 1e-9 {
				t.Fatalf("expected result %v, got %v", tc.expectedResult, gotResult)
			}
			if env.Data["operation"] != "divide" {
				t.Fatalf("expected operation 'divide', got %v", env.Data["operation"])
			}
		})
	}
}
