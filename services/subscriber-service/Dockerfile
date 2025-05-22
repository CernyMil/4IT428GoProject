# Start from the official Golang image for building
FROM golang:1.21.6-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/subscriber-service/main/api ./cmd/main.go

# Build the Go app
RUN go build -o subscriber-service .

# Use a minimal image for running
FROM alpine:latest

# Copy the built binary from builder
COPY --from=builder /app/subscriber-api .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api"]
