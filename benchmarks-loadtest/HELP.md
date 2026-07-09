# HELP - Running Benchmarks Manually

Step-by-step guide for running benchmarks on this project.

---

## 1. Go Micro-Benchmarks (No Server Needed)

These test individual components without a running server.

### Run All Benchmarks

```powershell
# From project root
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./models/ ./middleware/ ./repository/ ./controller/
```

### Run Per Package

```powershell
# Models (password hashing)
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./models/

# Middleware (CORS, logging, rate limit, etc.)
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./middleware/

# Repository (GORM database operations)
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./repository/

# Controller (HTTP handlers)
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./controller/
```

### Run Single Benchmark

```powershell
# Example: only CreateUser handler
go test -run='^$' -bench=BenchmarkCreateUser -benchmem -count=3 ./controller/

# Example: only bcrypt hash
go test -run='^$' -bench=BenchmarkHashPassword -benchmem -count=3 ./models/

# Example: only CORS middleware
go test -run='^$' -bench=BenchmarkCORS -benchmem -count=3 ./middleware/
```

### Save Results to File

```powershell
# Save to benchmarks/results/
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./models/ > benchmarks\results\go_models.txt 2>&1
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./middleware/ > benchmarks\results\go_middleware.txt 2>&1
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./repository/ > benchmarks\results\go_repository.txt 2>&1
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./controller/ > benchmarks\results\go_controller.txt 2>&1
```

### Compare Two Runs (Requires benchstat)

```powershell
# Install benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Save baseline
go test -run='^$' -bench='Benchmark' -benchmem -count=5 ./controller/ > benchmarks\results\baseline.txt

# Make changes, then save new run
go test -run='^$' -bench='Benchmark' -benchmem -count=5 ./controller/ > benchmarks\results\after.txt

# Compare
benchstat benchmarks\results\baseline.txt benchmarks\results\after.txt
```

---

## 2. Load Tests (Server Must Be Running)

These test the HTTP API under load. Start the server first.

### Start Server

```powershell
# Terminal 1: Start PostgreSQL
docker compose up -d

# Terminal 2: Start API
go run main.go
```

Server runs at `http://localhost:8080` by default.

---

### hey (Simple HTTP Load)

**Install:**
```powershell
go install github.com/rakyll/hey@latest
```

**Run:**
```powershell
# Default: 1000 requests, 50 concurrency
hey http://localhost:8080/health

# Custom: 5000 requests, 200 concurrency
hey -n 5000 -c 200 http://localhost:8080/health

# POST with JSON body
hey -n 100 -c 10 -m POST -H "Content-Type: application/json" ^
    -d "{\"name\":\"Test\",\"email\":\"test@example.com\",\"password\":\"secret123\"}" ^
    http://localhost:8080/api/v1/users

# GET users list
hey -n 1000 -c 50 http://localhost:8080/api/v1/users

# GET single user
hey -n 1000 -c 50 http://localhost:8080/api/v1/users/1
```

**Output columns:** Requests/sec, Average response time, Fastest/Slowest, Status code distribution, Latency histogram.

---

### vegeta (Constant Rate Attack)

**Install:**
```powershell
go install github.com/tsenart/vegeta@latest
```

**Run:**
```powershell
# 50 requests/sec for 30 seconds against /health
echo "GET http://localhost:8080/health" | vegeta attack -rate=50 -duration=30s | vegeta report

# 100 req/s for 60 seconds against users endpoint
echo "GET http://localhost:8080/api/v1/users" | vegeta attack -rate=100 -duration=60s | vegeta report

# Save results and analyze later
echo "GET http://localhost:8080/health" | vegeta attack -rate=50 -duration=30s > results.bin
vegeta report results.bin
vegeta plot results.bin > plot.html

# Multiple endpoints
(echo "GET http://localhost:8080/health"; echo "GET http://localhost:8080/api/v1/users") | vegeta attack -rate=25 -duration=20s | vegeta report
```

**Output:** Success rate, latency percentiles (p50, p95, p99), throughput.

---

### wrk (High Concurrency)

