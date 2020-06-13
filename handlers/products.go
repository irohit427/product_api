package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/irohit427/product_api/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}


func ListProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()
	err := products.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(Prod{}).(*data.Product)
	data.AddProduct(prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Invalid Url", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(Prod{}).(*data.Product)

	err = data.UpdateProduct(id, prod)

	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}
}

type Prod struct {}

func (p *Products) ValidateRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Printf("[ERROR]: ", err)
			http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.l.Printf("[ERROR]: ", err)
			http.Error(
					rw,
					fmt.Sprintf("Error validating product: %s", err),
					http.StatusBadRequest,
				)
			return
		}

		ctx := context.WithValue(r.Context(), Prod{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
