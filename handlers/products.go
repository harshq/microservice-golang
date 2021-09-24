package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/harshq/service/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.editProduct(rw, r)
		return
	}

	http.Error(rw, "Unsupported method!", http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter) {
	p.l.Println("Handle GET products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Invalid JSON data", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")
	pr := &data.Product{}
	err := pr.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	data.AddProduct(pr)
}

func (p *Products) editProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT product")
	url := r.URL.Path
	reg := regexp.MustCompile(`/([0-9]+)`)

	g := reg.FindAllStringSubmatch(url, -1)

	if len(g) > 1 {
		http.Error(rw, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(g[0][1])
	if err != nil {
		http.Error(rw, "Invalid parameter", http.StatusInternalServerError)
		return
	}

	pr := &data.Product{}
	err = pr.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	pr.ID = id

	err = data.UpdateProduct(pr)
	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

}
