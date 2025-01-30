// main.go
package grpc_blocker

import (
    "context"
    "net/http"
    "strings"
    "log"
)


// Config defines the plugin configuration
type Config struct {
    BlockedServices []string `json:"blockedServices,omitempty"`
    EnableLogging bool `json:"enableLogging,omitempty"`
}

// CreateConfig creates the default plugin configuration
func CreateConfig() *Config {
    return &Config{
        BlockedServices: []string{},
        EnableLogging: false,
    }
}

// Plugin holds the plugin instance configuration
type Plugin struct {
    next            http.Handler
    blockedServices []string
    enableLogging   bool
}

// New creates a new plugin instance
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
    return &Plugin{
        next:            next,
        blockedServices: config.BlockedServices,
        enableLogging: config.EnableLogging,
    }, nil
}

// ServeHTTP handles the HTTP requests
func (p *Plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    // Check if it's a gRPC request
    if isGRPCRequest(req) {
        servicePath := strings.TrimPrefix(req.URL.Path, "/")

        if p.enableLogging {
            log.Printf("[DEBUG] gRPC request received - Full path: %s", servicePath)
        }
        
        parts := strings.Split(servicePath, "/")
        if len(parts) >= 1 {
            serviceName := parts[0]
            if p.enableLogging {
                log.Printf("[DEBUG] Extracted service name: %s, blocked services: %v", serviceName, p.blockedServices)
            }

            for _, blocked := range p.blockedServices {
                if blocked == serviceName {
                    if p.enableLogging {
                        log.Printf("[DEBUG] Blocking request to service: %s", serviceName)
                    }
                    http.Error(rw, "This gRPC service is blocked", http.StatusForbidden)
                    return
                }
            }
        }
    }
    p.next.ServeHTTP(rw, req)
}


// isGRPCRequest checks if the request is a gRPC request
func isGRPCRequest(req *http.Request) bool {
    return strings.Contains(req.Header.Get("Content-Type"), "application/grpc")
}