package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

// API Key Authentication Examples

// ExampleAPIKeySimple demonstrates the simplest API key authentication (header-based, backward compatible)
func ExampleAPIKeySimple() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithAPIKey("demo-api-key-12345"),
		),
	)
}

// ExampleAPIKeyCustomHeader demonstrates API key authentication with a custom header name
func ExampleAPIKeyCustomHeader() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithAPIKey("demo-api-key-12345",
				scalargo.WithAPIKeyName("X-API-Key"),
			),
		),
	)
}

// ExampleAPIKeyQueryParameter demonstrates API key passed as a query parameter
func ExampleAPIKeyQueryParameter() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithAPIKeyQuery("api_key", "demo-key-12345"),
		),
	)
}

// ExampleAPIKeyCookie demonstrates API key passed as a cookie
func ExampleAPIKeyCookie() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithAPIKeyCookie("session_token", "demo-session-token"),
		),
	)
}

// HTTP Basic & Bearer Authentication Examples

// ExampleHTTPBasicAuth demonstrates HTTP Basic authentication
func ExampleHTTPBasicAuth() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithHTTPBasicAuth("demo-user", "demo-password"),
		),
	)
}

// ExampleHTTPBearerToken demonstrates HTTP Bearer token authentication
func ExampleHTTPBearerToken() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithHTTPBearerToken("demo-bearer-token-eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
		),
	)
}

// OAuth2 Authentication Examples

// ExampleOAuth2AuthorizationCode demonstrates OAuth2 Authorization Code flow with PKCE
func ExampleOAuth2AuthorizationCode() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithOAuth2AuthorizationCode(
				"https://auth.example.com/oauth/authorize",
				"https://auth.example.com/oauth/token",
				scalargo.WithOAuth2ClientID("demo-client-id"),
				scalargo.WithOAuth2RedirectURI("https://demo.example.com/callback"),
				scalargo.WithOAuth2PKCE(scalargo.PKCES256), // SHA-256 PKCE (recommended)
				scalargo.WithOAuth2Scopes("read:api", "write:api"),
			),
		),
		scalargo.WithTheme(scalargo.ThemeMoon),
	)
}

// ExampleOAuth2ClientCredentials demonstrates OAuth2 Client Credentials flow for service-to-service auth
func ExampleOAuth2ClientCredentials() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithOAuth2ClientCredentials(
				"https://auth.example.com/oauth/token",
				scalargo.WithOAuth2ClientID("service-account-id"),
				scalargo.WithOAuth2ClientSecret("demo-client-secret"),
				scalargo.WithOAuth2Scopes("read:api", "write:api"),
			),
		),
		scalargo.WithTheme(scalargo.ThemeDeepSpace),
	)
}

// ExampleOAuth2WithAdvancedOptions demonstrates OAuth2 with custom parameters and configuration
func ExampleOAuth2WithAdvancedOptions() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithOAuth2AuthorizationCode(
				"https://auth.example.com/oauth/authorize",
				"https://auth.example.com/oauth/token",
				scalargo.WithOAuth2ClientID("advanced-demo-client"),
				scalargo.WithOAuth2CustomToken("custom_access_token"),
				scalargo.WithOAuth2AdditionalAuthParams(map[string]string{
					"audience": "https://api.example.com",
					"prompt":   "consent",
				}),
				scalargo.WithOAuth2AdditionalTokenParams(map[string]string{
					"resource": "https://resource.example.com",
				}),
				scalargo.WithOAuth2CredentialsLocation(scalargo.OAuth2CredentialsHeader),
			),
		),
	)
}

// ExampleOAuth2PKCE demonstrates different PKCE modes
func ExampleOAuth2PKCE() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithOAuth2AuthorizationCode(
				"https://auth.example.com/oauth/authorize",
				"https://auth.example.com/oauth/token",
				scalargo.WithOAuth2ClientID("pkce-demo-client"),
				scalargo.WithOAuth2PKCE(scalargo.PKCES256), // Recommended: SHA-256 PKCE
				// Other options:
				// scalargo.WithOAuth2PKCE(scalargo.PKCEPlain) // Plain PKCE
				// scalargo.WithOAuth2PKCE(scalargo.PKCEDisabled) // Disabled
			),
		),
	)
}

// Multiple Security Schemes Examples

// ExampleMultipleSecuritySchemes demonstrates configuring multiple authentication methods
func ExampleMultipleSecuritySchemes() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			// Define multiple security schemes
			scalargo.WithSecurityScheme("api_key",
				scalargo.APIKeyScheme("X-API-Key", scalargo.APIKeyLocationHeader, "demo-api-key"),
			),
			scalargo.WithSecurityScheme("bearer_auth",
				scalargo.BearerScheme("demo-bearer-token"),
			),
			scalargo.WithSecurityScheme("basic_auth",
				scalargo.BasicScheme("demo-user", "demo-password"),
			),
			// Set the preferred authentication method
			scalargo.WithPreferredSecurityScheme("bearer_auth"),
		),
		scalargo.WithTheme(scalargo.ThemePurple),
	)
}

