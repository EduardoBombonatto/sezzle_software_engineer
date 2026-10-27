# AI Prompts Used During Development

In accordance with the assignment instructions (*"Share any prompts that you used in your work"*), this document records the prompts used to plan, design, test, and implement the Full-Stack Calculator application.

---

## Prompt 1 — Architecture & Incremental TDD Plan
> Build a full-stack calculator application with a React (TypeScript + Vite) frontend and a Go backend microservice consuming a REST API for basic and advanced arithmetic operations (`Addition`, `Subtraction`, `Multiplication`, `Division`, `Exponentiation`, `Square Root`, `Percentage`).
>
> Key engineering requirements:
> 1. **100% in English**: All code, comments, error messages, tests, API contracts, and documentation must be in English.
> 2. **Test-Driven Development (TDD)**: Follow Red → Green → Refactor for every operation across both backend and frontend.
> 3. **Clean, Readable, and Idiomatic Code**: Apply Clean Code principles (DRY, KISS, YAGNI, Guard Clauses, Single Level of Abstraction, no magic numbers) and SOLID principles (SRP, OCP, DIP).
> 4. **Exhaustive Pre-Calculation & Boundary Validation**: Check degenerate inputs without being asked—empty payloads, missing operands (`nil` vs `0` using `*float64`), single element on binary operations, unexpected second operand on unary operations, duplicate/equal elements, IEEE 754 `float64` extremes (`math.MaxFloat64`, `-math.MaxFloat64`, `SmallestNonzeroFloat64`), `-0.0` normalization, division by zero, negative square root, `0` raised to negative exponents, and negative bases with fractional exponents. Never execute the calculation if input validation fails.
> 5. **Standardized JSON Response Envelope**:
> ```json
> {
>   "erro": false,
>   "message": "Descriptive message",
>   "data": {},
>   "httpcode": 200,
>   "timestamp": "2026-07-16T22:40:00Z"
> }
> ```
> 6. **HTTP Client**: Use Axios exclusively in the frontend and validate inputs on the client before firing network requests.
> 7. **OpenAPI 3.0 Route**: Expose `/api/v1/openapi.json` and an interactive Swagger UI documentation page at `/api/v1/docs`.
> 8. **Incremental Git Workflow**: Implement and verify one operation per commit (`Addition` → `Subtraction` → `Multiplication` → `Division` → `Exponentiation` → `Square Root` → `Percentage` → `OpenAPI` → `Docker & Documentation`).

---

## Prompt 2 — UX Refinement for Error Notifications
> Move the validation and error warning messages out of the middle of the calculator card and display them above the calculator as a floating toast notification that automatically fades out and disappears after 4 seconds (4000ms), accompanied by automated timer tests.

---

## Prompt 3 — Step-by-Step Operation Implementation Prompts
> Implement each calculator operation incrementally using TDD (`Addition`, `Subtraction`, `Multiplication`, `Division`, `Exponentiation`, `Square Root`, `Percentage`), ensuring:
> - Backend table-driven unit tests (`internal/calculator`) and HTTP handler tests (`internal/api`) are written first.
> - Frontend validation unit tests (`validation.test.ts`), Axios client tests (`calculatorApi.test.ts`), and React Testing Library component tests (`Calculator.test.tsx`) are written and verified at every step.

---

## Prompt 4 — OpenAPI 3.0 & Docker Full-Stack Packaging
> Expose the OpenAPI 3.0.3 specification at `GET /api/v1/openapi.json` and interactive Swagger UI at `GET /api/v1/docs`. Then add multi-stage `Dockerfile`s for the Go backend and React/Nginx frontend, a root `docker-compose.yml` to run both services together, test coverage reports, and comprehensive `README.md` documentation.
