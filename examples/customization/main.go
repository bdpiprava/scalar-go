package main

import (
	"fmt"
	"log"
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
)

// This example demonstrates various customization options available in scalar-go.
// You can customize themes, layouts, visibility options, and add custom CSS.
func main() {
	// Directory containing your OpenAPI spec
	apiDir := "./data/loader"

	// Setup multiple endpoints with different customization examples
	setupThemeExamples(apiDir)
	setupLayoutExamples(apiDir)
	setupVisibilityExamples(apiDir)
	setupAdvancedCustomization(apiDir)

	fmt.Println("Customization examples available at:")
	fmt.Println("  Theme examples:")
	fmt.Println("    http://localhost:8082/theme/default")
	fmt.Println("    http://localhost:8082/theme/moon") 
	fmt.Println("    http://localhost:8082/theme/purple")
	fmt.Println("    http://localhost:8082/theme/solarized")
	fmt.Println("  Layout examples:")
	fmt.Println("    http://localhost:8082/layout/modern")
	fmt.Println("    http://localhost:8082/layout/classic")
	fmt.Println("  Visibility examples:")
	fmt.Println("    http://localhost:8082/visibility/hide-sidebar")
	fmt.Println("    http://localhost:8082/visibility/hide-models")
	fmt.Println("    http://localhost:8082/visibility/dark-mode")
	fmt.Println("  Advanced:")
	fmt.Println("    http://localhost:8082/advanced/custom-css")
	fmt.Println("    http://localhost:8082/advanced/all-options")
	fmt.Println("Press Ctrl+C to stop the server")

	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// setupThemeExamples demonstrates different theme options
func setupThemeExamples(apiDir string) {
	// Default theme
	http.HandleFunc("/theme/default", createHandler(apiDir, 
		scalargo.WithTheme(scalargo.ThemeDefault)))
	
	// Moon theme - dark theme with blue accents
	http.HandleFunc("/theme/moon", createHandler(apiDir,
		scalargo.WithTheme(scalargo.ThemeMoon)))
	
	// Purple theme - purple color scheme
	http.HandleFunc("/theme/purple", createHandler(apiDir,
		scalargo.WithTheme(scalargo.ThemePurple)))
	
	// Solarized theme - based on the popular Solarized color scheme
	http.HandleFunc("/theme/solarized", createHandler(apiDir,
		scalargo.WithTheme(scalargo.ThemeSolarized)))
}

// setupLayoutExamples demonstrates different layout options
func setupLayoutExamples(apiDir string) {
	// Modern layout (default) - contemporary design
	http.HandleFunc("/layout/modern", createHandler(apiDir,
		scalargo.WithLayout(scalargo.LayoutModern)))
	
	// Classic layout - traditional documentation layout
	http.HandleFunc("/layout/classic", createHandler(apiDir,
		scalargo.WithLayout(scalargo.LayoutClassic)))
}

// setupVisibilityExamples demonstrates visibility and UI control options
func setupVisibilityExamples(apiDir string) {
	// Hide the sidebar for a cleaner look
	http.HandleFunc("/visibility/hide-sidebar", createHandler(apiDir,
		scalargo.WithSidebarVisibility(false)))
	
	// Hide the models section to focus on endpoints
	http.HandleFunc("/visibility/hide-models", createHandler(apiDir,
		scalargo.WithHideModels()))
	
	// Enable dark mode by default
	http.HandleFunc("/visibility/dark-mode", createHandler(apiDir,
		scalargo.WithDarkMode()))
}

// setupAdvancedCustomization demonstrates advanced customization options
func setupAdvancedCustomization(apiDir string) {
	// Custom CSS override example
	customCSS := `
		/* Custom styling for the API documentation */
		.section-header {
			color: #e74c3c !important;
			font-weight: bold !important;
		}
		
		.api-client__request {
			background-color: #f8f9fa !important;
			border-left: 4px solid #007bff !important;
		}
		
		/* Custom button styling */
		button {
			border-radius: 8px !important;
		}
	`
	
	http.HandleFunc("/advanced/custom-css", createHandler(apiDir,
		scalargo.WithOverrideCSS(customCSS),
		scalargo.WithTheme(scalargo.ThemeDefault)))
	
	// Combination of multiple options
	http.HandleFunc("/advanced/all-options", createHandler(apiDir,
		scalargo.WithTheme(scalargo.ThemeMoon),
		scalargo.WithLayout(scalargo.LayoutModern),
		scalargo.WithDarkMode(),
		scalargo.WithHideDownloadButton(),
		scalargo.WithHiddenClients("fetch", "curl"), // Hide specific client examples
		scalargo.WithSearchHotKey("cmd+k"),
		scalargo.WithOverrideCSS(`
			.scalar-api-reference {
				--scalar-color-1: #2d3748;
				--scalar-color-2: #4a5568;
			}
		`)))
}

// createHandler is a helper function that creates an HTTP handler 
// with the specified scalar-go options
func createHandler(apiDir string, opts ...scalargo.Option) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate the HTML content with the provided options
		htmlContent, err := scalargo.New(apiDir, opts...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate documentation: %v", err), http.StatusInternalServerError)
			return
		}
		
		// Serve the HTML content
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	}
}