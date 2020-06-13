package main

import (
	"github.com/gorilla/mux"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/irohit427/product_api/handlers"
)

func main() {
	l := log.New(os.Stdout, "rest-api", log.LstdFlags)

	products := handlers.ListProducts(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", products.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", products.UpdateProduct)
	putRouter.Use(products.ValidateRequestMiddleware)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", products.AddProduct)
	postRouter.Use(products.ValidateRequestMiddleware)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	server.ListenAndServe()
}
