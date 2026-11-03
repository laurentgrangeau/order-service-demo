# `order-service` — Scion Multi-Agent Live Demo

This directory contains the Go microservice (`order-service`) and `.scion/` project configuration used for **Part 8: Live Demo** in the conference talk ([`slides.md`](../../slides.md)).

## Release v2.4 Blockers (Initial State)

Before the Scion task force runs, `order-service` has 3 intentional release blockers across isolated files so 3 worker agents can fix them concurrently without merge conflicts:

| Agent Name | Template Role (`--type`) | Harness (`--harness-config`) | Target Files | Release Blocker |
| --- | --- | --- | --- | --- |
| `lead-mgr` | `manager` | `antigravity` | `scion/release-v2.4` branch | Orchestrates workers, verifies tests, and rebases/merges branches |
| `sec-patch` | `sec-auditor` | `antigravity` | `internal/checkout/repo.go`<br>`internal/middleware/ratelimit.go` | 1. Raw `fmt.Sprintf` SQL injection in `FindByCustomer` & `CreateOrder`<br>2. Pass-through stub in `RateLimiter.Allow` |
| `grpc-mig` | `backend-go` | `gemini` | `proto/order/v2/order.proto`<br>`internal/api/grpc_server.go` | `OrderV2Server` RPC methods return `ErrUnimplementedRPC` |
| `qa-chaos` | `qa-engineer` | `gemini` | `internal/checkout/race_test.go` | Missing concurrent stress & inventory-conservation `go test -race` suite |

---

## Running the Live Demo

### Option A: Local Workstation Mode (Laptop / Docker)

```bash
cd demo/order-service

# 1. Start the local Scion workstation server (if not already running)
scion server start

# 2. Verify templates are discovered
scion templates list

# 3. Dispatch the 4-agent Task Force using the 'local' profile
scion start lead-mgr  "Coordinate v2.4 release sprint"                  --type manager     --harness-config antigravity --profile local
scion start sec-patch "Fix SQLi in checkout & add rate limit"           --type sec-auditor --harness-config antigravity --profile local
scion start grpc-mig  "Migrate REST order handlers to gRPC v2"          --type backend-go  --harness-config gemini      --profile local
scion start qa-chaos  "Write go test -race concurrency suite"           --type qa-engineer --harness-config gemini      --profile local
```

### Option B: GCP High Availability (HA) Hub Mode

```bash
cd demo/order-service

# 1. Authenticate & switch to your GCP HA Hub
scion hub add gcp-ha https://scion.corp.example.com
scion hub auth login
scion use --hub gcp-ha
scion hub list --check
scion hub brokers list

# 2. Dispatch the 4-agent Task Force on GKE Autopilot using '--profile gcp-ha'
scion start lead-mgr  "Coordinate v2.4 release sprint"                  --type manager     --harness-config antigravity --profile gcp-ha
scion start sec-patch "Fix SQLi in checkout & add rate limit"           --type sec-auditor --harness-config antigravity --profile gcp-ha
scion start grpc-mig  "Migrate REST order handlers to gRPC v2"          --type backend-go  --harness-config gemini      --profile gcp-ha
scion start qa-chaos  "Write go test -race concurrency suite"           --type qa-engineer --harness-config gemini      --profile gcp-ha
```

---

## Live Inspection Commands During the Talk

```bash
# List active agents and statuses
scion list

# Send real-time instructions to a running agent
scion msg qa-chaos "Use in-memory store inside your container for zero flakiness"

# Inspect recent terminal output or attach live to tmux
scion look grpc-mig
scion attach grpc-mig

# View structured lifecycle logs
scion logs sec-patch
scion logs grpc-mig
scion logs qa-chaos
```

---

## Resetting Between Rehearsals

Run `./reset-demo.sh` to delete running demo agents (`lead-mgr`, `sec-patch`, `grpc-mig`, `qa-chaos`) and reset any `scion/*` branches back to the initial unpatched state.
