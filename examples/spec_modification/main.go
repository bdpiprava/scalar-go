package main

import (
	"fmt"
	"log"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/model"
)

// This example demonstrates how to dynamically modify OpenAPI specifications
// before generating the documentation. This is useful for:
// - Adding or modifying API information at runtime
// - Customizing server URLs based on environment
// - Adding dynamic tags or descriptions
// - Filtering or modifying paths and operations
func main() {
	apiDir := "./data/loader"

	// Setup different examples of spec modification
	setupBasicModification(apiDir)
	setupServerModification(apiDir)
	setupDynamicInfo(apiDir)
	setupPathModification(apiDir)

	fmt.Println("Spec modification examples available at:")
	fmt.Println("  Basic modification:")
	fmt.Println("    http://localhost:8083/basic - Simple title and description change")
	fmt.Println("  Server modification:")
	fmt.Println("    http://localhost:8083/servers - Dynamic server URLs")
	fmt.Println("  Dynamic info:")
	fmt.Println("    http://localhost:8083/dynamic - Runtime information updates")
	fmt.Println("  Path modification:")
	fmt.Println("    http://localhost:8083/paths - Show documented paths count")
	fmt.Println("Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// setupBasicModification demonstrates basic spec information changes
func setupBasicModification(apiDir string) {
	http.HandleFunc("/basic", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalargo.New(apiDir,
			// WithSpecModifier allows you to modify the spec before rendering
			scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
				// Modify the API title and description
				spec.Info.Title = "Modified Pet Store API"
				spec.Info.Description = "This API specification has been dynamically modified to show custom title and description."
				spec.Info.Version = "2.0.0-modified"
				
				return spec
			}),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})
}

// setupServerModification demonstrates dynamic server URL modification
func setupServerModification(apiDir string) {
	http.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalargo.New(apiDir,
			scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
				// Add dynamic server URLs based on environment or request
				spec.Servers = []model.Server{
					{
						URL:         "https://api.example.com/v1",
						Description: "Production server",
					},
					{
						URL:         "https://staging-api.example.com/v1", 
						Description: "Staging server",
					},
					{
						URL:         "http://localhost:8080/api/v1",
						Description: "Local development server",
					},
				}
				
				return spec
			}),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})
}

// setupDynamicInfo demonstrates runtime information updates
func setupDynamicInfo(apiDir string) {
	http.HandleFunc("/dynamic", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalargo.New(apiDir,
			scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
				// Add dynamic information based on current state
				originalTitle := spec.Info.Title
				spec.Info.Title = fmt.Sprintf("%s (Generated at Runtime)", originalTitle)
				
				// Add current timestamp to description
				spec.Info.Description = fmt.Sprintf(
					"%s\n\n**Note:** This documentation was generated dynamically and includes runtime modifications.", 
					spec.Info.Description,
				)
				
				// Add a custom tag with dynamic information
				spec.Tags = append(spec.Tags, model.Tag{
					Name:        "runtime-info",
					Description: "This tag was added dynamically during spec modification",
				})
				
				return spec
			}),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})
}

// setupPathModification demonstrates working with API paths
func setupPathModification(apiDir string) {
	http.HandleFunc("/paths", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalargo.New(apiDir,
			scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
				// Get information about documented paths
				documentedPaths := spec.DocumentedPaths()
				
				// Update the description to include path information
				pathInfo := fmt.Sprintf("\n\n**API Statistics:**\n- Total endpoints: %d\n", len(documentedPaths))
				
				// Add details about each path
				pathInfo += "\n**Available Endpoints:**\n"
				for _, path := range documentedPaths {
					pathInfo += fmt.Sprintf("- %s %s\n", path.Method, path.Path)
				}
				
				spec.Info.Description = spec.Info.Description + pathInfo
				
				return spec
			}),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})
}