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

// ExampleAPIV1WithAuth demonstrates API v1 with basic authentication
func ExampleAPIV1WithAuth() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Pet Store API v1 (Authenticated)"),
			scalargo.WithKeyValue("version", "1.0.0"),
		),
		scalargo.WithAuthenticationOpts(
			scalargo.WithAPIKey("v1-api-key-demo"),
		),
		scalargo.WithTheme(scalargo.ThemeDefault),
	)
}

// ExampleAPIV2WithOAuth2 demonstrates API v2 with OAuth2 authentication
func ExampleAPIV2WithOAuth2() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Pet Store API v2 (OAuth2)"),
			scalargo.WithKeyValue("version", "2.0.0"),
		),
		scalargo.WithAuthenticationOpts(
			scalargo.WithOAuth2AuthorizationCode(
				"https://auth.petstore.example.com/authorize",
				"https://auth.petstore.example.com/token",
				scalargo.WithOAuth2ClientID("petstore-v2-client"),
				scalargo.WithOAuth2PKCE(scalargo.PKCES256),
				scalargo.WithOAuth2Scopes("pets:read", "pets:write"),
			),
		),
		scalargo.WithTheme(scalargo.ThemeSolarized),
		scalargo.WithDarkMode(),
	)
}

// ExampleProductionAPI demonstrates a production API with multiple environments and full auth
func ExampleProductionAPI() (string, error) {
	env := getEnvironment()
	var authURL, tokenURL string

	// Environment-specific OAuth2 endpoints
	switch env {
	case "production":
		authURL = "https://auth.production.example.com/authorize"
		tokenURL = "https://auth.production.example.com/token"
	case "staging":
		authURL = "https://auth.staging.example.com/authorize"
		tokenURL = "https://auth.staging.example.com/token"
	default:
		authURL = "https://auth.dev.example.com/authorize"
		tokenURL = "https://auth.dev.example.com/token"
	}

	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Production API Documentation"),
			scalargo.WithKeyValue("environment", env),
			scalargo.WithKeyValue("version", "3.0.0"),
			scalargo.WithKeyValue("generated", time.Now().Format(time.RFC3339)),
		),
		scalargo.WithAuthenticationOpts(
			// API Key for simple requests
			scalargo.WithSecurityScheme("api_key",
				scalargo.APIKeyScheme("X-API-Key", scalargo.APIKeyLocationHeader, ""),
			),
			// OAuth2 for user-authenticated requests
			scalargo.WithSecurityScheme("oauth2",
				scalargo.OAuth2Scheme(
					scalargo.OAuth2FlowAuthorizationCode,
					scalargo.OAuth2Config{
						AuthorizationURL: authURL,
						TokenURL:         tokenURL,
						UsePKCE:          scalargo.PKCES256,
						SelectedScopes:   []string{"read", "write", "admin"},
					},
				),
			),
			// Client Credentials for service-to-service
			scalargo.WithSecurityScheme("service_auth",
				scalargo.OAuth2Scheme(
					scalargo.OAuth2FlowClientCredentials,
					scalargo.OAuth2Config{
						TokenURL:       tokenURL,
						SelectedScopes: []string{"service:read", "service:write"},
					},
				),
			),
			scalargo.WithPreferredSecurityScheme("oauth2"),
		),
		scalargo.WithTheme(scalargo.ThemeDeepSpace),
		scalargo.WithDarkMode(),
		scalargo.WithLayout(scalargo.LayoutModern),
	)
}

// getEnvironment returns the current environment (mock implementation)
func getEnvironment() string {
	// In a real application, this might read from environment variables
	// or configuration files
	return "development"
}
