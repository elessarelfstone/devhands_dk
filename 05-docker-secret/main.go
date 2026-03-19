package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var (
	apiKey         string
	deploymentType string
)

func main() {
	// Read environment variables
	apiKey = os.Getenv("API_KEY")
	deploymentType = os.Getenv("DEPLOYMENT_TYPE")
	endpoint := os.Getenv("ENDPOINT_PATH")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Validate required env vars
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}

	// Set defaults
	if endpoint == "" {
		endpoint = "/v1/deployment-type"
	}
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "8080"
	}

	// Register endpoint handler
	http.HandleFunc(endpoint, deploymentTypeHandler)

	// Start server
	addr := host + ":" + port
	log.Printf("Server starting on %s%s", addr, endpoint)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func deploymentTypeHandler(w http.ResponseWriter, r *http.Request) {
	// Verify request method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Check API key
	clientKey := r.Header.Get("X-API-Key")
	if clientKey != apiKey {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Invalid API key",
			"status": http.StatusUnauthorized,
		})
		return
	}

	// Successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deployment_type": deploymentType,
		"status":          http.StatusOK,
	})
}
