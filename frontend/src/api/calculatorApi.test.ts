import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  calculateAddition,
  calculateDivision,
  calculateMultiplication,
  calculatePercentage,
  calculatePower,
  calculateSquareRoot,
  calculateSubtraction,
} from './calculatorApi';
import { httpClient } from './httpClient';

vi.mock('./httpClient', () => ({
  httpClient: {
    post: vi.fn(),
  },
}));

describe('calculatorApi - calculateAddition (Axios)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('sends POST request via Axios to /calculate/add and unwraps standardized envelope', async () => {
    vi.mocked(httpClient.post).mockResolvedValueOnce({
      data: {
        erro: false,
        message: 'Operation completed successfully',
        data: {
          operation: 'add',
          a: 10,
          b: 25,
          result: 35,
        },
        httpcode: 200,
        timestamp: '2026-10-26T21:00:00Z',
      },
    });

    const response = await calculateAddition(10, 25);

    expect(httpClient.post).toHaveBeenCalledTimes(1);
    expect(httpClient.post).toHaveBeenCalledWith('/calculate/add', { a: 10, b: 25 });
    expect(response).toEqual({
      operation: 'add',
      a: 10,
      b: 25,
      result: 35,
    });
  });

  it('does NOT call Axios if operands are non-finite (Fail-Fast validation)', async () => {
    await expect(calculateAddition(Number.NaN, 5)).rejects.toThrow(
      'Operands must be valid finite numbers'
    );
    expect(httpClient.post).not.toHaveBeenCalled();
  });

  it('surfaces backend standardized error message when erro is true', async () => {
    vi.mocked(httpClient.post).mockResolvedValueOnce({
      data: {
        erro: true,
        message: 'Arithmetic result overflowed maximum float64 bounds',
        data: null,
        httpcode: 422,
        timestamp: '2026-10-26T21:00:00Z',
      },
    });

    await expect(calculateAddition(1, 2)).rejects.toThrow(
      'Arithmetic result overflowed maximum float64 bounds'
    );
  });

  it('sends POST request via Axios to /calculate/subtract and unwraps standardized envelope', async () => {
    vi.mocked(httpClient.post).mockResolvedValueOnce({
      data: {
        erro: false,
        message: 'Operation completed successfully',
        data: {
          operation: 'subtract',
          a: 30,
          b: 12,
          result: 18,
        },
        httpcode: 200,
        timestamp: '2026-10-26T21:00:00Z',
      },
    });

    const response = await calculateSubtraction(30, 12);

    expect(httpClient.post).toHaveBeenCalledTimes(1);
    expect(httpClient.post).toHaveBeenCalledWith('/calculate/subtract', { a: 30, b: 12 });
    expect(response).toEqual({
      operation: 'subtract',
      a: 30,
      b: 12,
      result: 18,
    });
  });

  it('sends POST request via Axios to /calculate/multiply and unwraps standardized envelope', async () => {
    vi.mocked(httpClient.post).mockResolvedValueOnce({
      data: {
        erro: false,
        message: 'Operation completed successfully',
        data: {
          operation: 'multiply',
          a: 6,
          b: 7,
          result: 42,
        },
        httpcode: 200,
        timestamp: '2026-10-26T21:00:00Z',
      },
    });

    const response = await calculateMultiplication(6, 7);

    expect(httpClient.post).toHaveBeenCalledTimes(1);
    expect(httpClient.post).toHaveBeenCalledWith('/calculate/multiply', { a: 6, b: 7 });
    expect(response).toEqual({
      operation: 'multiply',
      a: 6,
      b: 7,
      result: 42,
    });
  });

  it('sends POST request via Axios to /calculate/divide and blocks division by zero before network call', async () => {
    await expect(calculateDivision(10, 0)).rejects.toThrow('Division by zero is not allowed');
    expect(httpClient.post).not.toHaveBeenCalled();

    vi.mocked(httpClient.post).mockResolvedValueOnce({
      data: {
        erro: false,
        message: 'Operation completed successfully',
        data: {
          operation: 'divide',
          a: 45,
          b: 6,
          result: 7.5,
        },
        httpcode: 200,
        timestamp: '2026-10-26T21:00:00Z',
      },
    });

    const response = await calculateDivision(45, 6);

    expect(httpClient.post).toHaveBeenCalledTimes(1);
    expect(httpClient.post).toHaveBeenCalledWith('/calculate/divide', { a: 45, b: 6 });
    expect(response).toEqual({
      operation: 'divide',
      a: 45,
      b: 6,
      result: 7.5,
    });
  });

  it('sends POST request via Axios to /calculate/power and blocks invalid domains before network call', async () => {
    await expect(calculatePower(0, -1)).rejects.toThrow(
      'Zero cannot be raised to a negative exponent'
    );
    await expect(calculatePower(-4, 0.5)).rejects.toThrow(
      'Negative base with fractional exponent is undefined in real numbers'
    );
    expect(httpClient.post).not.toHaveBeenCalled();

    vi.mocked(httpClient.post).mockResolvedValueOnce({
      data: {
        erro: false,
        message: 'Operation completed successfully',
        data: {
          operation: 'power',
          a: 2,
          b: 8,
          result: 256,
        },
        httpcode: 200,
        timestamp: '2026-10-26T21:00:00Z',
      },
    });

    const response = await calculatePower(2, 8);

    expect(httpClient.post).toHaveBeenCalledTimes(1);
    expect(httpClient.post).toHaveBeenCalledWith('/calculate/power', { a: 2, b: 8 });
    expect(response).toEqual({
      operation: 'power',
      a: 2,
      b: 8,
      result: 256,
    });
  });

  it('sends POST request via Axios to /calculate/sqrt and blocks negative numbers before network call', async () => {
    await expect(calculateSquareRoot(-9)).rejects.toThrow(
      'Cannot calculate square root of a negative number'
    );
    expect(httpClient.post).not.toHaveBeenCalled();

    vi.mocked(httpClient.post).mockResolvedValueOnce({
      data: {
        erro: false,
        message: 'Operation completed successfully',
        data: {
          operation: 'sqrt',
          a: 81,
          result: 9,
        },
        httpcode: 200,
        timestamp: '2026-10-26T21:00:00Z',
      },
    });

    const response = await calculateSquareRoot(81);

    expect(httpClient.post).toHaveBeenCalledTimes(1);
    expect(httpClient.post).toHaveBeenCalledWith('/calculate/sqrt', { a: 81 });
    expect(response).toEqual({
      operation: 'sqrt',
      a: 81,
      result: 9,
    });
  });

  it('sends POST request via Axios to /calculate/percentage and unwraps standardized envelope', async () => {
    vi.mocked(httpClient.post).mockResolvedValueOnce({
      data: {
        erro: false,
        message: 'Operation completed successfully',
        data: {
          operation: 'percentage',
          a: 15,
          b: 200,
          result: 30,
        },
        httpcode: 200,
        timestamp: '2026-10-26T21:00:00Z',
      },
    });

    const response = await calculatePercentage(15, 200);

    expect(httpClient.post).toHaveBeenCalledTimes(1);
    expect(httpClient.post).toHaveBeenCalledWith('/calculate/percentage', { a: 15, b: 200 });
    expect(response).toEqual({
      operation: 'percentage',
      a: 15,
      b: 200,
      result: 30,
    });
  });
});
