package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type DocumentResponse struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Calls   int64  `json:"calls"`
}

var (
	db  *sql.DB
	rdb *redis.Client
)

func main() {
	// Get port from environment
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	// Database connection
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "host=postgres user=wikiuser password=wikipass dbname=wikidb sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	// Redis connection
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	defer rdb.Close()

	// Test Redis connection
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Redis ping failed: %v", err)
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/document/", documentHandler)
	mux.HandleFunc("/health", healthHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start server
	go func() {
		log.Printf("HTTP server started on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server stopped")
}

func documentHandler(w http.ResponseWriter, r *http.Request) {
	// Extract document ID from URL
	id := r.URL.Path[len("/document/"):]
	if id == "" {
		http.Error(w, "Missing document ID", http.StatusBadRequest)
		return
	}

	// Query database
	var docID int32
	var title, content string
	err := db.QueryRowContext(r.Context(),
		"SELECT id, title, content FROM documents WHERE id = $1",
		id,
	).Scan(&docID, &title, &content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Document not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Increment call count
	count, err := rdb.Incr(r.Context(), "call_count").Result()
	if err != nil {
		log.Printf("Redis error: %v", err)
		count = -1
	}

	// Prepare response
	response := DocumentResponse{
		ID:      docID,
		Title:   title,
		Content: content,
		Calls:   count,
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Check database
	if err := db.PingContext(r.Context()); err != nil {
		http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		return
	}

	// Check Redis
	if _, err := rdb.Ping(r.Context()).Result(); err != nil {
		http.Error(w, "Redis unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
