package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Response struct สำหรับ JSON response
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// helloHandler - endpoint สำหรับ /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Hello, World!",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(response)
}

// healthHandler - endpoint สำหรับ health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{
		Message: "Service is healthy",
		Status:  "ok",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/health", healthHandler)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
