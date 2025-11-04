package scalargo_test

import (
	"encoding/json"
	"testing"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Enhanced API Key Configuration Tests

func TestWithAPIKey_BackwardCompatibility(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithAPIKey("my-token")(auth)

	apiKey, ok := auth["apiKey"].(map[string]any)
	require.True(t, ok, "apiKey should exist and be a map")
	assert.Equal(t, "my-token", apiKey["token"])
	assert.Equal(t, "Authorization", apiKey["name"])
	assert.Equal(t, "header", apiKey["in"])
}

func TestWithAPIKey_CustomHeaderName(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithAPIKey("my-token", scalargo.WithAPIKeyName("X-API-Key"))(auth)

	apiKey, ok := auth["apiKey"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "my-token", apiKey["token"])
	assert.Equal(t, "X-API-Key", apiKey["name"])
	assert.Equal(t, "header", apiKey["in"])
}

func TestWithAPIKey_QueryLocation(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithAPIKey("my-token", scalargo.WithAPIKeyLocation(scalargo.APIKeyLocationQuery))(auth)

	apiKey, ok := auth["apiKey"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "my-token", apiKey["token"])
	assert.Equal(t, "api_key", apiKey["name"]) // default for query
	assert.Equal(t, "query", apiKey["in"])
}

func TestWithAPIKey_CookieLocation(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithAPIKey("my-token", scalargo.WithAPIKeyLocation(scalargo.APIKeyLocationCookie))(auth)

	apiKey, ok := auth["apiKey"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "my-token", apiKey["token"])
	assert.Equal(t, "api_key", apiKey["name"]) // default for cookie
	assert.Equal(t, "cookie", apiKey["in"])
}

func TestWithAPIKey_CustomNameAndLocation(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithAPIKey("my-token",
		scalargo.WithAPIKeyName("apikey"),
		scalargo.WithAPIKeyLocation(scalargo.APIKeyLocationQuery),
	)(auth)

	apiKey, ok := auth["apiKey"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "my-token", apiKey["token"])
	assert.Equal(t, "apikey", apiKey["name"])
	assert.Equal(t, "query", apiKey["in"])
}

func TestWithAPIKeyHeader(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithAPIKeyHeader("X-Custom-Key", "token-123")(auth)

	apiKey, ok := auth["apiKey"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "token-123", apiKey["token"])
	assert.Equal(t, "X-Custom-Key", apiKey["name"])
	assert.Equal(t, "header", apiKey["in"])
}

func TestWithAPIKeyQuery(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithAPIKeyQuery("api_key", "token-456")(auth)

	apiKey, ok := auth["apiKey"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "token-456", apiKey["token"])
	assert.Equal(t, "api_key", apiKey["name"])
	assert.Equal(t, "query", apiKey["in"])
}

func TestWithAPIKeyCookie(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithAPIKeyCookie("session", "token-789")(auth)

	apiKey, ok := auth["apiKey"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "token-789", apiKey["token"])
	assert.Equal(t, "session", apiKey["name"])
	assert.Equal(t, "cookie", apiKey["in"])
}

// OAuth2 Configuration Tests

func TestWithOAuth2AuthorizationCode(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithOAuth2AuthorizationCode(
		"https://auth.example.com/authorize",
		"https://auth.example.com/token",
		scalargo.WithOAuth2ClientID("my-client-id"),
		scalargo.WithOAuth2RedirectURI("https://app.example.com/callback"),
		scalargo.WithOAuth2PKCE(scalargo.PKCES256),
		scalargo.WithOAuth2Scopes("read:api", "write:api"),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	oauth2Scheme, ok := schemes["oauth2"].(map[string]any)
	require.True(t, ok)

	flows, ok := oauth2Scheme["flows"].(map[string]any)
	require.True(t, ok)

	authCode, ok := flows["authorizationCode"].(*scalargo.OAuth2Config)
	require.True(t, ok)

	assert.Equal(t, "https://auth.example.com/authorize", authCode.AuthorizationURL)
	assert.Equal(t, "https://auth.example.com/token", authCode.TokenURL)
	assert.Equal(t, "my-client-id", authCode.ClientID)
	assert.Equal(t, "https://app.example.com/callback", authCode.RedirectURI)
	assert.Equal(t, scalargo.PKCES256, authCode.UsePKCE)
	assert.Equal(t, []string{"read:api", "write:api"}, authCode.SelectedScopes)
}

func TestWithOAuth2ClientCredentials(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithOAuth2ClientCredentials(
		"https://auth.example.com/token",
		scalargo.WithOAuth2ClientID("service-account"),
		scalargo.WithOAuth2ClientSecret("super-secret"),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	oauth2Scheme, ok := schemes["oauth2"].(map[string]any)
	require.True(t, ok)

	flows, ok := oauth2Scheme["flows"].(map[string]any)
	require.True(t, ok)

	clientCreds, ok := flows["clientCredentials"].(*scalargo.OAuth2Config)
	require.True(t, ok)

	assert.Equal(t, "https://auth.example.com/token", clientCreds.TokenURL)
	assert.Equal(t, "service-account", clientCreds.ClientID)
	assert.Equal(t, "super-secret", clientCreds.ClientSecret)
}

func TestWithOAuth2Implicit(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithOAuth2Implicit(
		"https://auth.example.com/authorize",
		scalargo.WithOAuth2ClientID("implicit-client"),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	oauth2Scheme, ok := schemes["oauth2"].(map[string]any)
	require.True(t, ok)

	flows, ok := oauth2Scheme["flows"].(map[string]any)
	require.True(t, ok)

	implicit, ok := flows["implicit"].(*scalargo.OAuth2Config)
	require.True(t, ok)

	assert.Equal(t, "https://auth.example.com/authorize", implicit.AuthorizationURL)
	assert.Equal(t, "implicit-client", implicit.ClientID)
}

func TestWithOAuth2Password(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithOAuth2Password(
		"https://auth.example.com/token",
		scalargo.WithOAuth2ClientID("password-client"),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	oauth2Scheme, ok := schemes["oauth2"].(map[string]any)
	require.True(t, ok)

	flows, ok := oauth2Scheme["flows"].(map[string]any)
	require.True(t, ok)

	password, ok := flows["password"].(*scalargo.OAuth2Config)
	require.True(t, ok)

	assert.Equal(t, "https://auth.example.com/token", password.TokenURL)
	assert.Equal(t, "password-client", password.ClientID)
}

func TestOAuth2PKCE_Modes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		pkceMode scalargo.OAuth2PKCEMode
		expected scalargo.OAuth2PKCEMode
	}{
		{
			name:     "PKCE S256 mode",
			pkceMode: scalargo.PKCES256,
			expected: "S256",
		},
		{
			name:     "PKCE plain mode",
			pkceMode: scalargo.PKCEPlain,
			expected: "plain",
		},
		{
			name:     "PKCE disabled",
			pkceMode: scalargo.PKCEDisabled,
			expected: "no",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			auth := make(scalargo.AuthenticationOption)
			scalargo.WithOAuth2AuthorizationCode(
				"https://auth.example.com/authorize",
				"https://auth.example.com/token",
				scalargo.WithOAuth2PKCE(tc.pkceMode),
			)(auth)

			schemes, ok := auth["securitySchemes"].(map[string]any)
			require.True(t, ok)

			oauth2Scheme, ok := schemes["oauth2"].(map[string]any)
			require.True(t, ok)

			flows, ok := oauth2Scheme["flows"].(map[string]any)
			require.True(t, ok)

			authCode, ok := flows["authorizationCode"].(*scalargo.OAuth2Config)
			require.True(t, ok)

			assert.Equal(t, tc.expected, authCode.UsePKCE)
		})
	}
}

func TestOAuth2_AdditionalParams(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithOAuth2AuthorizationCode(
		"https://auth.example.com/authorize",
		"https://auth.example.com/token",
		scalargo.WithOAuth2AdditionalAuthParams(map[string]string{
			"audience": "https://api.example.com",
		}),
		scalargo.WithOAuth2AdditionalTokenParams(map[string]string{
			"resource": "https://resource.example.com",
		}),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	oauth2Scheme, ok := schemes["oauth2"].(map[string]any)
	require.True(t, ok)

	flows, ok := oauth2Scheme["flows"].(map[string]any)
	require.True(t, ok)

	authCode, ok := flows["authorizationCode"].(*scalargo.OAuth2Config)
	require.True(t, ok)

	assert.Equal(t, map[string]string{"audience": "https://api.example.com"}, authCode.SecurityQuery)
	assert.Equal(t, map[string]string{"resource": "https://resource.example.com"}, authCode.SecurityBody)
}

func TestOAuth2_CustomToken(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithOAuth2AuthorizationCode(
		"https://auth.example.com/authorize",
		"https://auth.example.com/token",
		scalargo.WithOAuth2CustomToken("custom_access_token"),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	oauth2Scheme, ok := schemes["oauth2"].(map[string]any)
	require.True(t, ok)

	flows, ok := oauth2Scheme["flows"].(map[string]any)
	require.True(t, ok)

	authCode, ok := flows["authorizationCode"].(*scalargo.OAuth2Config)
	require.True(t, ok)

	assert.Equal(t, "custom_access_token", authCode.TokenName)
}

func TestOAuth2_CredentialsLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location scalargo.OAuth2CredentialsLocation
		expected string
	}{
		{
			name:     "credentials in header",
			location: scalargo.OAuth2CredentialsHeader,
			expected: "header",
		},
		{
			name:     "credentials in body",
			location: scalargo.OAuth2CredentialsBody,
			expected: "body",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			auth := make(scalargo.AuthenticationOption)
			scalargo.WithOAuth2ClientCredentials(
				"https://auth.example.com/token",
				scalargo.WithOAuth2CredentialsLocation(tc.location),
			)(auth)

			schemes, ok := auth["securitySchemes"].(map[string]any)
			require.True(t, ok)

			oauth2Scheme, ok := schemes["oauth2"].(map[string]any)
			require.True(t, ok)

			flows, ok := oauth2Scheme["flows"].(map[string]any)
			require.True(t, ok)

			clientCreds, ok := flows["clientCredentials"].(*scalargo.OAuth2Config)
			require.True(t, ok)

			assert.Equal(t, tc.expected, string(clientCreds.CredentialsLocation))
		})
	}
}

// Security Schemes Configuration Tests

func TestWithSecurityScheme_APIKey(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithSecurityScheme("api_key",
		scalargo.APIKeyScheme("X-API-Key", scalargo.APIKeyLocationHeader, "test-key"),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	apiKeyScheme, ok := schemes["api_key"].(scalargo.SecuritySchemeConfig)
	require.True(t, ok)

	assert.Equal(t, "apiKey", apiKeyScheme["type"])
	assert.Equal(t, "X-API-Key", apiKeyScheme["name"])
	assert.Equal(t, "header", apiKeyScheme["in"])
	assert.Equal(t, "test-key", apiKeyScheme["value"])
}

func TestWithSecurityScheme_Bearer(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithSecurityScheme("bearer_auth",
		scalargo.BearerScheme("my-bearer-token"),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	bearerScheme, ok := schemes["bearer_auth"].(scalargo.SecuritySchemeConfig)
	require.True(t, ok)

	assert.Equal(t, "http", bearerScheme["type"])
	assert.Equal(t, "bearer", bearerScheme["scheme"])
	assert.Equal(t, "my-bearer-token", bearerScheme["token"])
}

func TestWithSecurityScheme_Basic(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithSecurityScheme("basic_auth",
		scalargo.BasicScheme("admin", "password"),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	basicScheme, ok := schemes["basic_auth"].(scalargo.SecuritySchemeConfig)
	require.True(t, ok)

	assert.Equal(t, "http", basicScheme["type"])
	assert.Equal(t, "basic", basicScheme["scheme"])
	assert.Equal(t, "admin", basicScheme["username"])
	assert.Equal(t, "password", basicScheme["password"])
}

func TestWithSecurityScheme_OAuth2(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithSecurityScheme("oauth2",
		scalargo.OAuth2Scheme(
			scalargo.OAuth2FlowAuthorizationCode,
			scalargo.OAuth2Config{
				AuthorizationURL: "https://auth.example.com/authorize",
				TokenURL:         "https://auth.example.com/token",
				ClientID:         "my-client",
			},
		),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	oauth2Scheme, ok := schemes["oauth2"].(scalargo.SecuritySchemeConfig)
	require.True(t, ok)

	assert.Equal(t, "oauth2", oauth2Scheme["type"])

	flows, ok := oauth2Scheme["flows"].(map[string]any)
	require.True(t, ok)

	authCode, ok := flows["authorizationCode"].(scalargo.OAuth2Config)
	require.True(t, ok)

	assert.Equal(t, "https://auth.example.com/authorize", authCode.AuthorizationURL)
	assert.Equal(t, "https://auth.example.com/token", authCode.TokenURL)
	assert.Equal(t, "my-client", authCode.ClientID)
}

func TestMultipleSecuritySchemes(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithSecurityScheme("api_key",
		scalargo.APIKeyScheme("X-API-Key", scalargo.APIKeyLocationHeader, "key-123"),
	)(auth)
	scalargo.WithSecurityScheme("bearer_auth",
		scalargo.BearerScheme("token-456"),
	)(auth)
	scalargo.WithSecurityScheme("basic_auth",
		scalargo.BasicScheme("user", "pass"),
	)(auth)

	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)

	assert.Contains(t, schemes, "api_key")
	assert.Contains(t, schemes, "bearer_auth")
	assert.Contains(t, schemes, "basic_auth")
}

func TestWithPreferredSecurityScheme_Single(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithPreferredSecurityScheme("api_key")(auth)

	preferred := auth["preferredSecurityScheme"]
	assert.Equal(t, []any{"api_key"}, preferred)
}

func TestWithPreferredSecurityScheme_Multiple(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	scalargo.WithPreferredSecurityScheme("api_key", "bearer_auth")(auth)

	preferred := auth["preferredSecurityScheme"]
	assert.Equal(t, []any{"api_key", "bearer_auth"}, preferred)
}

// JSON Serialization Tests

func TestOAuth2Config_JSONSerialization(t *testing.T) {
	t.Parallel()

	config := scalargo.OAuth2Config{
		AuthorizationURL:    "https://auth.example.com/authorize",
		TokenURL:            "https://auth.example.com/token",
		ClientID:            "client-123",
		ClientSecret:        "secret-456",
		RedirectURI:         "https://app.example.com/callback",
		UsePKCE:             scalargo.PKCES256,
		SelectedScopes:      []string{"read", "write"},
		SecurityQuery:       map[string]string{"audience": "api"},
		SecurityBody:        map[string]string{"resource": "res"},
		TokenName:           "access_token",
		CredentialsLocation: scalargo.OAuth2CredentialsHeader,
	}

	jsonBytes, err := json.Marshal(config)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(jsonBytes, &result)
	require.NoError(t, err)

	assert.Equal(t, "https://auth.example.com/authorize", result["authorizationUrl"])
	assert.Equal(t, "https://auth.example.com/token", result["tokenUrl"])
	assert.Equal(t, "client-123", result["x-scalar-client-id"])
	assert.Equal(t, "secret-456", result["clientSecret"])
	assert.Equal(t, "https://app.example.com/callback", result["x-scalar-redirect-uri"])
	assert.Equal(t, "S256", result["x-usePkce"])
	assert.Equal(t, "access_token", result["x-tokenName"])
	assert.Equal(t, "header", result["x-scalar-credentials-location"])
}

// Integration Tests

func TestCompleteAuthenticationConfiguration(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)

	// Configure multiple schemes
	scalargo.WithSecurityScheme("api_key",
		scalargo.APIKeyScheme("X-API-Key", scalargo.APIKeyLocationHeader, "key-123"),
	)(auth)
	scalargo.WithSecurityScheme("bearer_auth",
		scalargo.BearerScheme("token-456"),
	)(auth)
	scalargo.WithSecurityScheme("oauth2",
		scalargo.OAuth2Scheme(
			scalargo.OAuth2FlowAuthorizationCode,
			scalargo.OAuth2Config{
				AuthorizationURL: "https://auth.example.com/authorize",
				TokenURL:         "https://auth.example.com/token",
				ClientID:         "my-client",
				UsePKCE:          scalargo.PKCES256,
				SelectedScopes:   []string{"read:api", "write:api"},
			},
		),
	)(auth)

	// Set preferred scheme
	scalargo.WithPreferredSecurityScheme("bearer_auth")(auth)

	// Verify all schemes exist
	schemes, ok := auth["securitySchemes"].(map[string]any)
	require.True(t, ok)
	assert.Len(t, schemes, 3)

	// Verify preferred scheme
	preferred := auth["preferredSecurityScheme"]
	assert.Equal(t, []any{"bearer_auth"}, preferred)
}
