package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/harshq/service/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// init logger
	l := log.New(os.Stdout, "product-service ", log.LstdFlags)
	// load env file
	err := godotenv.Load(".env")
	if err != nil {
		l.Printf("Error loading .env: %s \n", err)
		os.Exit(1)
	}

	// new servermuliplexer
	sm := mux.NewRouter()

	ph := handlers.NewProducts(l)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.Use(ph.MiddlewareProductValidation)
	postRouter.HandleFunc("/products", ph.AddProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.Use(ph.MiddlewareProductValidation)
	putRouter.HandleFunc("/products/{id:[0-9]+}", ph.EditProduct)

	// server instance
	s := &http.Server{
		Addr:         os.Getenv("BIND_ADDRESS"),
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// starting server on a go routine
	go func() {
		l.Println("Starting server on port 9000")
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s \n", err)
			os.Exit(1)
		}
	}()

	// exit chan
	sigChan := make(chan os.Signal)

	// push to sigchan on kill command
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// blocking till recieved a message
	sig := <-sigChan

	// exit message
	l.Println("Terminate signal received, Gracefully exiting...", sig)

	// graceful shutdown with 30 sec grace period
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
