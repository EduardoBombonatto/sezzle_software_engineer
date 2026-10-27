package api

import (
	"errors"
	"net/http"

	"github.com/eduardobombonatto/sezzle-calculator/backend/internal/calculator"
)

const (
	OperationAdd      = "add"
	OperationSubtract = "subtract"
	OperationMultiply = "multiply"
	OperationDivide   = "divide"
	OperationPower    = "power"
	OperationSqrt     = "sqrt"
)

// Handler manages HTTP endpoints for the calculator microservice.
type Handler struct {
	calculator calculator.Calculator
}

// NewHandler creates a new HTTP Handler with the given Calculator implementation.
func NewHandler(calc calculator.Calculator) *Handler {
	return &Handler{
		calculator: calc,
	}
}

// RegisterRoutes configures all HTTP routes and wraps them with standard middleware.
func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", h.handleHealth)
	mux.HandleFunc("/api/v1/calculate/add", h.handleBinaryOperation(OperationAdd, h.calculator.Add))
	mux.HandleFunc("/api/v1/calculate/subtract", h.handleBinaryOperation(OperationSubtract, h.calculator.Subtract))
	mux.HandleFunc("/api/v1/calculate/multiply", h.handleBinaryOperation(OperationMultiply, h.calculator.Multiply))
	mux.HandleFunc("/api/v1/calculate/divide", h.handleBinaryOperation(OperationDivide, h.calculator.Divide))
	mux.HandleFunc("/api/v1/calculate/power", h.handleBinaryOperation(OperationPower, h.calculator.Power))
	mux.HandleFunc("/api/v1/calculate/sqrt", h.handleUnaryOperation(OperationSqrt, h.calculator.Sqrt))

	return withMiddleware(mux)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed; use GET")
		return
	}

	writeSuccess(w, http.StatusOK, "Service is healthy", map[string]string{
		"status": "healthy",
	})
}

func (h *Handler) handleBinaryOperation(
	operationName string,
	compute func(operandA, operandB float64) (float64, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operandA, operandB, valErr := parseBinaryOperands(w, r)
		if valErr != nil {
			writeError(w, valErr.statusCode, valErr.message)
			return
		}

		result, err := compute(operandA, operandB)
		if err != nil {
			statusCode, msg := mapDomainError(err)
			writeError(w, statusCode, msg)
			return
		}

		writeSuccess(w, http.StatusOK, "Operation completed successfully", CalculationResult{
			Operation: operationName,
			OperandA:  operandA,
			OperandB:  &operandB,
			Result:    result,
		})
	}
}

func (h *Handler) handleUnaryOperation(
	operationName string,
	compute func(operandA float64) (float64, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		operandA, valErr := parseUnaryOperand(w, r)
		if valErr != nil {
			writeError(w, valErr.statusCode, valErr.message)
			return
		}

		result, err := compute(operandA)
		if err != nil {
			statusCode, msg := mapDomainError(err)
			writeError(w, statusCode, msg)
			return
		}

		writeSuccess(w, http.StatusOK, "Operation completed successfully", CalculationResult{
			Operation: operationName,
			OperandA:  operandA,
			Result:    result,
		})
	}
}

func mapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, calculator.ErrInvalidOperand):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, calculator.ErrDivisionByZero),
		errors.Is(err, calculator.ErrZeroToNegativePower),
		errors.Is(err, calculator.ErrNegativeBaseFractionalExponent),
		errors.Is(err, calculator.ErrNegativeSquareRoot),
		errors.Is(err, calculator.ErrOverflow):
		return http.StatusUnprocessableEntity, err.Error()
	default:
		return http.StatusUnprocessableEntity, err.Error()
	}
}
