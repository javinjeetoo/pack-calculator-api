# Pack Calculator API

A Go-based HTTP API that calculates the optimal pack allocation to fulfil an order, ensuring:

- only whole packs are shipped,
- the minimum number of items are sent to satisfy the order,
- and, within that constraint, the fewest number of packs are used.

This service is designed to be simple, testable, and production-oriented.

## Problem Summary

Given a set of available pack sizes and a requested number of items, the API determines the optimal combination of packs such that:

1. Packs cannot be broken.
2. No more items than necessary are shipped to fulfil the order.
3. Among valid solutions, the fewest number of packs is used.

Example pack sizes:

- 250
- 500
- 1000
- 2000
- 5000

---

## API Endpoints

### Health Check

`GET /health`

Response:

```json
{ "status": "ok" }
```

### Calculate Pack Allocation

`POST /calculate`

#### Request

```json
{
  "items": 12001
}
```

Optional override of pack sizes:

```json
{
  "items": 12001,
  "pack_sizes": [250, 500, 1000, 2000, 5000]
}
```

#### Response

```json
{
  "items_ordered": 12001,
  "items_shipped": 12250,
  "packs": {
    "5000": 2,
    "2000": 1,
    "250": 1
  },
  "total_packs": 4
}
```

## Algorithm Overview

The solution uses a dynamic programming approach to calculate the minimum number of packs required to reach each achievable item total.

The algorithm searches from the requested item count upwards to find the smallest valid shipped total, ensuring minimal over-shipment. For that shipped total, it selects the configuration with the fewest packs.

This guarantees correctness while remaining efficient for the expected input sizes.

## Running locally

### Prerequisites

- Go 1.22+

### Start the server

`go run ./cmd/server`

The API will be available at:
`http://localhost:8080`

### Testing

Run all unit tests with:
`go test ./...`

### Example Usage

```json
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"items":12001}'
```

### Configuration

Default pack sizes can be configured using an environment variable:

`export PACK_SIZES=250,500,1000,2000,5000`

If not set, the application uses sensible defaults.

## Running with Docker (Optional)

Build the Docker image:
`docker build -t pack-calculator-api .`

Run the container:
`docker run -p 8080:8080 pack-calculator-api`

The API will be available at:
`http://localhost:8080`

## Project Structure

```json
cmd/server/        # Application entry point
internal/api/      # HTTP handlers
internal/packs/    # Core pack allocation logic and tests
```

## Notes

- The service intentionally uses only the Go standard library to keep dependencies minimal.
- Pack sizes can be changed without code modifications.
- With more time, the service could be extended with OpenAPI documentation, metrics, tracing, and rate limiting.
