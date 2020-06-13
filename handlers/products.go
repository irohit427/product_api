package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/irohit427/coffee-shop/data"
)

type Products struct {
	l *log.Logger
}

func ListProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
	}

	if r.Method == http.MethodPut {
		reg := regexp.MustCompile(`/([0-9]+)`)
		path := r.URL.Path
		g := reg.FindAllStringSubmatch(path, -1)

		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(g[0][1])
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
		}

		p.updateProduct(id, rw, r)
	}
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()
	err := products.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}

	err = data.UpdateProduct(id, prod)

	if err != nil {
		http.Error(rw, "Product Not Found", http.StatusInternalServerError)
		return
	}
}
