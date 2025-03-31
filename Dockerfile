FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o data-collector-service .

# Create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/data-collector-service .
COPY --from=builder /app/migrations ./migrations

# Install PostgreSQL client for migrations
RUN apk add --no-cache postgresql-client

# Expose port
EXPOSE 8080

# Run the application
CMD ["./data-collector-service"]
