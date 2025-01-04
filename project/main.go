package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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

// Define a histogram to measure response times
var mirrorRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "mirror_request_duration_seconds",
		Help:    "Histogram of request durations for mirror handler",
		Buckets: prometheus.DefBuckets, // Use default bucket values
	},
	[]string{"status"},
)

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(mirrorRequestsTotal)
	prometheus.MustRegister(mirrorRequestDuration)
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
	fmt.Println("Received request")
	start := time.Now() // Record the start time

	// Introduce a random wait (simulate processing time)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	// Get the 'message' query parameter
	message := r.URL.Query().Get("message")

	status := "success"
	if message == "" {
		status = "failed"
		// Increment the counter for failed requests
		mirrorRequestsTotal.WithLabelValues(status).Inc()

		// If no 'message' parameter is provided, return a default response
		http.Error(w, "Please provide a message query parameter", http.StatusBadRequest)
	} else {
		// Increment the counter for successful requests
		mirrorRequestsTotal.WithLabelValues(status).Inc()

		// Create a response object
		response := map[string]string{"mirrored_message": message}

		// Set the content type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON response
		json.NewEncoder(w).Encode(response)
	}

	// Measure the request duration and observe it in the histogram
	duration := time.Since(start).Seconds()
	mirrorRequestDuration.WithLabelValues(status).Observe(duration)
	fmt.Println("Request duration:", duration)
}
