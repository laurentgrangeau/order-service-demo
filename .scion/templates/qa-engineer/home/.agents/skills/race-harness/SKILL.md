---
name: race-harness
description: Patterns for writing deterministic Go race-detector and inventory-conservation stress tests.
---

# Go Concurrency & Race Test Patterns

1. Use `sync.WaitGroup` and a start barrier (`close(startCh)`) so all worker goroutines hit `Repository.CreateOrder` at the exact same instant.
2. Assert that `successCount * quantityPerOrder + repo.RemainingStock(sku) == initialStock`.
3. Execute `go test -race ./internal/checkout/...` to verify zero data races.
