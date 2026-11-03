---
name: grpc-migration
description: Guidelines for migrating legacy REST handlers to idiomatic Go OrderService v2 RPCs.
---

# REST to gRPC v2 Migration Guide

1. Validate all incoming request fields (`ID`, `CustomerID`, `SKU`, `Quantity > 0`) before calling `Repository`.
2. Map repository domain errors (`checkout.ErrOrderNotFound`, `checkout.ErrInvalidInput`, `checkout.ErrInsufficientStock`) cleanly so callers can distinguish validation errors from not-found errors.
3. Cover happy-path and error-path RPC calls in `internal/api/grpc_server_test.go`.
