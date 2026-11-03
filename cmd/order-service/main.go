package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/laurentgrangeau/conf-scion/demo/order-service/internal/api"
	"github.com/laurentgrangeau/conf-scion/demo/order-service/internal/checkout"
	"github.com/laurentgrangeau/conf-scion/demo/order-service/internal/middleware"
)

func main() {
	db := checkout.NewInMemoryDB(map[string]int{
		"SKU-GPU-H100": 50,
		"SKU-TPU-V5P":  25,
	})
	repo := checkout.NewRepository(db)
	rest := api.NewLegacyRESTServer(repo)
	limiter := middleware.NewRateLimiter(10, 20)

	mux := http.NewServeMux()
	mux.Handle("/v1/orders", limiter.Wrap(http.HandlerFunc(rest.HandleCreateOrder)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	fmt.Printf("order-service listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
