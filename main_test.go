package grpc_blocker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockHandler implements http.Handler for testing
type mockHandler struct {
	called bool
}

func (h *mockHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.called = true
}

func TestPlugin_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		path           string
		headers        map[string]string
		expectedStatus int
		expectNextCall bool
	}{
		{
			name: "blocked grpc service",
			config: Config{
				BlockedServices: []string{"blocked.service"},
				EnableLogging:  false,
			},
			path: "/blocked.service/method",
			headers: map[string]string{
				"Content-Type": "application/grpc",
			},
			expectedStatus: http.StatusForbidden,
			expectNextCall: false,
		},
		{
			name: "allowed grpc service",
			config: Config{
				BlockedServices: []string{"blocked.service"},
				EnableLogging:  false,
			},
			path: "/allowed.service/method",
			headers: map[string]string{
				"Content-Type": "application/grpc",
			},
			expectedStatus: http.StatusOK,
			expectNextCall: true,
		},
		{
			name: "non-grpc request",
			config: Config{
				BlockedServices: []string{"blocked.service"},
				EnableLogging:  false,
			},
			path: "/blocked.service/method",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			expectedStatus: http.StatusOK,
			expectNextCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock handler for each test
			mock := &mockHandler{}
			
			// Create plugin instance
			plugin, err := New(context.Background(), mock, &tt.config, "test")
			if err != nil {
				t.Fatalf("Failed to create plugin: %v", err)
			}

			// Create test request
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			
			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			plugin.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Check if next handler was called
			if mock.called != tt.expectNextCall {
				t.Errorf("Expected next handler to be called: %v, got: %v", tt.expectNextCall, mock.called)
			}
		})
	}
}

func TestCreateConfig(t *testing.T) {
	config := CreateConfig()
	
	if config == nil {
		t.Fatal("Expected non-nil config")
	}
	
	if len(config.BlockedServices) != 0 {
		t.Errorf("Expected empty BlockedServices, got %v", config.BlockedServices)
	}
	
	if config.EnableLogging != false {
		t.Errorf("Expected EnableLogging to be false, got %v", config.EnableLogging)
	}
}

func TestIsGRPCRequest(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		expectedResult bool
	}{
		{
			name:           "grpc request",
			contentType:    "application/grpc",
			expectedResult: true,
		},
		{
			name:           "non-grpc request",
			contentType:    "application/json",
			expectedResult: false,
		},
		{
			name:           "empty content type",
			contentType:    "",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			result := isGRPCRequest(req)
			if result != tt.expectedResult {
				t.Errorf("Expected isGRPCRequest to return %v, got %v", tt.expectedResult, result)
			}
		})
	}
}