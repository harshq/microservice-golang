// Package classification for Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/harshq/service/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
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

		err = pr.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Validation failed: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, pr)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
