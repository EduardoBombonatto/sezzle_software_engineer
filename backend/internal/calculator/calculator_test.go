package calculator_test

import (
	"errors"
	"math"
	"testing"

	"github.com/eduardobombonatto/sezzle-calculator/backend/internal/calculator"
)

type binaryTestCase struct {
	name        string
	operandA    float64
	operandB    float64
	expected    float64
	expectedErr error
}

func runBinaryTestCases(
	t *testing.T,
	operation func(operandA, operandB float64) (float64, error),
	testCases []binaryTestCase,
) {
	t.Helper()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := operation(tc.operandA, tc.operandB)

			if tc.expectedErr != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if math.Abs(result-tc.expected) > 1e-12 && result != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, result)
			}

			if result == 0 && math.Signbit(result) {
				t.Fatalf("expected positive zero (0.0), got negative zero (-0.0)")
			}
		})
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	svc := calculator.NewService()

	runBinaryTestCases(t, svc.Add, []binaryTestCase{
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
	})
}

func TestSubtract(t *testing.T) {
	t.Parallel()

	svc := calculator.NewService()

	runBinaryTestCases(t, svc.Subtract, []binaryTestCase{
		{
			name:     "positive integers",
			operandA: 25,
			operandB: 10,
			expected: 15,
		},
		{
			name:     "subtracting larger from smaller yields negative",
			operandA: 10.5,
			operandB: 25.5,
			expected: -15.0,
		},
		{
			name:     "negative numbers",
			operandA: -15.5,
			operandB: -5.5,
			expected: -10.0,
		},
		{
			name:     "duplicate equal elements cancel to zero",
			operandA: 99.95,
			operandB: 99.95,
			expected: 0,
		},
		{
			name:     "duplicate extreme MaxFloat64 values cancel to zero",
			operandA: math.MaxFloat64,
			operandB: math.MaxFloat64,
			expected: 0,
		},
		{
			name:     "both operands zero",
			operandA: 0,
			operandB: 0,
			expected: 0,
		},
		{
			name:     "negative zero minus positive zero normalizes to positive zero",
			operandA: math.Copysign(0, -1),
			operandB: 0,
			expected: 0,
		},
		{
			name:     "smallest nonzero float64 values cancel to zero",
			operandA: math.SmallestNonzeroFloat64,
			operandB: math.SmallestNonzeroFloat64,
			expected: 0,
		},
		{
			name:        "positive float64 overflow when subtracting negative MaxFloat64",
			operandA:    math.MaxFloat64,
			operandB:    -math.MaxFloat64,
			expectedErr: calculator.ErrOverflow,
		},
		{
			name:        "negative float64 overflow when subtracting positive MaxFloat64",
			operandA:    -math.MaxFloat64,
			operandB:    math.MaxFloat64,
			expectedErr: calculator.ErrOverflow,
		},
		{
			name:        "NaN first operand rejected before calculation",
			operandA:    math.NaN(),
			operandB:    5,
			expectedErr: calculator.ErrInvalidOperand,
		},
		{
			name:        "NaN second operand rejected before calculation",
			operandA:    5,
			operandB:    math.NaN(),
			expectedErr: calculator.ErrInvalidOperand,
		},
		{
			name:        "infinite operand rejected before calculation",
			operandA:    math.Inf(1),
			operandB:    5,
			expectedErr: calculator.ErrInvalidOperand,
		},
	})
}

