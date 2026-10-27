import { describe, expect, it } from 'vitest';
import {
  formatNumberForDisplay,
  isUnaryOperation,
  validateBinaryOperation,
  validateUnaryOperation,
} from './validation';

describe('Frontend Calculator Validation (Addition)', () => {
  it('validates valid positive, negative, and decimal numbers', () => {
    const res = validateBinaryOperation('12.5', '-2.5', 'add');
    expect(res.isValid).toBe(true);
    expect(res.operandA).toBe(12.5);
    expect(res.operandB).toBe(-2.5);
  });

  it('validates duplicate equal numbers and zeroes', () => {
    const zeroRes = validateBinaryOperation('0', '-0', 'add');
    expect(zeroRes.isValid).toBe(true);
    expect(Object.is(zeroRes.operandA, 0)).toBe(true);
    expect(Object.is(zeroRes.operandB, 0)).toBe(true);

    const dupRes = validateBinaryOperation('42', '42', 'add');
    expect(dupRes.isValid).toBe(true);
    expect(dupRes.operandA).toBe(42);
    expect(dupRes.operandB).toBe(42);
  });

  it('rejects degenerate empty inputs', () => {
    expect(validateBinaryOperation('', '5', 'add')).toEqual({
      isValid: false,
      errorMessage: 'First operand is required',
    });
    expect(validateBinaryOperation('5', '   ', 'add')).toEqual({
      isValid: false,
      errorMessage: 'Second operand is required',
    });
  });

  it('rejects non-numeric strings and malformed decimals', () => {
    const res1 = validateBinaryOperation('abc', '5', 'add');
    expect(res1.isValid).toBe(false);
    expect(res1.errorMessage).toContain('valid finite number');

    const res2 = validateBinaryOperation('5', '12..3', 'add');
    expect(res2.isValid).toBe(false);
    expect(res2.errorMessage).toContain('valid finite number');
  });

  it('rejects extreme inputs exceeding IEEE 754 float64 limits before calling API', () => {
    const res = validateBinaryOperation('1e309', '1', 'add');
    expect(res.isValid).toBe(false);
    expect(res.errorMessage).toContain('valid finite number');
  });

  it('rejects addition that would overflow IEEE 754 float64 limits before calling API', () => {
    const maxFloat = String(Number.MAX_VALUE);
    const res = validateBinaryOperation(maxFloat, maxFloat, 'add');
    expect(res.isValid).toBe(false);
    expect(res.errorMessage).toContain('overflow');
  });

  it('normalizes negative zero in formatNumberForDisplay', () => {
    expect(formatNumberForDisplay(-0)).toBe('0');
    expect(formatNumberForDisplay(15.25)).toBe('15.25');
  });

  it('validates subtraction inputs and detects float64 subtraction overflow', () => {
    const validRes = validateBinaryOperation('25.5', '10.25', 'subtract');
    expect(validRes.isValid).toBe(true);
    expect(validRes.operandA).toBe(25.5);
    expect(validRes.operandB).toBe(10.25);

    const maxFloat = String(Number.MAX_VALUE);
    const overflowRes = validateBinaryOperation(maxFloat, `-${maxFloat}`, 'subtract');
    expect(overflowRes.isValid).toBe(false);
    expect(overflowRes.errorMessage).toContain('overflow');
  });

  it('validates multiplication inputs and detects float64 multiplication overflow', () => {
    const validRes = validateBinaryOperation('-7.5', '4', 'multiply');
    expect(validRes.isValid).toBe(true);
    expect(validRes.operandA).toBe(-7.5);
    expect(validRes.operandB).toBe(4);

    const maxFloat = String(Number.MAX_VALUE);
    const overflowRes = validateBinaryOperation(maxFloat, '2', 'multiply');
    expect(overflowRes.isValid).toBe(false);
    expect(overflowRes.errorMessage).toContain('overflow');
  });

  it('validates division inputs, rejects division by zero (+0 and -0), and detects overflow', () => {
    const validRes = validateBinaryOperation('45', '6', 'divide');
    expect(validRes.isValid).toBe(true);
    expect(validRes.operandA).toBe(45);
    expect(validRes.operandB).toBe(6);

    const divByZero = validateBinaryOperation('10', '0', 'divide');
    expect(divByZero).toEqual({
      isValid: false,
      errorMessage: 'Division by zero is not allowed',
    });

    const divByNegZero = validateBinaryOperation('10', '-0.0', 'divide');
    expect(divByNegZero).toEqual({
      isValid: false,
      errorMessage: 'Division by zero is not allowed',
    });

    const maxFloat = String(Number.MAX_VALUE);
    const overflowRes = validateBinaryOperation(maxFloat, '0.5', 'divide');
    expect(overflowRes.isValid).toBe(false);
    expect(overflowRes.errorMessage).toContain('overflow');
  });

  it('validates exponentiation inputs and blocks invalid domains before network call', () => {
    const validRes = validateBinaryOperation('2', '10', 'power');
    expect(validRes.isValid).toBe(true);
    expect(validRes.operandA).toBe(2);
    expect(validRes.operandB).toBe(10);

    const zeroNeg = validateBinaryOperation('0', '-2', 'power');
    expect(zeroNeg).toEqual({
      isValid: false,
      errorMessage: 'Zero cannot be raised to a negative exponent',
    });

    const negBaseFrac = validateBinaryOperation('-9', '0.5', 'power');
    expect(negBaseFrac).toEqual({
      isValid: false,
      errorMessage: 'Negative base with fractional exponent is undefined in real numbers',
    });

    const overflowRes = validateBinaryOperation('10', '1000', 'power');
    expect(overflowRes.isValid).toBe(false);
    expect(overflowRes.errorMessage).toContain('overflow');
  });

  it('validates unary square root operation and rejects negative inputs before network call', () => {
    expect(isUnaryOperation('sqrt')).toBe(true);
    expect(isUnaryOperation('add')).toBe(false);

    const validZero = validateUnaryOperation('-0.0', 'sqrt');
    expect(validZero.isValid).toBe(true);
    expect(Object.is(validZero.operandA, 0)).toBe(true);

    const validSquare = validateUnaryOperation('81', 'sqrt');
    expect(validSquare).toEqual({
      isValid: true,
      operandA: 81,
    });

    const negativeInput = validateUnaryOperation('-16', 'sqrt');
    expect(negativeInput).toEqual({
      isValid: false,
      errorMessage: 'Cannot calculate square root of a negative number',
    });

    const emptyInput = validateUnaryOperation('   ', 'sqrt');
    expect(emptyInput).toEqual({
      isValid: false,
      errorMessage: 'First operand is required',
    });
  });
});
