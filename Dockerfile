# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Set environment variables for Go
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Tidy modules and build the application
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /subscription-service ./cmd/server/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /subscription-service .

# Copy .env file (optional, environment variables from docker-compose take precedence)
COPY .env .env

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./subscription-service"]