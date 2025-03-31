# Data Collector Service

A microservice for collecting data from external sources and sending it for further processing as part of a product demand assessment system.

## Overview

Data Collector Service is responsible for:
1. Accepting data source configurations via REST API
2. Collecting data from configured external sources
3. Storing raw data in PostgreSQL database
4. Forwarding raw data to RabbitMQ for further processing

## API Endpoints

### POST /collect-config

Configure the data source for collection.

**Request:**
```json
{
  "source_url": "https://api.example.com/data",
  "api_key": "your-api-key"
}
```

**Response:**
```json
{
  "status": "success",
  "errorCode": 0,
  "error": null
}
```

### POST /collect

Trigger data collection from the configured source.

**Request:** Empty body

**Response:**
```json
{
  "status": "success",
  "errorCode": 0,
  "error": null
}
```

Or in case of an error:
```json
{
  "status": "error",
  "errorCode": 1001,
  "error": "Failed to connect to data source"
}
```

## Data Flow

1. Client sends configuration via `/collect-config` endpoint
2. Client triggers data collection via `/collect` endpoint
3. Service makes HTTP GET request to the configured source URL with API key in Authorization header
4. Raw data is received in JSON format (e.g., `{"product_id": "123", "sales": 100, "date": "2024-10-01"}`)
5. Raw data is stored in PostgreSQL database
6. Raw data is sent to RabbitMQ queue "raw_data_queue" for processing by other services

## Environment Variables

- `RABBITMQ_URL`: RabbitMQ connection string (default: "amqp://guest:guest@rabbitmq:5672/")
- `POSTGRES_DSN`: PostgreSQL connection string (default: "postgres://user:password@postgres:5432/db")

## Running with Docker

```bash
# Build the Docker image
docker build -t data-collector-service .

# Run the container
docker run -p 8080:8080 \
  -e RABBITMQ_URL="amqp://guest:guest@rabbitmq:5672/" \
  -e POSTGRES_DSN="postgres://user:password@postgres:5432/db" \
  data-collector-service
```

## Development

### Prerequisites
- Go 1.21+
- PostgreSQL
- RabbitMQ

### Setup
1. Clone the repository
2. Install dependencies: `go mod download`
3. Run the service: `go run main.go`