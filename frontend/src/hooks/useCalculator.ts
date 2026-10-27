import { useCallback, useEffect, useState } from 'react';
import {
  calculateAddition,
  calculateDivision,
  calculateMultiplication,
  calculatePower,
  calculateSquareRoot,
  calculateSubtraction,
} from '../api/calculatorApi';
import type {
  CalculationHistoryEntry,
  CalculationResponseData,
  OperationType,
} from '../types/calculator';
import {
  formatNumberForDisplay,
  isUnaryOperation,
  validateBinaryOperation,
  validateUnaryOperation,
} from '../utils/validation';

const MAX_HISTORY_ENTRIES = 10;
export const NOTIFICATION_DURATION_MS = 4000;

type BinaryOperationType = Exclude<OperationType, 'sqrt'>;

const OPERATION_SYMBOLS: Record<OperationType, string> = {
  add: '+',
  subtract: '-',
  multiply: '×',
  divide: '÷',
  power: '^',
  sqrt: '√',
};

const BINARY_OPERATION_EXECUTORS: Record<
  BinaryOperationType,
  (operandA: number, operandB: number) => Promise<CalculationResponseData>
> = {
  add: calculateAddition,
  subtract: calculateSubtraction,
  multiply: calculateMultiplication,
  divide: calculateDivision,
  power: calculatePower,
};

export function useCalculator() {
  const [operandA, setOperandA] = useState<string>('');
  const [operandB, setOperandB] = useState<string>('');
  const [selectedOperation, setSelectedOperation] = useState<OperationType>('add');
  const [result, setResult] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [history, setHistory] = useState<CalculationHistoryEntry[]>([]);

  useEffect(() => {
    if (!errorMessage) {
      return undefined;
    }

    const timerId = window.setTimeout(() => {
      setErrorMessage(null);
    }, NOTIFICATION_DURATION_MS);

    return () => {
      window.clearTimeout(timerId);
    };
  }, [errorMessage]);

  const dismissError = useCallback(() => {
    setErrorMessage(null);
  }, []);

  const clearAll = useCallback(() => {
    setOperandA('');
    setOperandB('');
    setResult(null);
    setErrorMessage(null);
  }, []);

  const executeCalculation = useCallback(async () => {
    setErrorMessage(null);

    if (isUnaryOperation(selectedOperation)) {
      const unaryValidation = validateUnaryOperation(operandA, selectedOperation);
      if (!unaryValidation.isValid || unaryValidation.operandA === undefined) {
        setResult(null);
        setErrorMessage(unaryValidation.errorMessage || 'Invalid calculator input');
        return;
      }

      setIsLoading(true);
      try {
        const data = await calculateSquareRoot(unaryValidation.operandA);
        const formattedResult = formatNumberForDisplay(data.result);
        setResult(formattedResult);

        const expression = `√(${formatNumberForDisplay(unaryValidation.operandA)})`;
        const entry: CalculationHistoryEntry = {
          id: `${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
          expression,
          result: data.result,
          operation: selectedOperation,
          timestamp: new Date().toISOString(),
        };
        setHistory((prev) => [entry, ...prev].slice(0, MAX_HISTORY_ENTRIES));
      } catch (error) {
        setResult(null);
        setErrorMessage(error instanceof Error ? error.message : 'Calculation failed');
      } finally {
        setIsLoading(false);
      }
      return;
    }

    const validation = validateBinaryOperation(operandA, operandB, selectedOperation);
    if (!validation.isValid || validation.operandA === undefined || validation.operandB === undefined) {
      setResult(null);
      setErrorMessage(validation.errorMessage || 'Invalid calculator input');
      return;
    }

    setIsLoading(true);
    try {
      const executor = BINARY_OPERATION_EXECUTORS[selectedOperation as BinaryOperationType];
      const data = await executor(validation.operandA, validation.operandB);
      const formattedResult = formatNumberForDisplay(data.result);
      setResult(formattedResult);

      const symbol = OPERATION_SYMBOLS[selectedOperation];
      const expression = `${formatNumberForDisplay(validation.operandA)} ${symbol} ${formatNumberForDisplay(validation.operandB)}`;
      const entry: CalculationHistoryEntry = {
        id: `${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
        expression,
        result: data.result,
        operation: selectedOperation,
        timestamp: new Date().toISOString(),
      };

      setHistory((prev) => [entry, ...prev].slice(0, MAX_HISTORY_ENTRIES));
    } catch (error) {
      setResult(null);
      setErrorMessage(error instanceof Error ? error.message : 'Calculation failed');
    } finally {
      setIsLoading(false);
    }
  }, [operandA, operandB, selectedOperation]);

  return {
    operandA,
    operandB,
    selectedOperation,
    result,
    errorMessage,
    isLoading,
    history,
    setOperandA,
    setOperandB,
    setSelectedOperation,
    dismissError,
    clearAll,
    executeCalculation,
  };
}
