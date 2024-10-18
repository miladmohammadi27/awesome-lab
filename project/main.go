package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define a counter metric to count the number of requests to the mirror handler
var mirrorRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "mirror_requests_total",
		Help: "Total number of requests to the mirror handler",
	},
	[]string{"status"},
)

func init() {
	// Register the metric with Prometheus
	prometheus.MustRegister(mirrorRequestsTotal)
}

func main() {
	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	// Start the Prometheus metrics server on port 9095
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Starting metrics server on port 9095...")
		log.Fatal(http.ListenAndServe(":9095", nil))
	}()

	http.HandleFunc("/mirror", mirrorHandler)

	fmt.Printf("Starting server on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func mirrorHandler(w http.ResponseWriter, r *http.Request) {
	// Get the 'message' query parameter
	message := r.URL.Query().Get("message")

	if message == "" {
		// Increment the counter for failed requests
		mirrorRequestsTotal.WithLabelValues("failed").Inc()

		// If no 'message' parameter is provided, return a default response
		http.Error(w, "Please provide a message query parameter", http.StatusBadRequest)
		return
	}

	// Increment the counter for successful requests
	mirrorRequestsTotal.WithLabelValues("success").Inc()

	// Create a response object
	response := map[string]string{"mirrored_message": message}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	json.NewEncoder(w).Encode(response)
}
