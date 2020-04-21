package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"products-api/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// create handlers
	producthandler := handlers.NewProducts(logger)

	// create new server mux, register handlers
	mux := http.NewServeMux()
	mux.Handle("/", producthandler)

	// create new server
	server := http.Server{
		Addr:         ":9090",
		Handler:      mux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	//http.ListenAndServe(":9090", mux)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Println("Received terminate, graceful shutdown", sig)

	timeout, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(timeout)
}