// ExampleMultipleSecuritySchemesWithOAuth2 demonstrates multiple auth methods including OAuth2
func ExampleMultipleSecuritySchemesWithOAuth2() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			// API Key scheme
			scalargo.WithSecurityScheme("api_key",
				scalargo.APIKeyScheme("X-API-Key", scalargo.APIKeyLocationHeader, "demo-key-12345"),
			),
			// Bearer token scheme
			scalargo.WithSecurityScheme("bearer_auth",
				scalargo.BearerScheme("demo-bearer-token-xyz"),
			),
			// OAuth2 scheme
			scalargo.WithSecurityScheme("oauth2",
				scalargo.OAuth2Scheme(
					scalargo.OAuth2FlowAuthorizationCode,
					scalargo.OAuth2Config{
						AuthorizationURL: "https://auth.example.com/authorize",
						TokenURL:         "https://auth.example.com/token",
						ClientID:         "demo-oauth-client",
						UsePKCE:          scalargo.PKCES256,
						SelectedScopes:   []string{"read:api", "write:api", "admin:api"},
					},
				),
			),
			// Prefer OAuth2 for this API
			scalargo.WithPreferredSecurityScheme("oauth2"),
		),
		scalargo.WithTheme(scalargo.ThemeSaturn),
		scalargo.WithLayout(scalargo.LayoutModern),
	)
}

// Real-World Scenarios

// ExampleGitHubStyleAuth demonstrates GitHub-style authentication with multiple options
func ExampleGitHubStyleAuth() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			// Personal Access Token (most common)
			scalargo.WithSecurityScheme("pat",
				scalargo.BearerScheme("ghp_demo_personal_access_token_xxxxxxxxxxxx"),
			),
			// OAuth2 App
			scalargo.WithSecurityScheme("oauth_app",
				scalargo.OAuth2Scheme(
					scalargo.OAuth2FlowAuthorizationCode,
					scalargo.OAuth2Config{
						AuthorizationURL: "https://github.com/login/oauth/authorize",
						TokenURL:         "https://github.com/login/oauth/access_token",
						ClientID:         "demo_github_client_id",
						SelectedScopes:   []string{"repo", "user", "workflow"},
					},
				),
			),
			// GitHub App
			scalargo.WithSecurityScheme("github_app",
				scalargo.BearerScheme("demo_github_app_installation_token"),
			),
			scalargo.WithPreferredSecurityScheme("pat"),
		),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("GitHub-Style API Documentation"),
			scalargo.WithKeyValue("description", "API with multiple authentication methods like GitHub"),
		),
		scalargo.WithTheme(scalargo.ThemeMoon),
		scalargo.WithDarkMode(),
	)
}

// ExampleStripeStyleAuth demonstrates Stripe-style authentication with test/live keys
func ExampleStripeStyleAuth() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			// Test mode secret key
			scalargo.WithSecurityScheme("test_key",
				scalargo.BearerScheme("sk_test_demo_51xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
			),
			// Live mode secret key
			scalargo.WithSecurityScheme("live_key",
				scalargo.BearerScheme("sk_live_demo_51xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
			),
			scalargo.WithPreferredSecurityScheme("test_key"),
		),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Stripe-Style API Documentation"),
			scalargo.WithKeyValue("description", "API with test/live key separation"),
		),
		scalargo.WithTheme(scalargo.ThemePurple),
	)
}

// ExampleAuth0StyleOAuth2 demonstrates Auth0-style OAuth2 configuration
func ExampleAuth0StyleOAuth2() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			scalargo.WithOAuth2AuthorizationCode(
				"https://your-tenant.auth0.com/authorize",
				"https://your-tenant.auth0.com/oauth/token",
				scalargo.WithOAuth2ClientID("demo_auth0_client_id"),
				scalargo.WithOAuth2RedirectURI("https://your-app.com/callback"),
				scalargo.WithOAuth2PKCE(scalargo.PKCES256),
				scalargo.WithOAuth2AdditionalAuthParams(map[string]string{
					"audience": "https://your-api.example.com",
				}),
				scalargo.WithOAuth2Scopes("openid", "profile", "email", "read:data", "write:data"),
			),
		),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Auth0-Protected API Documentation"),
			scalargo.WithKeyValue("description", "API secured with Auth0 OAuth2"),
		),
		scalargo.WithTheme(scalargo.ThemeDeepSpace),
	)
}

// ExampleProductionAPIWithFullAuth demonstrates a production-ready API with comprehensive auth setup
func ExampleProductionAPIWithFullAuth() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithAuthenticationOpts(
			// API Key for simple integrations
			scalargo.WithSecurityScheme("api_key",
				scalargo.APIKeyScheme("X-API-Key", scalargo.APIKeyLocationHeader, ""),
			),
			// Bearer token for token-based auth
			scalargo.WithSecurityScheme("bearer_token",
				scalargo.BearerScheme(""),
			),
			// OAuth2 for third-party integrations
			scalargo.WithSecurityScheme("oauth2",
				scalargo.OAuth2Scheme(
					scalargo.OAuth2FlowAuthorizationCode,
					scalargo.OAuth2Config{
						AuthorizationURL:    "https://auth.production-api.com/authorize",
						TokenURL:            "https://auth.production-api.com/token",
						ClientID:            "", // Client provides their own
						RedirectURI:         "", // Client provides their own
						UsePKCE:             scalargo.PKCES256,
						SelectedScopes:      []string{"read", "write"},
						CredentialsLocation: scalargo.OAuth2CredentialsHeader,
					},
				),
			),
			// Client Credentials for server-to-server
			scalargo.WithSecurityScheme("client_credentials",
				scalargo.OAuth2Scheme(
					scalargo.OAuth2FlowClientCredentials,
					scalargo.OAuth2Config{
						TokenURL:       "https://auth.production-api.com/token",
						ClientID:       "", // Client provides their own
						SelectedScopes: []string{"api:read", "api:write"},
					},
				),
			),
			scalargo.WithPreferredSecurityScheme("oauth2"),
		),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Production API Documentation"),
			scalargo.WithKeyValue("description", "Comprehensive API with multiple authentication methods"),
		),
		scalargo.WithTheme(scalargo.ThemeDefault),
		scalargo.WithLayout(scalargo.LayoutModern),
		scalargo.WithDarkMode(),
		scalargo.WithHideDownloadButton(),
	)
}
