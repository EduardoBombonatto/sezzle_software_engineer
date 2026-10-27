import { describe, expect, it } from 'vitest';
import { formatNumberForDisplay, validateBinaryOperation } from './validation';

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
});