func TestMultiply(t *testing.T) {
	t.Parallel()

	svc := calculator.NewService()

	runBinaryTestCases(t, svc.Multiply, []binaryTestCase{
		{
			name:     "positive integers",
			operandA: 6,
			operandB: 7,
			expected: 42,
		},
		{
			name:     "negative times negative yields positive",
			operandA: -8.5,
			operandB: -2.0,
			expected: 17.0,
		},
		{
			name:     "negative times positive yields negative",
			operandA: -12.5,
			operandB: 4.0,
			expected: -50.0,
		},
		{
			name:     "duplicate equal elements",
			operandA: 15,
			operandB: 15,
			expected: 225,
		},
		{
			name:     "both operands zero",
			operandA: 0,
			operandB: 0,
			expected: 0,
		},
		{
			name:     "negative number multiplied by zero normalizes -0.0 to 0.0",
			operandA: -999999.99,
			operandB: 0,
			expected: 0,
		},
		{
			name:     "negative MaxFloat64 multiplied by zero normalizes to 0.0",
			operandA: -math.MaxFloat64,
			operandB: 0,
			expected: 0,
		},
		{
			name:     "MaxFloat64 multiplied by identity 1 stays finite",
			operandA: math.MaxFloat64,
			operandB: 1,
			expected: math.MaxFloat64,
		},
		{
			name:        "positive float64 overflow when multiplying MaxFloat64 by 2",
			operandA:    math.MaxFloat64,
			operandB:    2,
			expectedErr: calculator.ErrOverflow,
		},
		{
			name:        "negative float64 overflow when multiplying MaxFloat64 by -2",
			operandA:    math.MaxFloat64,
			operandB:    -2,
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
			name:        "infinite operand rejected before calculation",
			operandA:    math.Inf(-1),
			operandB:    0,
			expectedErr: calculator.ErrInvalidOperand,
		},
	})
}

func TestDivide(t *testing.T) {
	t.Parallel()

	svc := calculator.NewService()

	runBinaryTestCases(t, svc.Divide, []binaryTestCase{
		{
			name:     "positive exact division",
			operandA: 42,
			operandB: 6,
			expected: 7,
		},
		{
			name:     "decimal quotient",
			operandA: 15,
			operandB: 4,
			expected: 3.75,
		},
		{
			name:     "negative divided by negative yields positive",
			operandA: -18,
			operandB: -3,
			expected: 6,
		},
		{
			name:     "negative divided by positive yields negative",
			operandA: -15,
			operandB: 2,
			expected: -7.5,
		},
		{
			name:     "duplicate non-zero elements yield 1",
			operandA: 77.7,
			operandB: 77.7,
			expected: 1,
		},
		{
			name:     "duplicate MaxFloat64 divided by MaxFloat64 yields 1",
			operandA: math.MaxFloat64,
			operandB: math.MaxFloat64,
			expected: 1,
		},
		{
			name:     "zero divided by positive number yields zero",
			operandA: 0,
			operandB: 42,
			expected: 0,
		},
		{
			name:     "zero divided by negative number normalizes -0.0 to 0.0",
			operandA: 0,
			operandB: -42,
			expected: 0,
		},
		{
			name:        "division by positive zero rejected before calculation",
			operandA:    100,
			operandB:    0,
			expectedErr: calculator.ErrDivisionByZero,
		},
		{
			name:        "division by negative zero rejected before calculation",
			operandA:    100,
			operandB:    math.Copysign(0, -1),
			expectedErr: calculator.ErrDivisionByZero,
		},
		{
			name:        "zero divided by zero rejected as division by zero",
			operandA:    0,
			operandB:    0,
			expectedErr: calculator.ErrDivisionByZero,
		},
		{
			name:        "division of MaxFloat64 by fractional number overflows",
			operandA:    math.MaxFloat64,
			operandB:    0.5,
			expectedErr: calculator.ErrOverflow,
		},
		{
			name:        "division of negative MaxFloat64 by fractional number overflows",
			operandA:    -math.MaxFloat64,
			operandB:    0.5,
			expectedErr: calculator.ErrOverflow,
		},
		{
			name:        "NaN first operand rejected before calculation",
			operandA:    math.NaN(),
			operandB:    2,
			expectedErr: calculator.ErrInvalidOperand,
		},
		{
			name:        "NaN second operand rejected before calculation",
			operandA:    2,
			operandB:    math.NaN(),
			expectedErr: calculator.ErrInvalidOperand,
		},
		{
			name:        "infinite operand rejected before calculation",
			operandA:    math.Inf(1),
			operandB:    2,
			expectedErr: calculator.ErrInvalidOperand,
		},
	})
}
