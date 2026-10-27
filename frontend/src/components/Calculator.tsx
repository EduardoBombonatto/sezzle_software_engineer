import type { FormEvent } from 'react';
import { useCalculator } from '../hooks/useCalculator';
import type { OperationType } from '../types/calculator';
import { formatNumberForDisplay, isUnaryOperation } from '../utils/validation';

interface OperationOption {
  type: OperationType;
  label: string;
  symbol: string;
}

const AVAILABLE_OPERATIONS: OperationOption[] = [
  { type: 'add', label: 'Add (+)', symbol: '+' },
  { type: 'subtract', label: 'Subtract (-)', symbol: '-' },
  { type: 'multiply', label: 'Multiply (×)', symbol: '×' },
  { type: 'divide', label: 'Divide (÷)', symbol: '÷' },
  { type: 'power', label: 'Power (xʸ)', symbol: '^' },
  { type: 'sqrt', label: 'Square Root (√x)', symbol: '√' },
];

export function Calculator() {
  const {
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
  } = useCalculator();

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    await executeCalculation();
  };

  return (
    <div className="calculator-page">
      <div className="notification-slot" aria-live="assertive">
        {errorMessage && (
          <div className="toast-notification" role="alert">
            <span className="toast-message">{errorMessage}</span>
            <button
              type="button"
              className="toast-close"
              onClick={dismissError}
              aria-label="Dismiss notification"
            >
              ×
            </button>
            <div className="toast-progress" />
          </div>
        )}
      </div>

      <main className="calculator-layout">
        <section className="calculator-card" aria-label="Calculator">
          <header className="calculator-header">
            <h1>Sezzle Calculator</h1>
            <p className="subtitle">Full-Stack Arithmetic Microservice Client</p>
          </header>

          <div className="display-panel" aria-live="polite">
            <div className="display-label">Result</div>
            <div className="display-value" data-testid="calculation-result">
              {result !== null ? result : '0'}
            </div>
          </div>

          <form className="calculator-form" onSubmit={handleSubmit} noValidate>
            <div className="input-group">
              <label htmlFor="operand-a">First Operand (a)</label>
              <input
                id="operand-a"
                type="text"
                inputMode="decimal"
                placeholder="Enter first number (e.g. 12.5)"
                value={operandA}
                onChange={(e) => setOperandA(e.target.value)}
                disabled={isLoading}
              />
            </div>

            <fieldset className="operations-fieldset">
              <legend>Operation</legend>
              <div className="operations-grid">
                {AVAILABLE_OPERATIONS.map((op) => (
                  <button
                    key={op.type}
                    type="button"
                    className={`op-button ${selectedOperation === op.type ? 'active' : ''}`}
                    onClick={() => setSelectedOperation(op.type)}
                    aria-pressed={selectedOperation === op.type}
                  >
                    {op.label}
                  </button>
                ))}
              </div>
            </fieldset>

            <div className="input-group">
              <label htmlFor="operand-b">Second Operand (b)</label>
              <input
                id="operand-b"
                type="text"
                inputMode="decimal"
                placeholder={
                  isUnaryOperation(selectedOperation)
                    ? 'Not required for Square Root (√x)'
                    : 'Enter second number (e.g. 4.5)'
                }
                value={isUnaryOperation(selectedOperation) ? '' : operandB}
                onChange={(e) => setOperandB(e.target.value)}
                disabled={isLoading || isUnaryOperation(selectedOperation)}
              />
            </div>

            <div className="action-buttons">
              <button
                type="button"
                className="secondary-button"
                onClick={clearAll}
                disabled={isLoading}
              >
                Clear
              </button>
              <button
                type="submit"
                className="primary-button"
                disabled={isLoading}
              >
                {isLoading ? 'Calculating...' : 'Calculate'}
              </button>
            </div>
          </form>
        </section>

        <aside className="history-card" aria-label="Calculation History">
          <h2>Recent Calculations</h2>
          {history.length === 0 ? (
            <p className="empty-history">No calculations yet.</p>
          ) : (
            <ul className="history-list" data-testid="history-list">
              {history.map((item) => (
                <li key={item.id} className="history-item">
                  <span className="history-expression">{item.expression}</span>
                  <span className="history-equals">{' = '}</span>
                  <strong className="history-result">
                    {formatNumberForDisplay(item.result)}
                  </strong>
                </li>
              ))}
            </ul>
          )}
        </aside>
      </main>
    </div>
  );
}
