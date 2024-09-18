package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	requestCount int
	mutex        sync.Mutex
)

func main() {
	http.HandleFunc("/", handler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Increment request count in a thread-safe manner
	mutex.Lock()
	requestCount++
	currentCount := requestCount
	mutex.Unlock()

	// Add 1-second latency every 3rd request
	if currentCount%3 == 0 {
		fmt.Println("Adding 1-second latency...")
		time.Sleep(1 * time.Second)
	}

	w.WriteHeader(200)

	// Send response
	fmt.Fprintf(w, "Request number: %d\n", currentCount)
	fmt.Println("Served request number:", currentCount)

}

