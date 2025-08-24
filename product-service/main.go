package main

import (
	"context"
	"fmt"
	"github/mtk14minou/product-service/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("product-api running in port 9090")
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	productHandler := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/", productHandler)

	//setting timeout

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	//run serve go routine

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}

	}()

	//gracefull shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	sig := <-sigChan
	l.Println("Received terminate, gracefull shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
