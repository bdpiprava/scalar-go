package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

// ExampleJavaScriptAPIMode demonstrates using the JavaScript API rendering mode
func ExampleJavaScriptAPIMode() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithRenderMode(scalargo.RenderModeJavaScriptAPI),
	)
}

// ExampleJavaScriptAPIWithCustomJS demonstrates custom JS injection with JavaScript API mode
func ExampleJavaScriptAPIWithCustomJS() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithRenderMode(scalargo.RenderModeJavaScriptAPI),
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

// ExampleHideSearchAndShowOperationId demonstrates display configuration options
func ExampleHideSearchAndShowOperationId() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithRenderMode(scalargo.RenderModeJavaScriptAPI),
		scalargo.WithHideSearch(true),
		scalargo.WithShowOperationId(true),
	)
}

// ExampleDefaultHttpClient demonstrates setting the default HTTP client for code examples
func ExampleDefaultHttpClient() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithRenderMode(scalargo.RenderModeJavaScriptAPI),
		scalargo.WithDefaultHttpClient("node", "undici"),
	)
}

// ExampleSortingOptions demonstrates sorting configuration for tags and operations
func ExampleSortingOptions() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithRenderMode(scalargo.RenderModeJavaScriptAPI),
		scalargo.WithTagsSorter(scalargo.SorterAlpha),
		scalargo.WithOperationsSorter(scalargo.SorterMethod),
		scalargo.WithOperationTitleSource(scalargo.OperationTitleSourcePath),
		scalargo.WithOrderSchemaPropertiesBy(scalargo.SchemaPropertiesOrderAlpha),
	)
}

// ExamplePersistAuth demonstrates enabling authentication persistence in localStorage
func ExamplePersistAuth() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithRenderMode(scalargo.RenderModeJavaScriptAPI),
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
			Title:   "API v1",
			Slug:    "v1",
			URL:     "https://example.com/api/v1/openapi.json",
			Default: true,
		},
		{
			Title: "API v2",
			Slug:  "v2",
			URL:   "https://example.com/api/v2/openapi.json",
		},
		{
			Title: "API v3 Beta",
			Slug:  "v3-beta",
			URL:   "https://example.com/api/v3/openapi.json",
		},
	}

	return scalargo.NewV2(
		scalargo.WithSpecURL("https://example.com/api/v1/openapi.json"),
		scalargo.WithRenderMode(scalargo.RenderModeJavaScriptAPI),
		scalargo.WithMultipleSources(sources...),
	)
}

// ExampleCustomCssInConfig demonstrates setting custom CSS via configuration object
func ExampleCustomCssInConfig() (string, error) {
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
		scalargo.WithRenderMode(scalargo.RenderModeJavaScriptAPI),
		scalargo.WithCustomCss(customCSS),
	)
}

// ExampleAdvancedConfiguration demonstrates combining multiple advanced features
func ExampleAdvancedConfiguration() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithRenderMode(scalargo.RenderModeJavaScriptAPI),
		scalargo.WithTheme(scalargo.ThemePurple),
		scalargo.WithLayout(scalargo.LayoutModern),
		scalargo.WithDarkMode(),
		scalargo.WithHideSearch(false),
		scalargo.WithShowOperationId(true),
		scalargo.WithDefaultHttpClient("node", "undici"),
		scalargo.WithTagsSorter(scalargo.SorterAlpha),
		scalargo.WithOperationsSorter(scalargo.SorterMethod),
		scalargo.WithPersistAuth(true),
		scalargo.WithHideDownloadButton(),
		scalargo.WithCustomHeadJS(`console.log('Scalar initializing...');`),
		scalargo.WithCustomBodyJS(`console.log('Scalar ready');`),
	)
}

// ExampleBackwardCompatibility demonstrates that existing code works without changes
func ExampleBackwardCompatibility() (string, error) {
	// This is the old way - still works perfectly (defaults to data-attribute mode)
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeDefault),
		scalargo.WithLayout(scalargo.LayoutModern),
	)
}
