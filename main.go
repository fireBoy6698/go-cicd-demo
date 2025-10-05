package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

// helloCatHandler - endpoint สำหรับ /hello-cat
func helloCatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Hello, Cat! 🐱",
		Status:  "success",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// helloDogHandler - endpoint สำหรับ /hello-dog
func helloDogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Hello, Dog! 🐶🐶",
		Status:  "success",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// helloBirdHandler - endpoint สำหรับ /hello-bird
func helloBirdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{
		Message: "Hello, Bird! 🐦🐦 fix bug !!",
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
	http.HandleFunc("/hello-cat", helloCatHandler)
	http.HandleFunc("/hello-dog", helloDogHandler)
	http.HandleFunc("/hello-bird", helloBirdHandler)

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get environment name
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	fmt.Printf("Server starting on port %s (environment: %s)\n", port, env)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
