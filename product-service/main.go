package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github/mtk14minou/product-service/handlers"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("product-api running on port 9090")
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	productHandler := handlers.NewProducts(l)

	// Create the main router
	sm := mux.NewRouter()

	// GET /products
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetProducts)

	// GET /products/id
	getByIdRouter := sm.Methods(http.MethodGet).Subrouter()
	getByIdRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.GetProductById)

	// POST /products (AddProduct)
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", productHandler.AddProduct)
	postRouter.Use(productHandler.MiddlewareProductValidation)

	// PUT /products/{id} (UpdateProduct)
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProduct)
	putRouter.Use(productHandler.MiddlewareProductValidation)

	// Configure the server
	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Run the server in a goroutine
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.Shutdown(tc)
}
