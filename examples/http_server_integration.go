package examples

import (
	"time"

	scalargo "github.com/bdpiprava/scalar-go"
)

// ExampleStaticDocumentation demonstrates simple static documentation generation
func ExampleStaticDocumentation() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
	)
}

// ExampleDynamicDocumentation demonstrates documentation with dynamic metadata
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

// getEnvironment returns the current environment (mock implementation)
func getEnvironment() string {
	// In a real application, this might read from environment variables
	// or configuration files
	return "development"
}
