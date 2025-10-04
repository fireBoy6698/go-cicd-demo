# CI/CD Workshop สำหรับ Go Lang 🚀

## เตรียมตัว

### ติดตั้ง Go
```bash
# ตรวจสอบว่าติดตั้ง Go แล้ว
go version

# ควรได้ Go 1.20 ขึ้นไป
```

---

## Workshop 1: สร้างโปรเจค Go พื้นฐาน

### ขั้นตอนที่ 1: สร้างโปรเจค
```bash
mkdir go-cicd-workshop
cd go-cicd-workshop
go mod init github.com/yourusername/go-cicd-workshop
```

### ขั้นตอนที่ 2: สร้าง Web Server ง่ายๆ

**main.go**
```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, CI/CD with Go!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "OK")
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/health", healthHandler)

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

### ขั้นตอนที่ 3: ทดสอบรันโปรเจค
```bash
go run main.go

# เปิดเบราว์เซอร์ไปที่ http://localhost:8080
```

---

## Workshop 2: เขียน Tests

### สร้างไฟล์ทดสอบ

**main_test.go**
```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHelloHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(helloHandler)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := "Hello, CI/CD with Go!"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}

func TestHealthHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/health", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(healthHandler)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := "OK"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}
```

### รัน Tests
```bash
# รัน tests ทั้งหมด
go test

# รัน tests พร้อมรายละเอียด
go test -v

# ดู test coverage
go test -cover

# สร้าง coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Workshop 3: ตั้งค่า GitHub Actions

### สร้างไฟล์ Workflow

**.github/workflows/go-ci.yml**
```yaml
name: Go CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test -v -cover ./...

    - name: Build
      run: go build -v ./...

    - name: Run go vet
      run: go vet ./...

    - name: Run staticcheck
      run: |
        go install honnef.co/go/tools/cmd/staticcheck@latest
        staticcheck ./...
```

---

## Workshop 4: เพิ่ม Code Quality Checks

### สร้าง Workflow ที่สมบูรณ์

**.github/workflows/go-quality.yml**
```yaml
name: Go Quality

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    name: Run Tests
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Test
      run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out

  lint:
    name: Lint Code
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest

  build:
    name: Build Binary
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Build
      run: |
        go build -v -o app ./...

    - name: Upload artifact
      uses: actions/upload-artifact@v3
      with:
        name: go-app
        path: app
```

---

## Workshop 5: Build และ Deploy

### Build สำหรับหลาย Platform

**.github/workflows/release.yml**
```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Build for multiple platforms
      run: |
        # Linux
        GOOS=linux GOARCH=amd64 go build -o build/app-linux-amd64

        # macOS
        GOOS=darwin GOARCH=amd64 go build -o build/app-darwin-amd64
        GOOS=darwin GOARCH=arm64 go build -o build/app-darwin-arm64

        # Windows
        GOOS=windows GOARCH=amd64 go build -o build/app-windows-amd64.exe

    - name: Create Release
      uses: softprops/action-gh-release@v1
      with:
        files: build/*
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Workshop 6: Docker Integration

### Dockerfile

**Dockerfile**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
```

### Build และรัน Docker

```bash
# Build image
docker build -t go-cicd-app .

# Run container
docker run -p 8080:8080 go-cicd-app

# ทดสอบ
curl http://localhost:8080
```

### GitHub Actions สำหรับ Docker

**.github/workflows/docker.yml**
```yaml
name: Docker Build

on:
  push:
    branches: [ main ]

jobs:
  docker:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Login to Docker Hub
      uses: docker/login-action@v2
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}

    - name: Build and push
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: yourusername/go-cicd-app:latest
```

---

## Workshop 7: REST API พร้อม Database

### โปรเจคที่ซับซ้อนขึ้น

