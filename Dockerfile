# ── Stage 1: builder ────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

# Install build dependencies (CGO requires gcc, SQLite requires headers)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Download dependencies first (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/api ./cmd/main.go

# ── Stage 2: runtime ────────────────────────────────────────────────────────────
FROM alpine:3.20 AS runtime

# Runtime dependencies
RUN apk add --no-cache sqlite-libs ca-certificates tzdata

# Create non-root user
RUN adduser -D appuser

WORKDIR /app

# Copy compiled binary from builder
COPY --from=builder /app/api .

# Data directory for SQLite database
RUN mkdir -p /app/data && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./api"]
