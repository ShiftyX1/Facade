# Shifty's Facade

[![Go Version](https://img.shields.io/github/go-mod/go-version/ShiftyX1/Facade)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/releases)
[![License](https://img.shields.io/github/license/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/blob/main/LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/ShiftyX1/Facade/ci.yml?branch=main)](https://github.com/ShiftyX1/Facade/actions)
[![GitHub Issues](https://img.shields.io/github/issues/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/issues)
[![GitHub Stars](https://img.shields.io/github/stars/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/stargazers)
[![Last Commit](https://img.shields.io/github/last-commit/ShiftyX1/Facade)](https://github.com/ShiftyX1/Facade/commits/main)

A simple and flexible mock HTTP server for developers. <br>
[Russian README](./README-ru.md)

## Installation

```bash
make build
```

## Quick Start

1. Create a configuration file `config.yaml`:

```yaml
server:
  port: 8080
  host: localhost

routes:
  - path: /users
    method: GET
    response:
      status: 200
      schema:
        type: array
        items:
          type: object
          properties:
            id: {type: integer, min: 1, max: 1000}
            name: {type: string, faker: name}
            email: {type: string, faker: email}
        count: 5

global:
  cors: true
  log_requests: true
```

2. Start the server:

```bash
make run
```

Or build and run manually:
```bash
make build
./build/facade -c configs/example.yaml
```

3. Test the API:

```bash
curl http://localhost:8080/users
```

## Command Line Options

- `-c, --config` - path to configuration file (default: configs/example.yaml)
- `-p, --port` - server port (overrides config)
- `--host` - server host (overrides config)
- `-v, --verbose` - verbose logging

## Configuration

### Basic Structure

```yaml
server:
  port: 8080
  host: localhost

routes:
  - path: /api/endpoint
    method: GET
    response:
      status: 200
      body: '{"message": "Hello World"}'

global:
  cors: true
  delay: 100ms
  log_requests: true
```

### Response Types

#### Static Response

```yaml
- path: /static
  method: GET
  response:
    status: 200
    body: '{"message": "Static response"}'
```

#### Templated Response

```yaml
- path: /users/{id}
  method: GET
  response:
    status: 200
    body: '{"id": {{.PathParams.id}}, "name": "User {{.PathParams.id}}"}'
```

#### Schema-based Generation

```yaml
- path: /users
  method: GET
  response:
    status: 200
    schema:
      type: array
      items:
        type: object
        properties:
          id: {type: integer, min: 1, max: 1000}
          name: {type: string, faker: name}
          email: {type: string, faker: email}
      count: 5
```

### Conditional Responses

```yaml
- path: /data
  method: GET
  response:
    status: 200
    conditions:
      - if: '{{.QueryParams.format}} == "xml"'
        body: '<data>XML response</data>'
      - if: '{{.QueryParams.format}} == "json"'
        body: '{"data": "JSON response"}'
      - default: true
        body: '{"data": "Default response"}'
```

### Template Functions

#### Request Data Access

- `{{.PathParams.id}}` - path parameters
- `{{.QueryParams.limit}}` - query parameters  
- `{{.Body.name}}` - request body data

#### Faker Data

- `{{fakerName}}` - random name
- `{{fakerEmail}}` - random email
- `{{fakerPhone}}` - random phone
- `{{fakerCompany}}` - random company
- `{{fakerUUID}}` - UUID

#### Random Data

- `{{randomInt 1 100}}` - random number
- `{{randomString 10}}` - random string
- `{{randomBool}}` - random boolean

#### Date and Time

- `{{now}}` - current time
- `{{now.Format "2006-01-02"}}` - formatted date

### Schema Types

#### Primitive Types

```yaml
properties:
  name: {type: string, faker: name}
  age: {type: integer, min: 18, max: 65}
  price: {type: number, min: 0.1, max: 999.99}
  active: {type: boolean}
```

#### Arrays

```yaml
schema:
  type: array
  items:
    type: object
    properties:
      id: {type: integer}
      name: {type: string}
  count: 10
```

#### Objects

```yaml
schema:
  type: object
  properties:
    user:
      type: object
      properties:
        id: {type: integer}
        profile:
          type: object
          properties:
            name: {type: string, faker: name}
```

### Faker Types

- `name` - person name
- `email` - email address
- `phone` - phone number
- `address` - address
- `company` - company name
- `word` - random word
- `sentence` - sentence
- `paragraph` - paragraph
- `url` - URL
- `uuid` - UUID
- `date` - date
- `timestamp` - timestamp

### System Endpoints

Facade provides system endpoints for management:

- `GET /_facade/health` - server health check
- `GET /_facade/state` - get entire state
- `DELETE /_facade/state` - clear state
- `GET /_facade/state/keys` - get state keys

### Stateful Mode

Facade automatically saves data from POST/PUT requests and can use it in subsequent GET requests.

Example of saving a user:

```bash
# Create user
curl -X POST http://localhost:8080/users \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Get state
curl http://localhost:8080/_facade/state
```

## Build and Run Commands

### Main Makefile Commands

```bash
make help          # Show all available commands
make build         # Build application to build/facade
make run           # Build and run with example configuration
make clean         # Clean build files
make deps          # Install/update dependencies
make build-all     # Cross-compile for all platforms
```

### Running with Different Configurations

```bash
# Run with example configuration
make run

# Run with custom configuration
./build/facade -c configs/frontend-dev.yaml -v

# Run with port override
./build/facade -c configs/testing.yaml -p 9000 -v
```

## Testing

### Basic Testing

After starting the server (`make run`), you can test the API:

```bash
# Check server health
curl http://localhost:8080/_facade/health

# Get list of users (schema generation)
curl http://localhost:8080/users | jq

# Get specific user (templating)
curl http://localhost:8080/users/123 | jq

# Create user (POST with data)
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}' | jq

# Check state after POST
curl http://localhost:8080/_facade/state | jq

# Test conditional logic
curl "http://localhost:8080/products?category=electronics" | jq
curl "http://localhost:8080/products" | jq

# Test query parameters
curl "http://localhost:8080/search?q=test&limit=10" | jq
```

### Testing testing.yaml Configuration

```bash
# Start test server
./build/facade -c configs/testing.yaml -v

# Successful payment
curl -X POST http://localhost:8081/api/payments \
  -H "Content-Type: application/json" \
  -d '{"amount": "500", "card_number": "4111111111111111"}' | jq

# Payment with verification (amount > 1000)
curl -X POST http://localhost:8081/api/payments \
  -H "Content-Type: application/json" \
  -d '{"amount": "1500", "card_number": "4111111111111111"}' | jq

# Declined payment
curl -X POST http://localhost:8081/api/payments \
  -H "Content-Type: application/json" \
  -d '{"amount": "100", "card_number": "4000000000000002"}' | jq

# Order statuses
curl http://localhost:8081/api/orders/1/status | jq  # pending
curl http://localhost:8081/api/orders/2/status | jq  # shipped
curl http://localhost:8081/api/orders/3/status | jq  # delivered
curl http://localhost:8081/api/orders/999/status | jq # not found

# Error simulation
curl http://localhost:8081/api/error/500 | jq
curl http://localhost:8081/api/error/404 | jq
curl http://localhost:8081/api/error/401 | jq

# Slow endpoint (5 second delay)
time curl http://localhost:8081/api/slow | jq
```

## Usage Examples

### Frontend Development

Set up Facade as a backend for frontend application development:

```yaml
server:
  port: 3001

routes:
  - path: /api/auth/login
    method: POST
    response:
      status: 200
      body: '{"token": "{{fakerUUID}}", "user": {"id": 1, "name": "{{.Body.username}}"}}'

  - path: /api/users
    method: GET
    response:
      status: 200
      schema:
        type: array
        items:
          type: object
          properties:
            id: {type: integer, min: 1, max: 1000}
            name: {type: string, faker: name}
            avatar: {type: string, faker: url}
        count: 20

global:
  cors: true
  delay: 200ms
```

### Integration Testing

Simulate various scenarios for testing:

```yaml
routes:
  - path: /api/payment
    method: POST
    response:
      status: 200
      conditions:
        - if: '{{.Body.amount}} > "1000"'
          body: '{"status": "pending", "id": "{{fakerUUID}}"}'
        - if: '{{.Body.card_number}} == "4111111111111111"'
          body: '{"status": "success", "id": "{{fakerUUID}}"}'
        - default: true
          status: 400
          body: '{"error": "Invalid card"}'
```

### API Demonstration

Create realistic APIs for demonstrations:

```yaml
routes:
  - path: /api/products
    method: GET
    response:
      status: 200
      schema:
        type: object
        properties:
          data:
            type: array
            items:
              type: object
              properties:
                id: {type: integer, min: 1, max: 1000}
                name: {type: string, faker: word}
                description: {type: string, faker: sentence}
                price: {type: number, min: 10, max: 500}
                image: {type: string, faker: url}
                category: {type: string, faker: word}
            count: 12
          total: {type: integer, min: 100, max: 1000}
          page: 1
          per_page: 12
```

## Build and Deployment

### Using Makefile (Recommended)

```bash
# Show all available commands
make help

# Quick start - build and run
make run

# Build only
make build

# Clean up
make clean

# Cross-compile for all platforms
make build-all
```

### Manual Build

```bash
# Local build
mkdir -p build
go build -o build/facade cmd/facade/main.go

# Run
./build/facade -c configs/example.yaml
```

## License

MIT License - see LICENSE file for details.