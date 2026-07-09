# Benchmarks & Load Testing

Performance testing suite for the Gin-GORM-Postgres API.

**[View HTML Benchmark Report](results/benchmark_report.html)**

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Go Micro-Benchmarks](#go-micro-benchmarks)
- [Load Testing Tools](#load-testing-tools)
- [Understanding Results](#understanding-results)

---

## Prerequisites

### 1. Start the Server

```bash
# Start PostgreSQL
docker compose up -d

# Run the API server
go run main.go
```

The server starts on `http://localhost:8080` by default.

### 2. Install Load Testing Tools

```bash
# hey + vegeta + benchstat (Go tools)
go install github.com/rakyll/hey@latest
go install github.com/tsenart/vegeta@latest
go install golang.org/x/perf/cmd/benchstat@latest

# k6 (Windows)
winget install GrafanaLabs.k6

# k6 (macOS)
brew install k6

# k6 (Linux)
sudo apt-get install -y k6

# wrk (Windows - requires WSL)
wsl -d Ubuntu-24.04 -u root -- apt-get install -y wrk

# wrk (macOS)
brew install wrk

# wrk (Linux)
sudo apt install wrk
```

### Tool Availability

| Tool | Purpose | Install |
|------|---------|---------|
| `go test -bench` | Go built-in benchmarks | Always available |
| `hey` | Simple HTTP load testing | `go install` |
| `vegeta` | Constant-rate attack tool | `go install` |
| `wrk` | High-performance HTTP benchmarking | apt/brew |
| `k6` | Scripted load testing with scenarios | winget/brew/apt |
| `benchstat` | Statistical comparison of Go benchmarks | `go install` |

---

## Quick Start

### Run Go Benchmarks (No Server Required)

```bash
# Run all benchmarks across all packages
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./models/ ./middleware/ ./repository/ ./controller/

# Run benchmarks for a specific package
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./models/
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./middleware/
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./repository/
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./controller/

# Run a single benchmark
go test -run='^$' -bench=BenchmarkCreateUser -benchmem -count=3 ./controller/
```

### Run Load Tests (Server Must Be Running)

```bash
# Windows
benchmarks\windows\hey.bat
benchmarks\windows\vegeta.bat
benchmarks\windows\wrk.bat
benchmarks\windows\k6.bat

# Linux / macOS
bash benchmarks/linux/hey.sh
bash benchmarks/linux/vegeta.sh
bash benchmarks/linux/wrk.sh
bash benchmarks/linux/k6.sh
```

---

## Go Micro-Benchmarks

Built-in Go benchmarks test individual components in isolation.

### What's Tested

| Package | Benchmarks | What It Measures |
|---------|-----------|------------------|
| `models/` | 3 | Bcrypt hash/verify operations |
| `middleware/` | 9 | RequestID, CORS, SecurityHeaders, RateLimit, Logger, Recovery, ErrorHandler, Full chain |
| `repository/` | 8 | GORM CRUD + parallel operations on SQLite |
| `controller/` | 7 | HTTP handlers with mock repo (Create, List, Get, Update, Delete, Full CRUD cycle) |

### Benchmark Results

Results are saved in `benchmarks/results/`:

| File | Contents |
|------|----------|
| `go_models.txt` | Password hashing benchmarks |
| `go_middleware.txt` | Middleware benchmarks |
| `go_repository.txt` | Repository/GORM benchmarks |
| `go_controller.txt` | Controller handler benchmarks |
| `benchmark_report.html` | Visual HTML report (open in browser) |

### Comparing Results

```bash
# Save baseline
go test -run='^$' -bench='Benchmark' -benchmem -count=5 ./models/ > benchmarks/results/baseline.txt

# After changes, save new results
go test -run='^$' -bench='Benchmark' -benchmem -count=5 ./models/ > benchmarks/results/after.txt

# Compare statistically
benchstat benchmarks/results/baseline.txt benchmarks/results/after.txt
```

---

## Load Testing Tools

### hey - Simple HTTP Load Test

Best for: Quick smoke tests, single endpoint testing.

```bash
# Default: 1000 requests, 50 concurrency
hey http://localhost:8080/health

# Custom: 5000 requests, 200 concurrency
hey -n 5000 -c 200 http://localhost:8080/api/v1/users

# POST request with JSON body
hey -n 1000 -c 50 -m POST -H "Content-Type: application/json" \
    -d '{"name":"Test","email":"test@example.com","password":"secret123"}' \
    http://localhost:8080/api/v1/users
```

### vegeta - Constant Rate Attack

Best for: Throughput testing at fixed requests-per-second, latency percentiles.

```bash
# Attack at 50 req/s for 30 seconds
echo "GET http://localhost:8080/health" | vegeta attack -rate=50 -duration=30s | vegeta report

# Full CRUD attack
echo "GET http://localhost:8080/api/v1/users" | vegeta attack -rate=100 -duration=60s | vegeta report

# Save and analyze results
echo "GET http://localhost:8080/health" | vegeta attack -rate=50 -duration=30s > results.bin
vegeta report results.bin
vegeta plot results.bin > plot.html
```

### wrk - High-Performance Benchmarking

Best for: High concurrency testing, raw throughput.

```bash
# Windows (via WSL)
wsl -d Ubuntu-24.04 -- wrk -t12 -c400 -d30s http://localhost:8080/health

# Linux / macOS
wrk -t12 -c400 -d30s http://localhost:8080/health

# Custom endpoint
wrk -t8 -c200 -d15s http://localhost:8080/api/v1/users
```

### k6 - Scripted Load Testing

Best for: Complex scenarios, ramping virtual users, custom thresholds.

```bash
# Run the included k6 script
k6 run benchmarks/k6.js

# Custom scenario
k6 run --vus 10 --duration 30s benchmarks/k6.js
```

The included `k6.js` script tests:
- **Health check**: 10 constant VUs for 30s
- **Full CRUD**: Ramping VUs (0 -> 20 -> 50 -> 0) over 40s
- **Thresholds**: p95 < 500ms, p99 < 1000ms, error rate < 10%

---

## Understanding Results

### Go Benchmark Output

```
BenchmarkCreateUser-8    5    339000000 ns/op    26275 B/op    120 allocs/op
│                   │    │    │                │              │              │
│                   │    │    │                │              │              └─ allocations per operation
│                   │    │    │                │              └─ bytes allocated per operation
│                   │    │    │                └─ nanoseconds per operation
│                   │    │    └─ number of iterations
│                   │    └─ GOMAXPROCS value
│                   └─ benchmark function name (suffix: -8 = GOMAXPROCS)
└─ package name
```

### Key Metrics

| Metric | Good | Okay | Bad |
|--------|------|------|-----|
| ns/op (read) | < 100 &mu;s | 100-500 &mu;s | > 500 &mu;s |
| ns/op (write) | < 1 ms | 1-10 ms | > 10 ms |
| B/op | < 10 KB | 10-100 KB | > 100 KB |
| allocs/op | < 20 | 20-100 | > 100 |

### What Affects Results

- **Machine load**: Close other applications during benchmarks
- **CPU throttling**: Run on AC power, not battery
- **GOMAXPROCS**: Usually matches CPU core count
- **GC pressure**: `count=3` or higher smooths out GC pauses
- **Database**: Repository benchmarks use SQLite in-memory; production PostgreSQL will differ

---

## Files

```
benchmarks/
  README.md                  - This file
  k6.js                      - k6 load test script
  results/
    go_models.txt            - Models benchmark raw output
    go_middleware.txt         - Middleware benchmark raw output
    go_repository.txt        - Repository benchmark raw output
    go_controller.txt        - Controller benchmark raw output
    benchmark_report.html    - Visual HTML report (open in browser)
  windows/
    hey.bat                  - hey load test runner
    vegeta.bat               - vegeta load test runner
    wrk.bat                  - wrk load test runner (via WSL)
    k6.bat                   - k6 load test runner
  linux/
    hey.sh                   - hey load test runner
    vegeta.sh                - vegeta load test runner
    wrk.sh                   - wrk load test runner
    k6.sh                    - k6 load test runner
```
