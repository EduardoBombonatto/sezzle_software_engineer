import axios from 'axios';
import type {
  ApiEnvelope,
  CalculationRequestPayload,
  CalculationResponseData,
} from '../types/calculator';
import { normalizeZero } from '../utils/validation';
import { httpClient } from './httpClient';

function assertFiniteNumbers(...values: number[]): void {
  for (const val of values) {
    if (typeof val !== 'number' || !Number.isFinite(val)) {
      throw new Error('Operands must be valid finite numbers');
    }
  }
}

function extractErrorMessage(error: unknown): string {
  if (axios.isAxiosError<ApiEnvelope<null>>(error)) {
    const serverMessage = error.response?.data?.message;
    if (serverMessage) {
      return serverMessage;
    }
    return error.message || 'Network error while communicating with calculator service';
  }
  if (error instanceof Error) {
    return error.message;
  }
  return 'Unexpected error occurred';
}

async function postCalculation(
  endpoint: string,
  payload: CalculationRequestPayload
): Promise<CalculationResponseData> {
  try {
    const response = await httpClient.post<ApiEnvelope<CalculationResponseData>>(
      endpoint,
      payload
    );

    const envelope = response.data;
    if (envelope.erro || !envelope.data) {
      throw new Error(envelope.message || 'Calculation failed');
    }

    return {
      ...envelope.data,
      result: normalizeZero(envelope.data.result),
    };
  } catch (error) {
    throw new Error(extractErrorMessage(error));
  }
}

export async function calculateAddition(
  operandA: number,
  operandB: number
): Promise<CalculationResponseData> {
  assertFiniteNumbers(operandA, operandB);
  return postCalculation('/calculate/add', {
    a: normalizeZero(operandA),
    b: normalizeZero(operandB),
  });
}

export async function calculateSubtraction(
  operandA: number,
  operandB: number
): Promise<CalculationResponseData> {
  assertFiniteNumbers(operandA, operandB);
  return postCalculation('/calculate/subtract', {
    a: normalizeZero(operandA),
    b: normalizeZero(operandB),
  });
}

export async function calculateMultiplication(
  operandA: number,
  operandB: number
): Promise<CalculationResponseData> {
  assertFiniteNumbers(operandA, operandB);
  return postCalculation('/calculate/multiply', {
    a: normalizeZero(operandA),
    b: normalizeZero(operandB),
  });
}

export async function calculateDivision(
  operandA: number,
  operandB: number
): Promise<CalculationResponseData> {
  assertFiniteNumbers(operandA, operandB);
  if (operandB === 0) {
    throw new Error('Division by zero is not allowed');
  }
  return postCalculation('/calculate/divide', {
    a: normalizeZero(operandA),
    b: normalizeZero(operandB),
  });
}
