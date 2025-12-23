# Stage 1: build
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go module files
COPY go/go.mod go/go.sum ./
RUN go mod download

# Copy source
COPY go/ ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# Stage 2: runtime
# Default to production build (use Dockerfile.dev for development)
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY config/ /app/config/
COPY templates/ /app/templates/

# Cloud Run requires PORT env, default to 8080
ENV PORT=8080
ENV ENV=prod
ENV LOG_LEVEL=info

ENTRYPOINT ["/app/server"]