**Install:**
```powershell
# Windows (via WSL)
wsl -d Ubuntu-24.04 -u root -- apt-get install -y wrk

# macOS
brew install wrk

# Linux
sudo apt install wrk
```

**Run:**
```powershell
# Windows (via WSL)
wsl -d Ubuntu-24.04 -- wrk -t12 -c400 -d30s http://localhost:8080/health

# Linux / macOS
wrk -t12 -c400 -d30s http://localhost:8080/health

# Lighter load
wrk -t4 -c50 -d15s http://localhost:8080/api/v1/users
```

**Flags:** `-t` threads, `-c` connections, `-d` duration.

---

### k6 (Scripted Scenarios)

**Install:**
```powershell
# Windows
winget install GrafanaLabs.k6

# macOS
brew install k6

# Linux
sudo apt-get install -y k6
```

**Run:**
```powershell
# Run the included script
k6 run benchmarks\k6.js

# Custom: 10 VUs for 30 seconds
k6 run --vus 10 --duration 30s benchmarks\k6.js

# With custom threshold output
k6 run --out json=benchmarks\results\k6_output.json benchmarks\k6.js
```

**What the included script does:**
1. **Health check scenario**: 10 constant VUs hitting `/health` for 30s
2. **CRUD scenario**: Ramping VUs (0 -> 20 -> 50 -> 0) doing Create -> List -> Get -> Update -> Delete
3. **Thresholds**: p95 < 500ms, p99 < 1000ms, error rate < 10%

---

## 3. Useful Flags Reference

### Go Benchmark Flags

| Flag | Purpose | Example |
|------|---------|---------|
| `-bench=` | Filter benchmarks by regex | `-bench=BenchmarkCreate` |
| `-benchmem` | Show memory allocation stats | `-benchmem` |
| `-benchtime=` | Duration per benchmark | `-benchtime=5s` |
| `-count=` | Number of benchmark runs | `-count=5` |
| `-cpu=` | Set GOMAXPROCS | `-cpu=1,2,4,8` |
| `-timeout=` | Timeout for the test | `-timeout=300s` |
| `-run='^$'` | Skip unit tests (only benchmarks) | Required to run only benchmarks |

### Common Patterns

```powershell
# Run benchmarks 5 times for statistical significance
go test -run='^$' -bench='Benchmark' -benchmem -count=5 ./controller/

# Run with different CPU counts to test scaling
go test -run='^$' -bench='Benchmark' -benchmem -cpu=1,2,4,8 ./repository/

# Run only benchmarks matching "Create"
go test -run='^$' -bench=BenchmarkCreate -benchmem ./controller/

# Longer benchmark time for stable results
go test -run='^$' -bench='Benchmark' -benchmem -benchtime=10s ./models/

# Skip slow benchmarks
go test -run='^$' -bench='Benchmark[^C]' -benchmem ./controller/
```

---

## 4. Troubleshooting

| Problem | Solution |
|---------|----------|
| `no such table: table_of_user` in repository benchmarks | Normal — SQLite in-memory DB is recreated per benchmark iteration |
| `panic: nil pointer dereference` in BenchmarkFullMiddlewareChain | Known issue — Logger middleware needs `log.Default()` initialized. Skip with `-bench='Benchmark[^F]'` |
| Benchmarks too fast / unstable | Increase `-benchtime=10s` or `-count=5` |
| Connection refused on load tests | Make sure server is running (`go run main.go`) |
| `hey`/`vegeta` not found | Run `go install github.com/rakyll/hey@latest` and `go install github.com/tsenart/vegeta@latest` |
| `wrk` not found on Windows | Use WSL: `wsl -d Ubuntu-24.04 -- wrk ...` |
| High memory in benchmarks | Close other applications, ensure machine is plugged in (not battery) |

---

## 5. Recommended Benchmark Workflow

```
1. Save baseline
   go test -run='^$' -bench='Benchmark' -benchmem -count=5 ./controller/ > baseline.txt

2. Make code changes

3. Save new results
   go test -run='^$' -bench='Benchmark' -benchmem -count=5 ./controller/ > after.txt

4. Compare
   benchstat baseline.txt after.txt

5. Open HTML report for visual review
   benchmarks/results/benchmark_report.html
```
