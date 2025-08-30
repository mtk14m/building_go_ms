package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github/mtk14minou/product-service/data"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

type KeyProduct struct{}

func NewProducts(l *log.Logger) *Products {
	return &Products{l: l}
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	listOfProducts := data.GetProducts()
	err := listOfProducts.ToJSON(w)
	if err != nil {
		p.l.Println("[ERROR] unable to marshal json:", err)
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
	p.l.Println("Get all products")
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")
	val := r.Context().Value(KeyProduct{})
	product, ok := val.(*data.Product)
	if !ok {
		p.l.Println("[ERROR] product not found in context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	product.CreatedOn = time.Now().UTC().String()
	product.UpdatedOn = product.CreatedOn
	p.l.Printf("product: %#v", product)
	data.AddProduct(product)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id, error in URI path", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT Products", id)
	val := r.Context().Value(KeyProduct{})
	product, ok := val.(*data.Product)
	if !ok {
		p.l.Println("[ERROR] product not found in context")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	product.ID = id
	product.UpdatedOn = time.Now().UTC().String()
	err = data.UpdateProduct(id, product)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	p.l.Printf("product: %#v", product)
}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := &data.Product{}
		err := product.FromJSON(r.Body)
		if err != nil {
			p.l.Println("[ERROR] unable to unmarshal json:", err)
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
