package main

import (
	"log"
	"net/http"
	"time"

	"go-name/internal/handler"
	"go-name/internal/service"
)

func main() {
	// 1. Inisialisasi Service & Handler
	pdfService := service.NewPDFService()
	pdfHandler := handler.NewPDFHandler(pdfService)

	// 2. Setup Router
	mux := http.NewServeMux()
	mux.HandleFunc("/api/generate", pdfHandler.Generate)

	// 3. Setup Server Configuration
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 4. Start Server
	log.Println("🚀 Server running on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on :8080: %v\n", err)
	}
}
