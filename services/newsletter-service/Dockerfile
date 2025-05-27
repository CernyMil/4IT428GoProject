# Use the official Go image as the base image
FROM golang:1.24 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o newsletter-service ./cmd/api

# Use a minimal base image for the final container
FROM gcr.io/distroless/base-debian11

# Set the working directory inside the container
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /app/newsletter-service .

# Expose the port your service listens on
EXPOSE 8082

# Command to run the application
CMD ["./newsletter-service"]