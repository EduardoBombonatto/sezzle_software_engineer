package calculator

import (
	"errors"
	"math"
)

var (
	// ErrInvalidOperand indicates that one or more input operands are NaN or Infinite.
	ErrInvalidOperand = errors.New("operands must be finite numbers")

	// ErrOverflow indicates that the arithmetic result exceeded IEEE 754 float64 bounds.
	ErrOverflow = errors.New("arithmetic result overflowed maximum float64 bounds")

	// ErrDivisionByZero indicates an attempt to divide a number by zero.
	ErrDivisionByZero = errors.New("division by zero is not allowed")

	// ErrZeroToNegativePower indicates an attempt to raise zero to a negative exponent.
	ErrZeroToNegativePower = errors.New("zero cannot be raised to a negative exponent")

	// ErrNegativeBaseFractionalExponent indicates an attempt to raise a negative base to a non-integer exponent.
	ErrNegativeBaseFractionalExponent = errors.New("negative base with fractional exponent is undefined in real numbers")
)

// Calculator defines the contract for arithmetic operations.
type Calculator interface {
	Add(operandA, operandB float64) (float64, error)
	Subtract(operandA, operandB float64) (float64, error)
	Multiply(operandA, operandB float64) (float64, error)
	Divide(operandA, operandB float64) (float64, error)
	Power(base, exponent float64) (float64, error)
}

// Service implements the Calculator domain logic with O(1) time and O(1) space complexity.
type Service struct{}

// NewService creates a new Calculator domain service instance.
func NewService() *Service {
	return &Service{}
}

// Add computes the sum of two finite float64 numbers after validating preconditions
// and verifying that the result does not overflow IEEE 754 float64 limits.
func (s *Service) Add(operandA, operandB float64) (float64, error) {
	if err := validateFiniteOperands(operandA, operandB); err != nil {
		return 0, err
	}

	return finalizeResult(operandA + operandB)
}

// Subtract computes the difference (operandA - operandB) of two finite float64 numbers
// after validating preconditions and checking for IEEE 754 float64 overflow.
func (s *Service) Subtract(operandA, operandB float64) (float64, error) {
	if err := validateFiniteOperands(operandA, operandB); err != nil {
		return 0, err
	}

	return finalizeResult(operandA - operandB)
}

// Multiply computes the product (operandA * operandB) of two finite float64 numbers
// after validating preconditions and checking for IEEE 754 float64 overflow.
func (s *Service) Multiply(operandA, operandB float64) (float64, error) {
	if err := validateFiniteOperands(operandA, operandB); err != nil {
		return 0, err
	}

	return finalizeResult(operandA * operandB)
}

// Divide computes the quotient (operandA / operandB) of two finite float64 numbers.
// It fails fast before executing division if operandB is zero (+0.0 or -0.0).
func (s *Service) Divide(operandA, operandB float64) (float64, error) {
	if err := validateFiniteOperands(operandA, operandB); err != nil {
		return 0, err
	}

	if operandB == 0 {
		return 0, ErrDivisionByZero
	}

	return finalizeResult(operandA / operandB)
}

// Power computes base raised to exponent (base^exponent).
// It validates domain preconditions (zero to negative power and negative base with fractional exponent)
// before invoking math.Pow.
func (s *Service) Power(base, exponent float64) (float64, error) {
	if err := validateFiniteOperands(base, exponent); err != nil {
		return 0, err
	}

	if base == 0 && exponent < 0 {
		return 0, ErrZeroToNegativePower
	}

	if base < 0 && math.Trunc(exponent) != exponent {
		return 0, ErrNegativeBaseFractionalExponent
	}

	return finalizeResult(math.Pow(base, exponent))
}

func validateFiniteOperands(operands ...float64) error {
	for _, value := range operands {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return ErrInvalidOperand
		}
	}
	return nil
}

func finalizeResult(result float64) (float64, error) {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrOverflow
	}

	return normalizeNegativeZero(result), nil
}

// normalizeNegativeZero converts IEEE 754 -0.0 into standard positive 0.0
// to guarantee consistent JSON serialization and equality comparisons.
func normalizeNegativeZero(value float64) float64 {
	if value == 0 {
		return 0
	}
	return value
}
