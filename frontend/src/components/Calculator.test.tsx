import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import * as calculatorApi from '../api/calculatorApi';
import { Calculator } from './Calculator';

vi.mock('../api/calculatorApi', () => ({
  calculateAddition: vi.fn(),
  calculateSubtraction: vi.fn(),
}));

describe('Calculator Component (Addition)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('renders calculator interface, inputs, and addition button', () => {
    render(<Calculator />);

    expect(screen.getByRole('heading', { name: /sezzle calculator/i })).toBeInTheDocument();
    expect(screen.getByLabelText(/first operand/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/second operand/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /add/i })).toBeInTheDocument();
  });

  it('calculates addition via Axios API and updates result and history', async () => {
    const user = userEvent.setup();
    vi.mocked(calculatorApi.calculateAddition).mockResolvedValueOnce({
      operation: 'add',
      a: 15.5,
      b: 4.5,
      result: 20,
    });

    render(<Calculator />);

    const inputA = screen.getByLabelText(/first operand/i);
    const inputB = screen.getByLabelText(/second operand/i);
    const calculateBtn = screen.getByRole('button', { name: /^calculate$/i });

    await user.clear(inputA);
    await user.type(inputA, '15.5');
    await user.clear(inputB);
    await user.type(inputB, '4.5');
    await user.click(calculateBtn);

    await waitFor(() => {
      expect(calculatorApi.calculateAddition).toHaveBeenCalledWith(15.5, 4.5);
    });

    expect(await screen.findByTestId('calculation-result')).toHaveTextContent('20');
    expect(screen.getByTestId('history-list')).toHaveTextContent('15.5 + 4.5 = 20');
  });

  it('shows notification above the calculator when operand is empty and auto-dismisses after 4 seconds', async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    const user = userEvent.setup({ advanceTimers: vi.advanceTimersByTime });
    render(<Calculator />);

    const calculateBtn = screen.getByRole('button', { name: /^calculate$/i });
    await user.click(calculateBtn);

    const alert = await screen.findByRole('alert');
    expect(alert).toHaveTextContent(/first operand is required/i);
    expect(calculatorApi.calculateAddition).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(4000);
    });

    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });

  it('displays backend error notification when API returns an error', async () => {
    const user = userEvent.setup();
    vi.mocked(calculatorApi.calculateAddition).mockRejectedValueOnce(
      new Error('Arithmetic result overflowed maximum float64 bounds')
    );

    render(<Calculator />);

    await user.type(screen.getByLabelText(/first operand/i), '10');
    await user.type(screen.getByLabelText(/second operand/i), '20');
    await user.click(screen.getByRole('button', { name: /^calculate$/i }));

    expect(await screen.findByRole('alert')).toHaveTextContent(
      /arithmetic result overflowed maximum float64 bounds/i
    );
  });

  it('calculates subtraction via Axios API when Subtract (-) is selected', async () => {
    const user = userEvent.setup();
    vi.mocked(calculatorApi.calculateSubtraction).mockResolvedValueOnce({
      operation: 'subtract',
      a: 50,
      b: 18,
      result: 32,
    });

    render(<Calculator />);

    await user.click(screen.getByRole('button', { name: /subtract/i }));
    await user.type(screen.getByLabelText(/first operand/i), '50');
    await user.type(screen.getByLabelText(/second operand/i), '18');
    await user.click(screen.getByRole('button', { name: /^calculate$/i }));

    await waitFor(() => {
      expect(calculatorApi.calculateSubtraction).toHaveBeenCalledWith(50, 18);
    });

    expect(await screen.findByTestId('calculation-result')).toHaveTextContent('32');
    expect(screen.getByTestId('history-list')).toHaveTextContent('50 - 18 = 32');
  });
});
