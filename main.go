package main

import (
	"context"
	"fmt"
	"log"
	"naevis/search"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Client

// Security headers middleware
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set HTTP headers for enhanced security
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store, must-revalidate, private")
		next.ServeHTTP(w, r) // Call the next handler
	})
}

func main() {
	baseRouter := httprouter.New()
	baseRouter.GET("/health", Index)
	baseRouter.GET("/ac", search.Autocompleter)
	baseRouter.GET("/api/search/:entityType", search.SearchHandler)
	baseRouter.POST("/emitted", search.EventHandler)

	// Apply middleware only once at setup
	handler := securityHeaders(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(baseRouter))

	server := &http.Server{
		Addr:              ":7000",
		Handler:           handler,
		ReadTimeout:       5 * time.Second,   // Prevents slow client attacks
		WriteTimeout:      10 * time.Second,  // Prevents slow responses
		IdleTimeout:       120 * time.Second, // Keeps connections open for reuse
		ReadHeaderTimeout: 2 * time.Second,   // Limits header parsing time
	}

	// Start server in a goroutine to handle graceful shutdown
	go func() {
		log.Println("Server started on port 7000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port 7000: %v", err)
		}
	}()

	// Graceful shutdown listener
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	<-shutdownChan
	log.Println("Shutting down gracefully...")

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped")
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "200")
}
