package http_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	httpClient "github.com/vixyninja/go-blocks/http"
)

func TestNewHTTPClient(t *testing.T) {
	// Test with default config
	client := httpClient.NewHTTPClient(nil)
	if client == nil {
		t.Fatal("Expected client, got nil")
	}

	// Test with custom config
	config := &httpClient.Config{
		Timeout:   10 * time.Second,
		UserAgent: "Test-Client/1.0",
		DefaultHeaders: map[string]string{
			"X-Test": "value",
		},
	}
	client = httpClient.NewHTTPClient(config)
	if client == nil {
		t.Fatal("Expected client, got nil")
	}
}

func TestHTTPClient_Get(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("Expected /test, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Custom") != "value" {
			t.Errorf("Expected X-Custom: value, got %s", r.Header.Get("X-Custom"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	// Create client
	client := httpClient.NewHTTPClient(nil)

	// Test GET request
	ctx := context.Background()
	headers := map[string]string{"X-Custom": "value"}
	resp, err := client.Get(ctx, server.URL+"/test", headers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Expected no error reading body, got %v", err)
	}
	if string(body) != "test response" {
		t.Errorf("Expected 'test response', got %s", string(body))
	}
}

func TestHTTPClient_Post(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Expected no error reading body, got %v", err)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			t.Fatalf("Expected valid JSON, got error: %v", err)
		}

		if data["name"] != "test" {
			t.Errorf("Expected name: test, got %v", data["name"])
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1, "name": "test"}`))
	}))
	defer server.Close()

	// Create client
	client := httpClient.NewHTTPClient(nil)

	// Test POST request
	ctx := context.Background()
	data := map[string]interface{}{"name": "test"}
	jsonBody, err := httpClient.JSONRequest(data)
	if err != nil {
		t.Fatalf("Expected no error creating JSON body, got %v", err)
	}

	headers := httpClient.JSONHeaders()
	resp, err := client.Post(ctx, server.URL+"/test", jsonBody, headers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}
}

func TestHTTPClient_Put(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("updated"))
	}))
	defer server.Close()

	// Create client
	client := httpClient.NewHTTPClient(nil)

	// Test PUT request
	ctx := context.Background()
	body := strings.NewReader("test data")
	headers := map[string]string{"Content-Type": "text/plain"}

	resp, err := client.Put(ctx, server.URL+"/test", body, headers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHTTPClient_Delete(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Create client
	client := httpClient.NewHTTPClient(nil)

	// Test DELETE request
	ctx := context.Background()
	resp, err := client.Delete(ctx, server.URL+"/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestJSONRequest(t *testing.T) {
	data := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	reader, err := httpClient.JSONRequest(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Expected no error reading body, got %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	if result["name"] != "test" {
		t.Errorf("Expected name: test, got %v", result["name"])
	}
	if result["value"] != float64(123) {
		t.Errorf("Expected value: 123, got %v", result["value"])
	}
}

func TestJSONHeaders(t *testing.T) {
	headers := httpClient.JSONHeaders()

	expected := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	for key, value := range expected {
		if headers[key] != value {
			t.Errorf("Expected %s: %s, got %s", key, value, headers[key])
		}
	}
}

func TestFormHeaders(t *testing.T) {
	headers := httpClient.FormHeaders()

	if headers["Content-Type"] != "application/x-www-form-urlencoded" {
		t.Errorf("Expected Content-Type: application/x-www-form-urlencoded, got %s", headers["Content-Type"])
	}
}

func TestBearerTokenHeaders(t *testing.T) {
	token := "test-token"
	headers := httpClient.BearerTokenHeaders(token)

	expected := "Bearer " + token
	if headers["Authorization"] != expected {
		t.Errorf("Expected Authorization: %s, got %s", expected, headers["Authorization"])
	}
}

func TestReadResponseBody(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response body"))
	}))
	defer server.Close()

	// Make request
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	// Test ReadResponseBody
	body, err := httpClient.ReadResponseBody(resp)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if body != "test response body" {
		t.Errorf("Expected 'test response body', got %s", body)
	}
}

func TestParseJSONResponse(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "name": "test", "active": true}`))
	}))
	defer server.Close()

	// Make request
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	// Test ParseJSONResponse
	type TestData struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Active bool   `json:"active"`
	}

	var data TestData
	if err := httpClient.ParseJSONResponse(resp, &data); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if data.ID != 1 {
		t.Errorf("Expected ID: 1, got %d", data.ID)
	}
	if data.Name != "test" {
		t.Errorf("Expected Name: test, got %s", data.Name)
	}
	if data.Active != true {
		t.Errorf("Expected Active: true, got %v", data.Active)
	}
}

