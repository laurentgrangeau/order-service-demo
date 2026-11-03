package checkout

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInsufficientStock  = errors.New("insufficient inventory stock")
	ErrUnsafeSQLDetected  = errors.New("unsafe raw SQL interpolation detected")
)

// OrderRecord represents a stored checkout order.
type OrderRecord struct {
	ID         string
	CustomerID string
	SKU        string
	Quantity   int
	TotalCents int64
	Status     string
	CreatedAt  time.Time
}

// QueryLogEntry records executed SQL statements and bound parameters for auditing and testing.
type QueryLogEntry struct {
	Query string
	Args  []any
}

// SQLExecutor abstracts SQL execution so repo logic can be unit-tested and audited deterministically.
type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) error
	QueryContext(ctx context.Context, query string, args ...any) ([]OrderRecord, error)
}

// InMemoryDB is an in-memory SQL executor & store used for local testing and demo execution.
type InMemoryDB struct {
	mu        sync.Mutex
	orders    map[string]OrderRecord
	stock     map[string]int
	queryLog  []QueryLogEntry
}

// NewInMemoryDB initializes an InMemoryDB with default SKU stock levels.
func NewInMemoryDB(initialStock map[string]int) *InMemoryDB {
	stockCopy := make(map[string]int, len(initialStock))
	for k, v := range initialStock {
		stockCopy[k] = v
	}
	return &InMemoryDB{
		orders: make(map[string]OrderRecord),
		stock:  stockCopy,
	}
}

// LastQueries returns a copy of executed SQL queries for security/parameterization tests.
func (db *InMemoryDB) LastQueries() []QueryLogEntry {
	db.mu.Lock()
	defer db.mu.Unlock()
	out := make([]QueryLogEntry, len(db.queryLog))
	copy(out, db.queryLog)
	return out
}

func (db *InMemoryDB) ExecContext(ctx context.Context, query string, args ...any) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.queryLog = append(db.queryLog, QueryLogEntry{Query: query, Args: args})
	return nil
}

func (db *InMemoryDB) QueryContext(ctx context.Context, query string, args ...any) ([]OrderRecord, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.queryLog = append(db.queryLog, QueryLogEntry{Query: query, Args: args})

	var results []OrderRecord
	for _, ord := range db.orders {
		results = append(results, ord)
	}
	return results, nil
}

// Repository manages order persistence and inventory reservation during checkout.
type Repository struct {
	mu sync.RWMutex
	db *InMemoryDB
}

// NewRepository creates a checkout Repository backed by InMemoryDB.
func NewRepository(db *InMemoryDB) *Repository {
	return &Repository{db: db}
}

// FindByCustomer retrieves orders for a customer.
//
// BLOCKER(sec-patch): This method currently uses fmt.Sprintf string concatenation to build
// the SQL query, creating a SQL injection vulnerability (G201/G202 in gosec).
// Replace with a parameterized query ($1 placeholder) and input validation.
func (r *Repository) FindByCustomer(ctx context.Context, customerID string) ([]OrderRecord, error) {
	if strings.TrimSpace(customerID) == "" {
		return nil, ErrInvalidInput
	}

	// VULNERABLE: Raw string interpolation allows SQL injection payloads like `' OR '1'='1`
	rawQuery := fmt.Sprintf("SELECT id, customer_id, sku, quantity, total_cents, status FROM orders WHERE customer_id = '%s' ORDER BY created_at DESC", customerID)
	if _, err := r.db.QueryContext(ctx, rawQuery); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []OrderRecord
	for _, ord := range r.db.orders {
		if ord.CustomerID == customerID {
			out = append(out, ord)
		}
	}
	return out, nil
}

// CreateOrder reserves SKU stock and persists an OrderRecord.
//
// BLOCKER(sec-patch): Uses fmt.Sprintf for INSERT statement instead of parameterized placeholders ($1..$6).
func (r *Repository) CreateOrder(ctx context.Context, rec OrderRecord) (OrderRecord, error) {
	if rec.ID == "" || rec.CustomerID == "" || rec.SKU == "" || rec.Quantity <= 0 {
		return OrderRecord{}, ErrInvalidInput
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	available := r.db.stock[rec.SKU]
	if available < rec.Quantity {
		return OrderRecord{}, ErrInsufficientStock
	}
	r.db.stock[rec.SKU] = available - rec.Quantity

	if rec.Status == "" {
		rec.Status = "CONFIRMED"
	}
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = time.Now().UTC()
	}

	// VULNERABLE: Raw string concatenation in INSERT statement
	insertSQL := fmt.Sprintf(
		"INSERT INTO orders (id, customer_id, sku, quantity, total_cents, status) VALUES ('%s', '%s', '%s', %d, %d, '%s')",
		rec.ID, rec.CustomerID, rec.SKU, rec.Quantity, rec.TotalCents, rec.Status,
	)
	if err := r.db.ExecContext(ctx, insertSQL); err != nil {
		return OrderRecord{}, err
	}

	r.db.orders[rec.ID] = rec
	return rec, nil
}

// GetOrder fetches a single order by ID.
func (r *Repository) GetOrder(ctx context.Context, id string) (OrderRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ord, ok := r.db.orders[id]
	if !ok {
		return OrderRecord{}, ErrOrderNotFound
	}
	return ord, nil
}

// CancelOrder transitions an order to CANCELLED and restores reserved inventory stock.
func (r *Repository) CancelOrder(ctx context.Context, id string) (OrderRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ord, ok := r.db.orders[id]
	if !ok {
		return OrderRecord{}, ErrOrderNotFound
	}
	if ord.Status != "CANCELLED" {
		r.db.stock[ord.SKU] += ord.Quantity
		ord.Status = "CANCELLED"
		r.db.orders[id] = ord
	}
	return ord, nil
}

// RemainingStock returns the current available stock for a SKU.
func (r *Repository) RemainingStock(sku string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.db.stock[sku]
}
