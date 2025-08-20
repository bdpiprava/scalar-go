package main

import (
	"fmt"
	"log"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
)

// This example demonstrates loading OpenAPI specifications from multiple files
// organized in a structured directory layout. This is useful for large APIs
// where you want to split schemas, paths, and responses into separate files.
//
// Expected directory structure:
// /data/loader-multiple-files/
//     ├── api.yml            // main OpenAPI spec file
//     ├── schemas/           // schema definitions
//     │   ├── Pet.yml
//     │   ├── Pets.yaml
//     │   └── Error.yaml
//     ├── paths/             // API endpoint definitions
//     │   ├── pets.yaml
//     │   └── pet-by-id.yml
//     └── responses/         // response definitions
//         └── Error.yaml
func main() {
	// Directory containing the segmented OpenAPI specification
	// The library automatically combines all files into a single spec
	specDir := "./data/loader-multiple-files"

	// Use NewV2 for more explicit configuration
	// WithSpecDir specifies the directory containing the spec files
	// WithBaseFileName specifies the main spec file (default: "api.yaml")
	htmlContent, err := scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName("api.yml"), // Main spec file name
	)
	if err != nil {
		log.Fatalf("Failed to load multi-file specification: %v", err)
	}

	// Alternative approach using the legacy New function
	// This is equivalent to the above but uses the older API
	// htmlContent, err := scalargo.New(specDir, scalargo.WithBaseFileName("api.yml"))

	// Set up HTTP server to serve the documentation
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})

	// Provide different endpoints to demonstrate the same content
	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})

	fmt.Println("Multi-file API documentation available at:")
	fmt.Println("  http://localhost:8081/")
	fmt.Println("  http://localhost:8081/docs")
	fmt.Println("Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}