# Stage 1: Build the Go binary
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/calculator-server ./cmd/server

# Stage 2: Minimal non-root runtime image
FROM alpine:3.20

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/calculator-server /app/calculator-server

USER appuser

EXPOSE 8080

ENV PORT=8080

CMD ["/app/calculator-server"]
