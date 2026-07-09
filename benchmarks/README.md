# Load Testing Benchmarks

Load testing scripts for the Gin-GORM-Postgres API using **hey**, **vegeta**, **wrk**, and **k6**.

## Prerequisites

Start the server first:

```bash
# From project root
docker compose up -d        # Start PostgreSQL
go run main.go              # Start API on :8080
```

## Installed Tools

| Tool | Status | Location |
|------|--------|----------|
| **hey** | Installed | `$HOME/go/bin/hey.exe` |
| **vegeta** | Installed | `$HOME/go/bin/vegeta.exe` |
| **benchstat** | Installed | `$HOME/go/bin/benchstat.exe` |
| **k6** | Installed | `C:\Program Files\k6\k6.exe` |
| **wrk** | Installed | WSL Ubuntu (`wsl -d Ubuntu-24.04 -- wrk`) |

## Install Tools (Manual)

```bash
# hey + vegeta + benchstat (Go tools)
go install github.com/rakyll/hey@latest
go install github.com/tsenart/vegeta@latest
go install golang.org/x/perf/cmd/benchstat@latest

# k6 (Windows)
winget install GrafanaLabs.k6

# wrk (Windows - requires WSL)
wsl -d Ubuntu-24.04 -u root -- apt-get install -y wrk

# wrk (macOS)
brew install wrk

# wrk (Linux)
sudo apt install wrk
```

## Usage

All scripts respect `BASE_URL` (default: `http://localhost:8080`).

### Windows

```cmd
benchmarks\windows\hey.bat
benchmarks\windows\vegeta.bat
benchmarks\windows\wrk.bat
benchmarks\windows\k6.bat
```

### Linux / macOS

```bash
bash benchmarks/linux/hey.sh
bash benchmarks/linux/vegeta.sh
bash benchmarks/linux/wrk.sh
bash benchmarks/linux/k6.sh
```

### Override Defaults

```cmd
REM Windows
set BASE_URL=http://localhost:9090
set REQUESTS=5000
set CONCURRENCY=100
benchmarks\windows\hey.bat
```

```bash
# Linux / macOS
BASE_URL=http://localhost:9090 REQUESTS=5000 CONCURRENCY=100 bash benchmarks/linux/hey.sh
```

## Go Benchmarks

Micro-benchmarks using Go's built-in `testing.B`:

```bash
# Run all benchmarks with memory stats (Windows PowerShell)
go test -run='^$' -bench='Benchmark' -benchmem -count=3 ./models/ ./middleware/ ./repository/ ./controller/

# Run specific benchmark
go test -run='^$' -bench=BenchmarkCreateUser -benchmem ./controller/

# Compare runs with benchstat
go test -run='^$' -bench='Benchmark' -benchmem -count=5 ./models/ > bench_v1.txt
# ... make changes ...
go test -run='^$' -bench='Benchmark' -benchmem -count=5 ./models/ > bench_v2.txt
benchstat bench_v1.txt bench_v2.txt
```

## Files

```
benchmarks/
  windows/
    hey.bat              - hey load test
    vegeta.bat           - vegeta load test
    wrk.bat              - wrk load test (via WSL)
    k6.bat               - k6 runner
  linux/
    hey.sh               - hey load test
    vegeta.sh            - vegeta load test
    wrk.sh               - wrk load test
    k6.sh                - k6 runner
  k6.js                  - k6 script
  README.md              - This file
```

## What Each Tool Measures

| Tool | Type | Best For |
|------|------|----------|
| **hey** | Simple HTTP load | Quick smoke tests, single endpoint |
| **vegeta** | Constant rate attack | Throughput at fixed RPS, latency percentiles |
| **wrk** | Multi-threaded | High concurrency, raw throughput |
| **k6** | Scripted scenarios | Complex flows, ramping VUs, thresholds |
| **benchstat** | Statistical analysis | Comparing Go benchmark results |
