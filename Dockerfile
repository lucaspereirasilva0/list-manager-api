# Use a multi-stage build for a compact image

# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for static compilation, making the binary self-contained
# -o /app/list-manager-api specifies the output path and name of the binary
# ./cmd/api specifies the main package to build
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/list-manager-api ./cmd/api

# Stage 2: Create the final, compact image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/list-manager-api .

# Expose the port the application listens on
# Based on cmd/api/main.go, defaultPort is 8081
EXPOSE 8085

# Command to run the application
# The application reads MONGO_URI and MONGO_DB_NAME from environment variables
CMD ["./list-manager-api"] 