package api

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	ContentTypeHeader = "Content-Type"
	ContentTypeJSON   = "application/json; charset=utf-8"
)

// Envelope represents the standardized API JSON response contract.
type Envelope struct {
	Erro      bool        `json:"erro"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	HTTPCode  int         `json:"httpcode"`
	Timestamp string      `json:"timestamp"`
}

// CalculationResult represents the structured data payload returned on successful calculations.
type CalculationResult struct {
	Operation string   `json:"operation"`
	OperandA  float64  `json:"a"`
	OperandB  *float64 `json:"b,omitempty"`
	Result    float64  `json:"result"`
}

func writeSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	writeEnvelope(w, Envelope{
		Erro:      false,
		Message:   message,
		Data:      data,
		HTTPCode:  statusCode,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeEnvelope(w, Envelope{
		Erro:      true,
		Message:   message,
		Data:      nil,
		HTTPCode:  statusCode,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func writeEnvelope(w http.ResponseWriter, payload Envelope) {
	w.Header().Set(ContentTypeHeader, ContentTypeJSON)
	w.WriteHeader(payload.HTTPCode)
	_ = json.NewEncoder(w).Encode(payload)
}
