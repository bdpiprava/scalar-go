package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

// ExampleJavaScriptAPIMode demonstrates using the JavaScript API rendering mode (now the default)
func ExampleJavaScriptAPIMode() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		// No need to specify WithRenderMode - JavaScript API is now the default
	)
}

// ExampleJavaScriptAPIWithCustomJS demonstrates custom JS injection with JavaScript API mode
func ExampleJavaScriptAPIWithCustomJS() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithCustomHeadJS(`
			// Custom initialization code in head
			console.log('Scalar is about to initialize');
		`),
		scalargo.WithCustomBodyJS(`
			// Add event listener after Scalar initializes
			document.addEventListener('DOMContentLoaded', () => {
				console.log('Page loaded with Scalar');
			});
		`),
	)
}

// ExampleDefaultHTTPClient demonstrates setting the default HTTP client for code examples
func ExampleDefaultHTTPClient() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithDefaultHTTPClient("node", "undici"),
	)
}

// ExamplePersistAuth demonstrates enabling authentication persistence in localStorage
func ExamplePersistAuth() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithPersistAuth(true),
		scalargo.WithAuthenticationOpts(
			scalargo.WithAPIKey("demo-api-key"),
		),
	)
}

// ExampleMultipleSources demonstrates configuring multiple OpenAPI document sources
func ExampleMultipleSources() (string, error) {
	sources := []scalargo.DocumentSource{
		{
			Title:   "Scalar Galaxy API",
			Slug:    "scalar-galaxy",
			URL:     "https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml",
			Default: true,
		},
		{
			Title: "Petstore API",
			Slug:  "petstore",
			URL:   "https://petstore3.swagger.io/api/v3/openapi.json",
		},
		{
			Title: "GitHub REST API",
			Slug:  "github",
			URL:   "https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.json",
		},
	}

	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		scalargo.WithMultipleSources(sources...),
	)
}

// ExampleCustomCSSInConfig demonstrates setting custom CSS via configuration object
func ExampleCustomCSSInConfig() (string, error) {
	customCSS := `
		.scalar-card {
			border-radius: 12px;
			box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		}
		.scalar-button {
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		}
	`

	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithCustomCSS(customCSS),
	)
}

// ExampleAdvancedConfiguration demonstrates combining multiple advanced features
func ExampleAdvancedConfiguration() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithLayout(scalargo.LayoutModern),
		scalargo.WithDarkMode(),
		scalargo.WithHideSearch(false),
		scalargo.WithShowOperationID(true),
		scalargo.WithDefaultHTTPClient("node", "undici"),
		scalargo.WithTagsSorter(scalargo.SorterAlpha),
		scalargo.WithOperationsSorter(scalargo.SorterMethod),
		scalargo.WithPersistAuth(true),
		scalargo.WithHideDownloadButton(),
		scalargo.WithCustomHeadJS(`console.log('Scalar initializing...');`),
		scalargo.WithCustomBodyJS(`console.log('Scalar ready');`),
	)
}

// ExampleDataAttributeMode demonstrates using the legacy data-attribute rendering mode for backward compatibility
func ExampleDataAttributeMode() (string, error) {
	// Explicitly use the legacy data-attribute mode (pre-v1.0 behavior)
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithRenderMode(scalargo.RenderModeDataAttribute),
		scalargo.WithLayout(scalargo.LayoutModern),
	)
}

// ExampleToolbarVisibility demonstrates controlling the developer toolbar visibility
// Options: ShowToolbarAlways (all environments), ShowToolbarLocalhost (local dev only), ShowToolbarNever (hidden)
func ExampleToolbarVisibility() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		// ShowToolbarAlways: Visible in all environments for debugging
		// ShowToolbarLocalhost: Visible only on localhost (Scalar's default)
		// ShowToolbarNever: Completely hidden (this library's default)
		scalargo.WithShowToolbar(scalargo.ShowToolbarNever),
	)
}

// ExampleHideSearch demonstrates hiding the search functionality
func ExampleHideSearch() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithHideSearch(true),
	)
}

// ExampleShowOperationID demonstrates displaying operation IDs in the UI
func ExampleShowOperationID() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithShowOperationID(true),
	)
}

// ExampleTagsSorter demonstrates sorting tags alphabetically in the sidebar
func ExampleTagsSorter() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTagsSorter(scalargo.SorterAlpha),
	)
}

// ExampleOperationsSorter demonstrates sorting operations within tags
func ExampleOperationsSorter() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		// SorterAlpha: Alphabetically by title
		// SorterMethod: By HTTP method (GET, POST, PUT, DELETE, etc.)
		scalargo.WithOperationsSorter(scalargo.SorterMethod),
	)
}

// ExampleOperationTitleSource demonstrates choosing operation title source
func ExampleOperationTitleSource() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		// OperationTitleSourceSummary: Use operation summary
		// OperationTitleSourcePath: Use operation path
		scalargo.WithOperationTitleSource(scalargo.OperationTitleSourcePath),
	)
}

// ExampleSchemaPropertiesOrder demonstrates controlling schema property order
func ExampleSchemaPropertiesOrder() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		// SchemaPropertiesOrderAlpha: Sort alphabetically
		// SchemaPropertiesOrderPreserve: Preserve original spec order
		scalargo.WithOrderSchemaPropertiesBy(scalargo.SchemaPropertiesOrderAlpha),
	)
}
