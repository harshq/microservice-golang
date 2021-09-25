package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/harshq/service/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Invalid JSON data", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")
	pr := &data.Product{}
	err := pr.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	data.AddProduct(pr)
}

func (p *Products) EditProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT product")

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Invalid parameter", http.StatusInternalServerError)
		return
	}

	pr := r.Context().Value(KeyProduct{}).(*data.Product)
	pr.ID = id

	err = data.UpdateProduct(pr)
	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		pr := &data.Product{}
		err := pr.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Invalid JSON data", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, pr)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
