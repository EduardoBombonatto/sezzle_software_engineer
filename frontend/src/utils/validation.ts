import type { OperationType, ValidationResult } from '../types/calculator';

const NUMERIC_PATTERN = /^[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d+)?$/;

export function normalizeZero(value: number): number {
  return value === 0 ? 0 : value;
}

export function parseStrictFiniteNumber(raw: string): number | null {
  const trimmed = raw.trim();
  if (trimmed.length === 0 || !NUMERIC_PATTERN.test(trimmed)) {
    return null;
  }

  const parsed = Number(trimmed);
  if (!Number.isFinite(parsed)) {
    return null;
  }

  return normalizeZero(parsed);
}

export function validateBinaryOperation(
  rawOperandA: string,
  rawOperandB: string,
  operation: OperationType
): ValidationResult {
  if (rawOperandA.trim().length === 0) {
    return { isValid: false, errorMessage: 'First operand is required' };
  }

  if (rawOperandB.trim().length === 0) {
    return { isValid: false, errorMessage: 'Second operand is required' };
  }

  const operandA = parseStrictFiniteNumber(rawOperandA);
  if (operandA === null) {
    return {
      isValid: false,
      errorMessage: 'First operand must be a valid finite number within float64 bounds',
    };
  }

  const operandB = parseStrictFiniteNumber(rawOperandB);
  if (operandB === null) {
    return {
      isValid: false,
      errorMessage: 'Second operand must be a valid finite number within float64 bounds',
    };
  }

  if (operation === 'add') {
    const previewSum = operandA + operandB;
    if (!Number.isFinite(previewSum)) {
      return {
        isValid: false,
        errorMessage: 'Arithmetic result overflowed maximum float64 bounds',
      };
    }
  }

  return {
    isValid: true,
    operandA,
    operandB,
  };
}

export function formatNumberForDisplay(value: number): string {
  const normalized = normalizeZero(value);
  return String(normalized);
}
