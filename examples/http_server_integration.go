package examples

import (
	"encoding/json"
	"net/http"
	"time"

	scalargo "github.com/bdpiprava/scalar-go"
)

// ExampleStaticDocumentation demonstrates simple static documentation generation
// @example Static Documentation
// @description This example shows how to create simple static documentation with proper headers and caching for production use.
func ExampleStaticDocumentation() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
	)
}

// ExampleDynamicDocumentation demonstrates documentation with dynamic metadata
// @example Dynamic Documentation
// @description This example shows how to generate documentation with custom metadata based on request or environment, including dynamic titles and timestamps.
func ExampleDynamicDocumentation() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Pet Store API Documentation"),
			scalargo.WithKeyValue("version", "1.0.0"),
			scalargo.WithKeyValue("environment", getEnvironment()),
			scalargo.WithKeyValue("generated", time.Now().Format(time.RFC3339)),
		),
		scalargo.WithTheme(scalargo.ThemeMoon),
		scalargo.WithLayout(scalargo.LayoutModern),
	)
}

// ExampleURLBasedDocumentation demonstrates loading specs from external URLs
// @example URL-Based Documentation
// @description This example shows how to load specification from an external URL, useful for documentation services or when specs are hosted elsewhere.
func ExampleURLBasedDocumentation() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("External API Documentation"),
			scalargo.WithKeyValue("source", "External URL"),
		),
		scalargo.WithTheme(scalargo.ThemePurple),
	)
}

// ExampleAPIV1 demonstrates serving documentation for API version 1
// @example API Version 1
// @description This example shows how to serve documentation for multiple API versions with version-specific metadata and theming.
func ExampleAPIV1() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Pet Store API v1"),
			scalargo.WithKeyValue("version", "1.0.0"),
			scalargo.WithKeyValue("deprecated", "false"),
		),
		scalargo.WithTheme(scalargo.ThemeDefault),
	)
}

// ExampleAPIV2 demonstrates serving documentation for API version 2 with external spec
// @example API Version 2
// @description This example shows how to serve documentation for a newer API version using external specifications with enhanced theming and dark mode.
func ExampleAPIV2() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Pet Store API v2"),
			scalargo.WithKeyValue("version", "2.0.0"),
			scalargo.WithKeyValue("deprecated", "false"),
		),
		scalargo.WithTheme(scalargo.ThemeSolarized),
		scalargo.WithDarkMode(),
	)
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
