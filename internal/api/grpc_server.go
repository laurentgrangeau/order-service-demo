package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/laurentgrangeau/conf-scion/demo/order-service/internal/checkout"
)

var (
	ErrUnimplementedRPC = errors.New("gRPC v2 method not implemented yet")
)

// LegacyRESTServer implements the v1 JSON/REST order endpoints.
// BLOCKER(grpc-mig): Migrate these handlers to the v2 gRPC service surface in grpc_server.go.
type LegacyRESTServer struct {
	repo *checkout.Repository
}

func NewLegacyRESTServer(repo *checkout.Repository) *LegacyRESTServer {
	return &LegacyRESTServer{repo: repo}
}

func (s *LegacyRESTServer) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req checkout.OrderRecord
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ord, err := s.repo.CreateOrder(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ord)
}

// OrderV2Server is the gRPC v2 service implementation corresponding to proto/order/v2/order.proto.
//
// BLOCKER(grpc-mig): Complete the gRPC v2 handlers (CreateOrder, GetOrder, CancelOrder, ListCustomerOrders, etc.)
// with full request validation, status code mapping, and unit tests.
type OrderV2Server struct {
	repo *checkout.Repository
}

func NewOrderV2Server(repo *checkout.Repository) *OrderV2Server {
	return &OrderV2Server{repo: repo}
}

func (s *OrderV2Server) CreateOrder(ctx context.Context, req checkout.OrderRecord) (checkout.OrderRecord, error) {
	return checkout.OrderRecord{}, ErrUnimplementedRPC
}

func (s *OrderV2Server) GetOrder(ctx context.Context, id string) (checkout.OrderRecord, error) {
	return checkout.OrderRecord{}, ErrUnimplementedRPC
}

func (s *OrderV2Server) CancelOrder(ctx context.Context, id string) (checkout.OrderRecord, error) {
	return checkout.OrderRecord{}, ErrUnimplementedRPC
}
