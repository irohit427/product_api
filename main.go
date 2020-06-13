package main

import (
	"log"
	"os"
	"time"

	"net/http"

	"github.com/irohit427/coffee-shop/handlers"
)

func main() {
	l := log.New(os.Stdout, "rest-api", log.LstdFlags)

	products := handlers.ListProducts(l)

	sm := http.NewServeMux()

	sm.Handle("/", products)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	server.ListenAndServe()
}
