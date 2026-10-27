package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eduardobombonatto/sezzle-calculator/backend/internal/api"
	"github.com/eduardobombonatto/sezzle-calculator/backend/internal/calculator"
)

const (
	defaultPort         = "8080"
	readHeaderTimeout   = 5 * time.Second
	readTimeout         = 10 * time.Second
	writeTimeout        = 10 * time.Second
	idleTimeout         = 60 * time.Second
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	calcService := calculator.NewService()
	apiHandler := api.NewHandler(calcService)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           apiHandler.RegisterRoutes(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	log.Printf("Calculator API server listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
