package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

const githubAPI = "https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.json"

// ExampleScalarGalaxy demonstrates loading the Scalar Galaxy API specification from CDN
func ExampleScalarGalaxy() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
	)
}

// ExamplePetstore demonstrates loading the classic Petstore API specification
func ExamplePetstore() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL("https://petstore3.swagger.io/api/v3/openapi.json"),
	)
}

// ExampleGitHubAPI demonstrates loading GitHub API documentation from public OpenAPI spec
func ExampleGitHubAPI() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL(githubAPI),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("GitHub REST API"),
			scalargo.WithKeyValue("description", "Complete GitHub REST API documentation"),
		),
		scalargo.WithTheme(scalargo.ThemeDefault),
	)
}

// ExampleOpenAIAPI demonstrates loading external API documentation with custom metadata
func ExampleOpenAIAPI() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("OpenAI API (Demo)"),
			scalargo.WithKeyValue("description", "Example of loading external API documentation"),
			scalargo.WithKeyValue("note", "This is a demo using Scalar Galaxy spec"),
		),
		scalargo.WithTheme(scalargo.ThemeAlternate),
		scalargo.WithLayout(scalargo.LayoutClassic),
	)
}

// ExampleCustomizedExternal demonstrates URL loading with extensive branding customization
func ExampleCustomizedExternal() (string, error) {
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

	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Customized External API Docs"),
			scalargo.WithKeyValue("company", "Your Company Name"),
			scalargo.WithKeyValue("customized", "true"),
		),
		scalargo.WithTheme(scalargo.ThemeMoon),
		scalargo.WithLayout(scalargo.LayoutModern),
		scalargo.WithOverrideCSS(customCSS),
		scalargo.WithDarkMode(),
		scalargo.WithHideDownloadButton(),
		scalargo.WithSearchHotKey("cmd+k"),
		scalargo.WithHiddenClients("curl", "wget"),
	)
}
