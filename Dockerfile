FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o data-collector-service .

# Create a minimal production image
FROM alpine:3.18

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/data-collector-service .

# Copy the dataset directory
COPY --from=builder /app/dataset ./dataset

# Create .env file from example if needed
COPY --from=builder /app/.env.example ./.env

# Set executable permissions
RUN chmod +x ./data-collector-service

# Run as non-root user for better security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Command to run the executable
CMD ["./data-collector-service"]
