package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/harshq/service/data"
)

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
