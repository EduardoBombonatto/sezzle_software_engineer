import { beforeEach, describe, expect, it, vi } from 'vitest';
import { calculateAddition } from './calculatorApi';
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
});
