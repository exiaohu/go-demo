# Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
# CGO_ENABLED=0 for static binary (required for scratch/distroless)
# -ldflags="-s -w" to strip debug symbols and reduce size
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server cmd/main.go

# Create a minimal runtime image
# Use distroless/static for smallest size and best security
FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config/config.yaml .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./server", "server"]
