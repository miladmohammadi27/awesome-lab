package main

import (
	"fmt"
	"net/http"
	"sync"
)

const (
	url            = "http://localhost:9092/mirror?message=HelloWorld"
	requestCount   = 1000 // Total number of requests to send
	goroutineCount = 4    // Number of goroutines
)

func sendRequest(wg *sync.WaitGroup, start int) {
	defer wg.Done()

	for i := start; i < requestCount; i += goroutineCount {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error sending request %d: %v\n", i, err)
		}
		resp.Body.Close() // Close the response body to avoid resource leaks
		fmt.Printf("Request %d sent...\n", i)
	}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go sendRequest(&wg, i)
	}

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("Finished sending requests")
}
