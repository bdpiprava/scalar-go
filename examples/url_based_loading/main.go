package main

import (
	"fmt"
	"log"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
)

// This example demonstrates loading OpenAPI specifications from remote URLs.
// This is useful when:
// - Your API specs are hosted separately (e.g., in a CDN or documentation service)
// - You want to display documentation for external APIs
// - You're building a documentation aggregator
// - Your specs are stored in a central repository
func main() {
	setupExternalSpecExamples()
	setupMultipleExternalAPIs()
	setupWithCustomization()

	fmt.Println("URL-based documentation loading examples:")
	fmt.Println("  External specs:")
	fmt.Println("    http://localhost:8085/scalar-galaxy - Scalar Galaxy API example")
	fmt.Println("    http://localhost:8085/petstore - Classic Petstore API")
	fmt.Println("  Multiple external APIs:")
	fmt.Println("    http://localhost:8085/apis/github - GitHub API documentation")
	fmt.Println("    http://localhost:8085/apis/openai - OpenAI API documentation")
	fmt.Println("  Customized external docs:")
	fmt.Println("    http://localhost:8085/custom/docs - External spec with custom styling")
	fmt.Println("Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(":8085", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// setupExternalSpecExamples demonstrates basic URL-based spec loading
func setupExternalSpecExamples() {
	// Load Scalar Galaxy API specification from CDN
	http.HandleFunc("/scalar-galaxy", func(w http.ResponseWriter, r *http.Request) {
		content, err := scalargo.NewV2(
			// Specify the URL where the OpenAPI spec is hosted
			scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load Scalar Galaxy spec: %v", err), http.StatusInternalServerError)
			log.Printf("Error loading Scalar Galaxy spec: %v", err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(content))
	})

	// Load classic Petstore API specification
	http.HandleFunc("/petstore", func(w http.ResponseWriter, r *http.Request) {
		content, err := scalargo.NewV2(
			// Using the classic Petstore OpenAPI spec from the official repository
			scalargo.WithSpecURL("https://petstore3.swagger.io/api/v3/openapi.json"),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load Petstore spec: %v", err), http.StatusInternalServerError)
			log.Printf("Error loading Petstore spec: %v", err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(content))
	})
}

// setupMultipleExternalAPIs demonstrates loading different external API specifications
func setupMultipleExternalAPIs() {
	// GitHub API documentation (using a public OpenAPI spec)
	http.HandleFunc("/apis/github", func(w http.ResponseWriter, r *http.Request) {
		content, err := scalargo.NewV2(
			scalargo.WithSpecURL("https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.json"),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("GitHub REST API"),
				scalargo.WithKeyValue("description", "Complete GitHub REST API documentation"),
			),
			scalargo.WithTheme(scalargo.ThemeDefault),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load GitHub API spec: %v", err), http.StatusInternalServerError)
			log.Printf("Error loading GitHub API spec: %v", err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(content))
	})

	// OpenAI API documentation (hypothetical - using Scalar Galaxy as example)
	http.HandleFunc("/apis/openai", func(w http.ResponseWriter, r *http.Request) {
		content, err := scalargo.NewV2(
			// Using Scalar Galaxy as a placeholder for demonstration
			scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("OpenAI API (Demo)"),
				scalargo.WithKeyValue("description", "Example of loading external API documentation"),
				scalargo.WithKeyValue("note", "This is a demo using Scalar Galaxy spec"),
			),
			scalargo.WithTheme(scalargo.ThemeAlternate),
			scalargo.WithLayout(scalargo.LayoutClassic),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load OpenAI API spec: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(content))
	})
}

// setupWithCustomization demonstrates URL loading with extensive customization
func setupWithCustomization() {
	http.HandleFunc("/custom/docs", func(w http.ResponseWriter, r *http.Request) {
		// Custom CSS for branding the external documentation
		customCSS := `
			/* Custom branding for external API docs */
			.scalar-api-reference {
				--scalar-color-1: #1a202c;
				--scalar-color-2: #2d3748;
				--scalar-color-3: #4a5568;
				--scalar-color-accent: #3182ce;
			}
			
			.section-header {
				color: #3182ce !important;
				border-bottom: 2px solid #3182ce !important;
			}
			
			/* Custom styling for request panels */
			.api-client__request {
				background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
				color: white !important;
				border-radius: 8px !important;
			}
			
			/* Custom button styling */
			button {
				background: #3182ce !important;
				border: none !important;
				border-radius: 6px !important;
				transition: all 0.2s ease !important;
			}
			
			button:hover {
				background: #2c5282 !important;
				transform: translateY(-1px) !important;
			}
		`

		content, err := scalargo.NewV2(
			// Load external specification
			scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
			
			// Add comprehensive customization
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Customized External API Docs"),
				scalargo.WithKeyValue("company", "Your Company Name"),
				scalargo.WithKeyValue("customized", "true"),
			),
			
			// Apply theme and layout
			scalargo.WithTheme(scalargo.ThemeMoon),
			scalargo.WithLayout(scalargo.LayoutModern),
			
			// Apply custom CSS
			scalargo.WithOverrideCSS(customCSS),
			
			// Configure UI options
			scalargo.WithDarkMode(),
			scalargo.WithHideDownloadButton(),
			scalargo.WithSearchHotKey("cmd+k"),
			
			// Hide specific client examples
			scalargo.WithHiddenClients("curl", "wget"),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load customized external spec: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		// Add custom headers for external content
		w.Header().Set("X-Content-Source", "external-url")
		w.Header().Set("X-Customized", "true")
		
		w.Write([]byte(content))
	})
}