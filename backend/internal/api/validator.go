package api

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"strings"
)

const (
	MaxRequestBodyBytes = 4096
)

// CalculationRequest uses pointer fields so omitted or null JSON fields
// are distinguishable from explicit 0.0 values.
type CalculationRequest struct {
	OperandA *float64 `json:"a"`
	OperandB *float64 `json:"b"`
}

type validationError struct {
	statusCode int
	message    string
}

func (e *validationError) Error() string {
	return e.message
}

func validatePostJSONRequest(w http.ResponseWriter, r *http.Request) (*CalculationRequest, *validationError) {
	if r.Method != http.MethodPost {
		return nil, &validationError{
			statusCode: http.StatusMethodNotAllowed,
			message:    "Method not allowed; use POST",
		}
	}

	contentType := strings.ToLower(strings.TrimSpace(r.Header.Get(ContentTypeHeader)))
	if !strings.HasPrefix(contentType, "application/json") {
		return nil, &validationError{
			statusCode: http.StatusUnsupportedMediaType,
			message:    "Content-Type must be application/json",
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodyBytes)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req CalculationRequest
	if err := decoder.Decode(&req); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, &validationError{
				statusCode: http.StatusBadRequest,
				message:    "Request body must not be empty",
			}
		}
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return nil, &validationError{
				statusCode: http.StatusRequestEntityTooLarge,
				message:    "Request body exceeds the maximum allowed size of 4096 bytes",
			}
		}
		return nil, &validationError{
			statusCode: http.StatusBadRequest,
			message:    "Invalid JSON payload: operands must be valid numbers and no unknown fields are allowed",
		}
	}

	var extraToken interface{}
	if err := decoder.Decode(&extraToken); !errors.Is(err, io.EOF) {
		return nil, &validationError{
			statusCode: http.StatusBadRequest,
			message:    "Request body must contain only a single JSON object",
		}
	}

	return &req, nil
}

func parseBinaryOperands(w http.ResponseWriter, r *http.Request) (float64, float64, *validationError) {
	req, valErr := validatePostJSONRequest(w, r)
	if valErr != nil {
		return 0, 0, valErr
	}

	if req.OperandA == nil || req.OperandB == nil {
		return 0, 0, &validationError{
			statusCode: http.StatusBadRequest,
			message:    "Both operands 'a' and 'b' are required for this operation",
		}
	}

	operandA := *req.OperandA
	operandB := *req.OperandB

	if math.IsNaN(operandA) || math.IsInf(operandA, 0) || math.IsNaN(operandB) || math.IsInf(operandB, 0) {
		return 0, 0, &validationError{
			statusCode: http.StatusBadRequest,
			message:    "Operands 'a' and 'b' must be finite numbers",
		}
	}

	return operandA, operandB, nil
}

func parseUnaryOperand(w http.ResponseWriter, r *http.Request) (float64, *validationError) {
	req, valErr := validatePostJSONRequest(w, r)
	if valErr != nil {
		return 0, valErr
	}

	if req.OperandA == nil {
		return 0, &validationError{
			statusCode: http.StatusBadRequest,
			message:    "Operand 'a' is required for this unary operation",
		}
	}

	if req.OperandB != nil {
		return 0, &validationError{
			statusCode: http.StatusBadRequest,
			message:    "Unary operation does not accept a second operand 'b'",
		}
	}

	operandA := *req.OperandA
	if math.IsNaN(operandA) || math.IsInf(operandA, 0) {
		return 0, &validationError{
			statusCode: http.StatusBadRequest,
			message:    "Operand 'a' must be a finite number",
		}
	}

	return operandA, nil
}
