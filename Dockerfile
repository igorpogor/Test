# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Set environment variables for Go
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies with retry logic
RUN --mount=type=cache,target=/go/pkg/mod \
    sh -c 'go mod download || (sleep 5 && go mod download)'

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /subscription-service ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /subscription-service .

# Copy migrations
COPY migrations /app/migrations

# Install wait-for-it script for database health check
RUN apk add --no-cache bash curl netcat-openbsd
COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./subscription-service"]
