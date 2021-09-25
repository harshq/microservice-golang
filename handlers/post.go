package handlers

import (
	"net/http"

	"github.com/harshq/service/data"
)

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")
	pr := r.Context().Value(KeyProduct{}).(*data.Product)

	data.AddProduct(pr)
}
