# Role: Backend Go Engineer (`backend-go`)

## Assigned Scope
1. **Expand `proto/order/v2/order.proto`:**
   - Ensure `OrderService` v2 defines the complete RPC surface (`CreateOrder`, `GetOrder`, `CancelOrder`, `ListCustomerOrders`) and message types.
2. **Implement v2 Handlers in `internal/api/grpc_server.go`:**
   - Replace `ErrUnimplementedRPC` stubs on `OrderV2Server` (`CreateOrder`, `GetOrder`, `CancelOrder`, `ListCustomerOrders`) with full implementations delegating to `checkout.Repository`.
   - Add table-driven tests in `internal/api/grpc_server_test.go`.
3. Do not modify `internal/middleware/ratelimit.go` (owned by `sec-patch`).
4. Run `go test -v ./...` and commit your changes cleanly.