func TestIsSuccess(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{200, true},
		{201, true},
		{299, true},
		{300, false},
		{400, false},
		{500, false},
	}

	for _, test := range tests {
		resp := &http.Response{StatusCode: test.statusCode}
		result := httpClient.IsSuccess(resp)
		if result != test.expected {
			t.Errorf("Expected IsSuccess(%d) = %v, got %v", test.statusCode, test.expected, result)
		}
	}
}

func TestIsClientError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{399, false},
		{400, true},
		{404, true},
		{499, true},
		{500, false},
	}

	for _, test := range tests {
		resp := &http.Response{StatusCode: test.statusCode}
		result := httpClient.IsClientError(resp)
		if result != test.expected {
			t.Errorf("Expected IsClientError(%d) = %v, got %v", test.statusCode, test.expected, result)
		}
	}
}

func TestIsServerError(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{499, false},
		{500, true},
		{503, true},
		{599, true},
		{600, false},
	}

	for _, test := range tests {
		resp := &http.Response{StatusCode: test.statusCode}
		result := httpClient.IsServerError(resp)
		if result != test.expected {
			t.Errorf("Expected IsServerError(%d) = %v, got %v", test.statusCode, test.expected, result)
		}
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		baseURL  string
		params   map[string]string
		expected string
	}{
		{
			"https://api.example.com",
			nil,
			"https://api.example.com",
		},
		{
			"https://api.example.com",
			map[string]string{},
			"https://api.example.com",
		},
		{
			"https://api.example.com",
			map[string]string{"page": "1"},
			"https://api.example.com?page=1",
		},
		{
			"https://api.example.com",
			map[string]string{"page": "1", "limit": "10"},
			"https://api.example.com?page=1&limit=10",
		},
	}

	for _, test := range tests {
		result := httpClient.BuildURL(test.baseURL, test.params)
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestCloseResponse(t *testing.T) {
	// Test with nil response
	httpClient.CloseResponse(nil) // Should not panic

	// Test with valid response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	httpClient.CloseResponse(resp) // Should not panic
}

func TestDefaultConfig(t *testing.T) {
	config := httpClient.DefaultConfig()

	if config.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", config.Timeout)
	}

	if config.UserAgent != "Gaia-Client/1.0" {
		t.Errorf("Expected UserAgent 'Gaia-Client/1.0', got %s", config.UserAgent)
	}

	if config.DefaultHeaders == nil {
		t.Error("Expected DefaultHeaders to be initialized")
	}
}

func TestHTTPClient_WithCustomConfig(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check User-Agent
		if r.Header.Get("User-Agent") != "Custom-Client/2.0" {
			t.Errorf("Expected User-Agent: Custom-Client/2.0, got %s", r.Header.Get("User-Agent"))
		}

		// Check default header
		if r.Header.Get("X-Default") != "default-value" {
			t.Errorf("Expected X-Default: default-value, got %s", r.Header.Get("X-Default"))
		}

		// Check custom header (should override default)
		if r.Header.Get("X-Custom") != "custom-value" {
			t.Errorf("Expected X-Custom: custom-value, got %s", r.Header.Get("X-Custom"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	// Create client with custom config
	config := &httpClient.Config{
		Timeout:   5 * time.Second,
		UserAgent: "Custom-Client/2.0",
		DefaultHeaders: map[string]string{
			"X-Default": "default-value",
			"X-Custom":  "default-custom-value",
		},
	}
	client := httpClient.NewHTTPClient(config)

	// Test request with custom headers (should override default)
	ctx := context.Background()
	headers := map[string]string{"X-Custom": "custom-value"}

	resp, err := client.Get(ctx, server.URL, headers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
