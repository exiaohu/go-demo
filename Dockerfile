# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache make git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build-linux

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (if any)
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/go-demo_unix ./go-demo
COPY --from=builder /app/config/config.yaml ./config.yaml

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./go-demo", "server"]
