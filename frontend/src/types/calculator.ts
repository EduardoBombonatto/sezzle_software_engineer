export type OperationType = 'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'sqrt';

export interface ApiEnvelope<T> {
  erro: boolean;
  message: string;
  data: T | null;
  httpcode: number;
  timestamp: string;
}

export interface CalculationRequestPayload {
  a: number;
  b?: number;
}

export interface CalculationResponseData {
  operation: OperationType;
  a: number;
  b?: number;
  result: number;
}

export interface CalculationHistoryEntry {
  id: string;
  expression: string;
  result: number;
  operation: OperationType;
  timestamp: string;
}

export interface ValidationResult {
  isValid: boolean;
  errorMessage?: string;
  operandA?: number;
  operandB?: number;
}
