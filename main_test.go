package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	// สร้าง request ทดสอบ
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	// เรียกใช้ handler
	helloHandler(w, req)

	// ตรวจสอบ status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// ตรวจสอบ Content-Type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// ตรวจสอบ response body
	var response Response
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedMessage := "Hello, World!"
	if response.Message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, response.Message)
	}

	expectedStatus := "success"
	if response.Status != expectedStatus {
		t.Errorf("Expected status '%s', got '%s'", expectedStatus, response.Status)
	}
}

func TestHealthHandler(t *testing.T) {
	// สร้าง request ทดสอบ
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// เรียกใช้ handler
	healthHandler(w, req)

	// ตรวจสอบ status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// ตรวจสอบ response body
	var response Response
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	expectedMessage := "Service is healthy"
	if response.Message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, response.Message)
	}

	expectedStatus := "ok"
	if response.Status != expectedStatus {
		t.Errorf("Expected status '%s', got '%s'", expectedStatus, response.Status)
	}
}
