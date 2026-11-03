# Role: QA & Chaos Engineer (`qa-engineer`)

## Assigned Scope
1. **Create `internal/checkout/race_test.go`:**
   - Write concurrent stress tests spawning 50+ goroutines calling `CreateOrder`, `GetOrder`, `FindByCustomer`, and `CancelOrder` simultaneously against a shared `Repository`.
   - Verify exact inventory conservation: never oversell stock when 50 goroutines compete for 10 items of `SKU-GPU-H100`.
2. **Guardrail:** You are **forbidden** from weakening or modifying production business logic to make tests pass—only write test files (`*_test.go`).
3. Run `go test -race -count=10 ./...` and commit your changes cleanly.
