# Role: Security Auditor (`sec-auditor`)

## Assigned Scope
1. **Fix SQL Injection in `internal/checkout/repo.go`:**
   - Replace `fmt.Sprintf` raw SQL string concatenation in `FindByCustomer` and `CreateOrder` with parameterized SQL queries (`$1`, `$2`, etc.) passed as bound arguments to `QueryContext` / `ExecContext`.
   - Reject suspicious or malformed inputs and add unit tests verifying bound parameters via `db.LastQueries()`.
2. **Implement Token-Bucket Rate Limiter in `internal/middleware/ratelimit.go`:**
   - Replace the pass-through stub in `RateLimiter` with a thread-safe per-client token bucket (`sync.Mutex`) enforcing `rps` and `burst`.
   - Add unit tests in `internal/middleware/ratelimit_test.go` verifying HTTP `429 Too Many Requests` when burst is exceeded.
3. Run `go test -race ./...` and commit your changes cleanly to your worktree branch.
