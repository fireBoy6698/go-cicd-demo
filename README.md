# Hello World API

Simple Go REST API with CI/CD using GitHub Actions

## Endpoints

- `GET /hello` - Returns Hello, World! message
- `GET /health` - Health check endpoint

## Getting Started

### Run locally

```bash
# Run server
go run main.go

# Run tests
go test -v

# Build
go build -o app
./app
```

### Test API

```bash
curl http://localhost:8080/hello
curl http://localhost:8080/health
```

## CI/CD

This project uses GitHub Actions for automated testing and building.

- On every push/PR to `main` branch:
  - Runs all tests
  - Generates coverage report
  - Builds the binary
