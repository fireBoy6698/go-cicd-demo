package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//version := "1.0.1"

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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// healthHandler - endpoint สำหรับ health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{
		Message: "Service is healthy",
		Status:  "ok",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// helloGoHandler - endpoint สำหรับ /hello-go
func helloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Hello, Go!",
		Status:  "success",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/hello-go", helloGoHandler)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
