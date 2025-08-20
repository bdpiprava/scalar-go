package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	scalargo "github.com/bdpiprava/scalar-go"
)

// This example demonstrates how to integrate scalar-go with different HTTP server
// frameworks and patterns. It shows various integration approaches including:
// - Multiple documentation endpoints
// - URL-based spec loading  
// - Metadata customization
// - Server configuration with timeouts
func main() {
	// Configure HTTP server with proper timeouts
	server := &http.Server{
		Addr:         ":8084",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Setup different integration patterns
	setupStaticDocumentation()
	setupDynamicDocumentation()
	setupURLBasedDocumentation()
	setupMultipleAPIs()
	setupHealthAndStatus()

	fmt.Println("HTTP Server Integration examples:")
	fmt.Println("  Static documentation:")
	fmt.Println("    http://localhost:8084/docs - Basic static docs")
	fmt.Println("  Dynamic documentation:")
	fmt.Println("    http://localhost:8084/api/docs - Dynamic with metadata")
	fmt.Println("  URL-based documentation:")
	fmt.Println("    http://localhost:8084/external/docs - Load from external URL")
	fmt.Println("  Multiple APIs:")
	fmt.Println("    http://localhost:8084/v1/docs - API version 1 docs")
	fmt.Println("    http://localhost:8084/v2/docs - API version 2 docs")
	fmt.Println("  Utility endpoints:")
	fmt.Println("    http://localhost:8084/health - Health check")
	fmt.Println("    http://localhost:8084/status - Service status")
	fmt.Println("Press Ctrl+C to stop the server")

	log.Fatal(server.ListenAndServe())
}

// setupStaticDocumentation creates a simple static documentation endpoint
func setupStaticDocumentation() {
	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		// Simple static documentation generation
		content, err := scalargo.New("./data/loader")
		if err != nil {
			http.Error(w, "Failed to generate documentation", http.StatusInternalServerError)
			log.Printf("Documentation generation error: %v", err)
			return
		}
		
		// Set proper headers
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
		
		w.Write([]byte(content))
	})
}

// setupDynamicDocumentation creates documentation with dynamic metadata
func setupDynamicDocumentation() {
	http.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		// Generate documentation with custom metadata
		content, err := scalargo.NewV2(
			scalargo.WithSpecDir("./data/loader"),
			// Add custom metadata based on request or environment
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Pet Store API Documentation"),
				scalargo.WithKeyValue("version", "1.0.0"),
				scalargo.WithKeyValue("environment", getEnvironment()),
				scalargo.WithKeyValue("generated", time.Now().Format(time.RFC3339)),
			),
			// Customize the appearance
			scalargo.WithTheme(scalargo.ThemeMoon),
			scalargo.WithLayout(scalargo.LayoutModern),
		)
		if err != nil {
			http.Error(w, "Failed to generate API documentation", http.StatusInternalServerError)
			log.Printf("API documentation error: %v", err)
			return
		}
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(content))
	})
}

// setupURLBasedDocumentation demonstrates loading specs from external URLs
func setupURLBasedDocumentation() {
	http.HandleFunc("/external/docs", func(w http.ResponseWriter, r *http.Request) {
		// Load specification from an external URL
		// This is useful for documentation services or when specs are hosted elsewhere
		content, err := scalargo.NewV2(
			scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("External API Documentation"),
				scalargo.WithKeyValue("source", "External URL"),
			),
			scalargo.WithTheme(scalargo.ThemePurple),
		)
		if err != nil {
			http.Error(w, "Failed to load external documentation", http.StatusInternalServerError)
			log.Printf("External documentation error: %v", err)
			return
		}
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(content))
	})
}

// setupMultipleAPIs demonstrates serving documentation for multiple API versions
func setupMultipleAPIs() {
	// API Version 1 documentation
	http.HandleFunc("/v1/docs", func(w http.ResponseWriter, r *http.Request) {
		content, err := scalargo.NewV2(
			scalargo.WithSpecDir("./data/loader"),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Pet Store API v1"),
				scalargo.WithKeyValue("version", "1.0.0"),
				scalargo.WithKeyValue("deprecated", "false"),
			),
			scalargo.WithTheme(scalargo.ThemeDefault),
		)
		if err != nil {
			http.Error(w, "Failed to generate v1 documentation", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(content))
	})

	// API Version 2 documentation (using external spec as example)
	http.HandleFunc("/v2/docs", func(w http.ResponseWriter, r *http.Request) {
		content, err := scalargo.NewV2(
			scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Pet Store API v2"),
				scalargo.WithKeyValue("version", "2.0.0"),
				scalargo.WithKeyValue("deprecated", "false"),
			),
			scalargo.WithTheme(scalargo.ThemeSolarized),
			scalargo.WithDarkMode(),
		)
		if err != nil {
			http.Error(w, "Failed to generate v2 documentation", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(content))
	})
}

// setupHealthAndStatus adds utility endpoints for monitoring
func setupHealthAndStatus() {
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"service":   "scalar-go-docs",
		}
		json.NewEncoder(w).Encode(response)
	})

	// Status endpoint with more detailed information
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"service":     "scalar-go-docs",
			"version":     "1.0.0",
			"environment": getEnvironment(),
			"uptime":      time.Since(startTime).String(),
			"endpoints": map[string]string{
				"docs":          "/docs",
				"api_docs":      "/api/docs", 
				"external_docs": "/external/docs",
				"v1_docs":       "/v1/docs",
				"v2_docs":       "/v2/docs",
			},
		}
		json.NewEncoder(w).Encode(response)
	})
}

// Global variable to track service start time
var startTime = time.Now()

// getEnvironment returns the current environment (mock implementation)
func getEnvironment() string {
	// In a real application, this might read from environment variables
	// or configuration files
	return "development"
}