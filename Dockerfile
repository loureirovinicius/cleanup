# Use the official Golang image
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to cache dependencies
COPY go.mod go.sum ./

# Download dependencies only (this is meant to cache dependencies)
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o cleanup ./cmd/main.go

# Build minimal image
FROM alpine:3.18 AS official

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/cleanup /app/cleanup

# Set the entrypoint
ENTRYPOINT ["sh"]
