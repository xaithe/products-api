package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"products-api/data"
	"strconv"

	"github.com/gorilla/mux"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts , given a logger, creates a products handler
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts returns products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	// fetch products from datastore
	listproducts := data.GetProducts()

	// serialise to JSON
	err := listproducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// AddProduct adds a product to the data store
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

// UpdateProducts updates the product of a given ID
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
}

// KeyProduct is a product key for contexts
type KeyProduct struct{}

// MiddlewareProductValidation validates new or updated products
func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			http.Error(
				rw,
				fmt.Sprintf("Unable to validate product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		context := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(context)

		next.ServeHTTP(rw, r)
	})
}