**main.go**
```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type Store struct {
    mu    sync.RWMutex
    users map[int]User
    nextID int
}

func NewStore() *Store {
    return &Store{
        users: make(map[int]User),
        nextID: 1,
    }
}

func (s *Store) CreateUser(name, email string) User {
    s.mu.Lock()
    defer s.mu.Unlock()

    user := User{
        ID:    s.nextID,
        Name:  name,
        Email: email,
    }
    s.users[s.nextID] = user
    s.nextID++
    return user
}

func (s *Store) GetUsers() []User {
    s.mu.RLock()
    defer s.mu.RUnlock()

    users := make([]User, 0, len(s.users))
    for _, u := range s.users {
        users = append(users, u)
    }
    return users
}

type Server struct {
    store *Store
}

func (srv *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        users := srv.store.GetUsers()
        json.NewEncoder(w).Encode(users)

    case http.MethodPost:
        var req struct {
            Name  string `json:"name"`
            Email string `json:"email"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        user := srv.store.CreateUser(req.Name, req.Email)
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)

    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func main() {
    store := NewStore()
    server := &Server{store: store}

    http.HandleFunc("/users", server.handleUsers)
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Tests สำหรับ REST API

**main_test.go**
```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateUser(t *testing.T) {
    store := NewStore()
    server := &Server{store: store}

    reqBody := bytes.NewBufferString(`{"name":"John","email":"john@example.com"}`)
    req := httptest.NewRequest(http.MethodPost, "/users", reqBody)
    w := httptest.NewRecorder()

    server.handleUsers(w, req)

    if w.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", w.Code)
    }

    var user User
    json.NewDecoder(w.Body).Decode(&user)

    if user.Name != "John" {
        t.Errorf("Expected name John, got %s", user.Name)
    }
}

func TestGetUsers(t *testing.T) {
    store := NewStore()
    store.CreateUser("Alice", "alice@example.com")
    store.CreateUser("Bob", "bob@example.com")

    server := &Server{store: store}

    req := httptest.NewRequest(http.MethodGet, "/users", nil)
    w := httptest.NewRecorder()

    server.handleUsers(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }

    var users []User
    json.NewDecoder(w.Body).Decode(&users)

    if len(users) != 2 {
        t.Errorf("Expected 2 users, got %d", len(users))
    }
}
```

---

## Workshop 8: Environment Variables และ Config

### การจัดการ Config

**config.go**
```go
package main

import (
    "os"
    "strconv"
)

type Config struct {
    Port     string
    DBHost   string
    DBPort   int
    LogLevel string
}

func LoadConfig() *Config {
    dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

    return &Config{
        Port:     getEnv("PORT", "8080"),
        DBHost:   getEnv("DB_HOST", "localhost"),
        DBPort:   dbPort,
        LogLevel: getEnv("LOG_LEVEL", "info"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### .env ไฟล์ (สำหรับ local development)

**.env**
```bash
PORT=8080
DB_HOST=localhost
DB_PORT=5432
LOG_LEVEL=debug
```

### .github/workflows ที่ใช้ secrets

```yaml
name: Deploy

on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest

    env:
      PORT: ${{ secrets.PORT }}
      DB_HOST: ${{ secrets.DB_HOST }}
      DB_PORT: ${{ secrets.DB_PORT }}

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Run tests
      run: go test -v ./...

    - name: Build
      run: go build -o app

    - name: Deploy to server
      run: |
        echo "Deploying with config:"
        echo "Port: $PORT"
        echo "DB Host: $DB_HOST"
        # เพิ่มคำสั่ง deploy ที่นี่
```

---

## Workshop 9: Integration Tests

### Docker Compose สำหรับ Testing

**docker-compose.test.yml**
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
    depends_on:
      - postgres

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=testuser
      - POSTGRES_PASSWORD=testpass
      - POSTGRES_DB=testdb
    ports:
      - "5432:5432"
```

### Integration Test Script

**scripts/integration-test.sh**
```bash
#!/bin/bash

echo "Starting integration tests..."

# Start services
docker-compose -f docker-compose.test.yml up -d

# Wait for services to be ready
sleep 5

# Run tests
curl -f http://localhost:8080/health || exit 1

echo "Testing create user..."
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com"}' || exit 1

echo "Testing get users..."
curl -f http://localhost:8080/users || exit 1

# Cleanup
docker-compose -f docker-compose.test.yml down

echo "Integration tests passed!"
```

---

## Workshop 10: Performance Testing

### Benchmark Tests

**benchmark_test.go**
```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func BenchmarkHelloHandler(b *testing.B) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        helloHandler(w, req)
    }
}

func BenchmarkHealthHandler(b *testing.B) {
    req := httptest.NewRequest(http.MethodGet, "/health", nil)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        healthHandler(w, req)
    }
}
```

### รัน Benchmarks

```bash
# รัน benchmark
go test -bench=.

# รัน benchmark พร้อม memory profiling
go test -bench=. -benchmem

# รัน benchmark และสร้าง CPU profile
go test -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

---

## สรุป Checklist

### ✅ CI/CD ที่สมบูรณ์สำหรับ Go

- [ ] Go module (`go.mod`)
- [ ] Unit tests (`*_test.go`)
- [ ] Benchmark tests
- [ ] Code coverage > 80%
- [ ] GitHub Actions workflow
- [ ] Linting (golangci-lint)
- [ ] Static analysis (staticcheck, go vet)
- [ ] Dockerfile
- [ ] Docker Compose
- [ ] Integration tests
- [ ] Environment configuration
- [ ] Health check endpoint
- [ ] Logging
- [ ] Error handling
- [ ] Documentation

---

## คำสั่งที่ใช้บ่อย

```bash
# Development
go run main.go
go test -v ./...
go test -cover ./...

# Building
go build -o app
GOOS=linux GOARCH=amd64 go build -o app-linux

# Testing
go test -v ./...
go test -race ./...
go test -bench=.

# Code Quality
go fmt ./...
go vet ./...
golangci-lint run

# Docker
docker build -t myapp .
docker run -p 8080:8080 myapp

# Modules
go mod init
go mod tidy
go mod download
```

---

## ทรัพยากรเพิ่มเติม

📚 **เอกสารสำคัญ:**
- [Go Documentation](https://go.dev/doc/)
- [Go Testing](https://go.dev/doc/tutorial/add-a-test)
- [GitHub Actions for Go](https://github.com/actions/setup-go)

🛠️ **Tools แนะนำ:**
- [golangci-lint](https://golangci-lint.run/) - Linter
- [Air](https://github.com/cosmtrek/air) - Live reload
- [Delve](https://github.com/go-delve/delve) - Debugger

---

## Next Steps

1. ✅ เริ่มจาก Workshop 1-3 ก่อน
2. ✅ เพิ่ม GitHub Actions
3. ✅ เพิ่ม Docker support
4. ✅ ทำ Integration tests
5. ✅ Deploy ขึ้น production

**เริ่มต้นได้เลย! สร้างโปรเจค Go แรกและเพิ่ม CI/CD ทีละขั้นตอน 🎯**
