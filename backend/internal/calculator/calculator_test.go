package calculator_test

import (
	"errors"
	"math"
	"testing"

	"github.com/eduardobombonatto/sezzle-calculator/backend/internal/calculator"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	svc := calculator.NewService()

	testCases := []struct {
		name        string
		operandA    float64
		operandB    float64
		expected    float64
		expectedErr error
	}{
		{
			name:     "positive integers",
			operandA: 12,
			operandB: 8,
			expected: 20,
		},
		{
			name:     "negative numbers",
			operandA: -15.5,
			operandB: -4.5,
			expected: -20.0,
		},
		{
			name:     "mixed sign numbers",
			operandA: 100.25,
			operandB: -40.25,
			expected: 60.0,
		},
		{
			name:     "duplicate equal elements",
			operandA: 42.5,
			operandB: 42.5,
			expected: 85.0,
		},
		{
			name:     "both operands zero",
			operandA: 0,
			operandB: 0,
			expected: 0,
		},
		{
			name:     "negative zero normalization",
			operandA: math.Copysign(0, -1),
			operandB: math.Copysign(0, -1),
			expected: 0,
		},
		{
			name:     "opposite extreme float64 values cancel to zero",
			operandA: math.MaxFloat64,
			operandB: -math.MaxFloat64,
			expected: 0,
		},
		{
			name:     "smallest nonzero float64 values",
			operandA: math.SmallestNonzeroFloat64,
			operandB: math.SmallestNonzeroFloat64,
			expected: 2 * math.SmallestNonzeroFloat64,
		},
		{
			name:     "max float64 plus zero stays finite",
			operandA: math.MaxFloat64,
			operandB: 0,
			expected: math.MaxFloat64,
		},
		{
			name:        "positive float64 overflow returns ErrOverflow",
			operandA:    math.MaxFloat64,
			operandB:    math.MaxFloat64,
			expectedErr: calculator.ErrOverflow,
		},
		{
			name:        "negative float64 overflow returns ErrOverflow",
			operandA:    -math.MaxFloat64,
			operandB:    -math.MaxFloat64,
			expectedErr: calculator.ErrOverflow,
		},
		{
			name:        "NaN first operand rejected before calculation",
			operandA:    math.NaN(),
			operandB:    10,
			expectedErr: calculator.ErrInvalidOperand,
		},
		{
			name:        "NaN second operand rejected before calculation",
			operandA:    10,
			operandB:    math.NaN(),
			expectedErr: calculator.ErrInvalidOperand,
		},
		{
			name:        "positive infinity first operand rejected before calculation",
			operandA:    math.Inf(1),
			operandB:    1,
			expectedErr: calculator.ErrInvalidOperand,
		},
		{
			name:        "negative infinity second operand rejected before calculation",
			operandA:    1,
			operandB:    math.Inf(-1),
			expectedErr: calculator.ErrInvalidOperand,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := svc.Add(tc.operandA, tc.operandB)

			if tc.expectedErr != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if result != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, result)
			}

			if result == 0 && math.Signbit(result) {
				t.Fatalf("expected positive zero (0.0), got negative zero (-0.0)")
			}
		})
	}
}
