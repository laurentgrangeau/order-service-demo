package checkout

import (
	"context"
	"testing"
)

func TestRepositoryCreateAndGetOrder(t *testing.T) {
	db := NewInMemoryDB(map[string]int{"SKU-GPU-H100": 10})
	repo := NewRepository(db)

	ctx := context.Background()
	created, err := repo.CreateOrder(ctx, OrderRecord{
		ID:         "ord-100",
		CustomerID: "cust-42",
		SKU:        "SKU-GPU-H100",
		Quantity:   2,
		TotalCents: 600000,
	})
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	if created.Status != "CONFIRMED" {
		t.Fatalf("expected status CONFIRMED, got %s", created.Status)
	}
	if got := repo.RemainingStock("SKU-GPU-H100"); got != 8 {
		t.Fatalf("expected remaining stock 8, got %d", got)
	}

	fetched, err := repo.GetOrder(ctx, "ord-100")
	if err != nil {
		t.Fatalf("GetOrder failed: %v", err)
	}
	if fetched.CustomerID != "cust-42" {
		t.Fatalf("expected customer cust-42, got %s", fetched.CustomerID)
	}
}
