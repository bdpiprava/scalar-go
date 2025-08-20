package main

import (
	"fmt"
	"log"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
)

// This example demonstrates the most basic usage of scalar-go:
// generating HTML documentation from a single OpenAPI specification file.
func main() {
	// Define the directory containing your OpenAPI spec file
	// The library will look for "api.yaml" by default
	apiDir := "./data/loader"

	// Generate HTML documentation from the OpenAPI spec
	// This creates a complete HTML page with embedded Scalar UI
	htmlContent, err := scalargo.New(apiDir)
	if err != nil {
		log.Fatalf("Failed to generate API documentation: %v", err)
	}

	// Set up a simple HTTP server to serve the documentation
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set the content type to HTML
		w.Header().Set("Content-Type", "text/html")
		
		// Write the generated HTML content
		w.Write([]byte(htmlContent))
	})

	// Start the server
	fmt.Println("API documentation available at: http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop the server")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}