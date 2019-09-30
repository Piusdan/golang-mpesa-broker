package main

import (
	"log"
	"net/http"
	"time"
	"github.com/Piusdan/payments-gateway/paymentsgateway"
)

func main() {
	router := paymentsgateway.NewRouter()
	server := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

