package scalargo

type AuthenticationOption map[string]any

type AuthOption func(AuthenticationOption)

// APIKeyLocation specifies where the API key is transmitted
type APIKeyLocation string

const (
	APIKeyLocationHeader APIKeyLocation = "header"
	APIKeyLocationQuery  APIKeyLocation = "query"
	APIKeyLocationCookie APIKeyLocation = "cookie"
)

// OAuth2FlowType represents different OAuth2 flow types
type OAuth2FlowType string

const (
	OAuth2FlowAuthorizationCode OAuth2FlowType = "authorizationCode"
	OAuth2FlowClientCredentials OAuth2FlowType = "clientCredentials"
	OAuth2FlowImplicit          OAuth2FlowType = "implicit"
	OAuth2FlowPassword          OAuth2FlowType = "password"
)

// OAuth2PKCEMode represents PKCE configuration modes
type OAuth2PKCEMode string

const (
	PKCES256     OAuth2PKCEMode = "S256"
	PKCEPlain    OAuth2PKCEMode = "plain"
	PKCEDisabled OAuth2PKCEMode = "no"
)

// OAuth2CredentialsLocation specifies where credentials are sent
type OAuth2CredentialsLocation string

const (
	OAuth2CredentialsHeader OAuth2CredentialsLocation = "header"
	OAuth2CredentialsBody   OAuth2CredentialsLocation = "body"
)

// APIKeyConfig holds API key configuration
type APIKeyConfig struct {
	Token string         `json:"token"`
	Name  string         `json:"name,omitempty"`
	In    APIKeyLocation `json:"in,omitempty"`
}

// OAuth2Config holds OAuth2 flow configuration
type OAuth2Config struct {
	AuthorizationURL    string                    `json:"authorizationUrl,omitempty"`
	TokenURL            string                    `json:"tokenUrl,omitempty"`
	ClientID            string                    `json:"x-scalar-client-id,omitempty"`
	ClientSecret        string                    `json:"clientSecret,omitempty"`
	RedirectURI         string                    `json:"x-scalar-redirect-uri,omitempty"`
	UsePKCE             OAuth2PKCEMode            `json:"x-usePkce,omitempty"`
	SelectedScopes      []string                  `json:"selectedScopes,omitempty"`
	SecurityQuery       map[string]string         `json:"x-scalar-security-query,omitempty"`
	SecurityBody        map[string]string         `json:"x-scalar-security-body,omitempty"`
	TokenName           string                    `json:"x-tokenName,omitempty"`
	CredentialsLocation OAuth2CredentialsLocation `json:"x-scalar-credentials-location,omitempty"`
}

// SecuritySchemeConfig represents configuration for a single security scheme
type SecuritySchemeConfig map[string]any

// APIKeyOption configures API key authentication
type APIKeyOption func(*APIKeyConfig)

// OAuth2Option configures OAuth2 authentication
type OAuth2Option func(*OAuth2Config)

// WithCustomSecurity sets the custom security toggle to true
func WithCustomSecurity() AuthOption {
	return func(o AuthenticationOption) {
		o["customSecurity"] = true
	}
}

// WithPreferredSecurityScheme sets the preferred security scheme
// Acceptable values:
// 1. Single security scheme:     "my_custom_security_scheme"
// 2. Multiple security schemes:  "my_custom_security_scheme", "another_security_scheme"
// 3. Complex security schemes:   ["my_custom_security_scheme", "another_security_scheme"], "yet-another_security_scheme"
func WithPreferredSecurityScheme(schemes ...any) AuthOption {
	return func(o AuthenticationOption) {
		o["preferredSecurityScheme"] = schemes
	}
}

// WithHTTPBasicAuth sets the HTTP Basic Auth options
func WithHTTPBasicAuth(username, password string) AuthOption {
	return func(o AuthenticationOption) {
		// Initialize securitySchemes if it doesn't exist
		if o["securitySchemes"] == nil {
			o["securitySchemes"] = make(map[string]any)
		}

		// Type assert and add httpBasic scheme
		if schemes, ok := o["securitySchemes"].(map[string]any); ok {
			schemes["httpBasic"] = map[string]any{
				"username": username,
				"password": password,
			}
		}
	}
}

// WithHTTPBearerToken sets the HTTP Bearer Token options
func WithHTTPBearerToken(token string) AuthOption {
	return func(o AuthenticationOption) {
		// Initialize securitySchemes if it doesn't exist
		if o["securitySchemes"] == nil {
			o["securitySchemes"] = make(map[string]any)
		}

		// Type assert and add httpBearer scheme
		if schemes, ok := o["securitySchemes"].(map[string]any); ok {
			schemes["httpBearer"] = map[string]any{
				"token": token,
			}
		}
	}
}

// WithAPIKeyName sets the parameter name for the API key
func WithAPIKeyName(name string) APIKeyOption {
	return func(c *APIKeyConfig) {
		c.Name = name
	}
}

// WithAPIKeyLocation sets where the API key is sent
func WithAPIKeyLocation(location APIKeyLocation) APIKeyOption {
	return func(c *APIKeyConfig) {
		c.In = location
	}
}

// WithAPIKey configures API key authentication with optional configuration
func WithAPIKey(token string, opts ...APIKeyOption) AuthOption {
	config := &APIKeyConfig{
		Token: token,
		In:    APIKeyLocationHeader, // default
		Name:  "Authorization",      // default for header
	}

	for _, opt := range opts {
		opt(config)
	}

	// Adjust default name based on location if not explicitly customized
	if config.In == APIKeyLocationQuery && config.Name == "Authorization" {
		config.Name = "api_key"
	}
	if config.In == APIKeyLocationCookie && config.Name == "Authorization" {
		config.Name = "api_key"
	}

	return func(o AuthenticationOption) {
		apiKeyMap := map[string]any{
			"token": config.Token,
		}
		if config.Name != "" {
			apiKeyMap["name"] = config.Name
		}
		if config.In != "" {
			apiKeyMap["in"] = string(config.In)
		}
		o["apiKey"] = apiKeyMap
	}
}

// WithAPIKeyHeader is a shorthand for header-based API keys
func WithAPIKeyHeader(name, token string) AuthOption {
	return WithAPIKey(token, WithAPIKeyName(name), WithAPIKeyLocation(APIKeyLocationHeader))
}

// WithAPIKeyQuery is a shorthand for query parameter-based API keys
func WithAPIKeyQuery(name, token string) AuthOption {
	return WithAPIKey(token, WithAPIKeyName(name), WithAPIKeyLocation(APIKeyLocationQuery))
}

// WithAPIKeyCookie is a shorthand for cookie-based API keys
func WithAPIKeyCookie(name, token string) AuthOption {
	return WithAPIKey(token, WithAPIKeyName(name), WithAPIKeyLocation(APIKeyLocationCookie))
}

// OAuth2 Configuration Options

// WithOAuth2ClientID sets the OAuth2 client ID
func WithOAuth2ClientID(clientID string) OAuth2Option {
	return func(c *OAuth2Config) {
		c.ClientID = clientID
	}
}

// WithOAuth2ClientSecret sets the OAuth2 client secret
func WithOAuth2ClientSecret(clientSecret string) OAuth2Option {
	return func(c *OAuth2Config) {
		c.ClientSecret = clientSecret
	}
}

// WithOAuth2RedirectURI sets the OAuth2 redirect URI
func WithOAuth2RedirectURI(redirectURI string) OAuth2Option {
	return func(c *OAuth2Config) {
		c.RedirectURI = redirectURI
	}
}

// WithOAuth2PKCE enables PKCE with the specified mode
func WithOAuth2PKCE(mode OAuth2PKCEMode) OAuth2Option {
	return func(c *OAuth2Config) {
		c.UsePKCE = mode
	}
}

// WithOAuth2Scopes sets the pre-selected OAuth2 scopes
func WithOAuth2Scopes(scopes ...string) OAuth2Option {
	return func(c *OAuth2Config) {
		c.SelectedScopes = scopes
	}
}

// WithOAuth2CustomToken sets a custom token field name
func WithOAuth2CustomToken(tokenName string) OAuth2Option {
	return func(c *OAuth2Config) {
		c.TokenName = tokenName
	}
}

// WithOAuth2AdditionalAuthParams adds additional query parameters to the authorization endpoint
func WithOAuth2AdditionalAuthParams(params map[string]string) OAuth2Option {
	return func(c *OAuth2Config) {
		c.SecurityQuery = params
	}
}

// WithOAuth2AdditionalTokenParams adds additional body parameters to the token endpoint
func WithOAuth2AdditionalTokenParams(params map[string]string) OAuth2Option {
	return func(c *OAuth2Config) {
		c.SecurityBody = params
	}
}

// WithOAuth2CredentialsLocation sets where client credentials are sent
func WithOAuth2CredentialsLocation(location OAuth2CredentialsLocation) OAuth2Option {
	return func(c *OAuth2Config) {
		c.CredentialsLocation = location
	}
}

// WithOAuth2AuthorizationCode configures OAuth2 Authorization Code flow
func WithOAuth2AuthorizationCode(authorizationURL, tokenURL string, opts ...OAuth2Option) AuthOption {
	config := &OAuth2Config{
		AuthorizationURL: authorizationURL,
		TokenURL:         tokenURL,
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(o AuthenticationOption) {
		// Initialize securitySchemes if it doesn't exist
		if o["securitySchemes"] == nil {
			o["securitySchemes"] = make(map[string]any)
		}

		if schemes, ok := o["securitySchemes"].(map[string]any); ok {
			schemes["oauth2"] = map[string]any{
				"flows": map[string]any{
					"authorizationCode": config,
				},
			}
		}
	}
}

// WithOAuth2ClientCredentials configures OAuth2 Client Credentials flow
func WithOAuth2ClientCredentials(tokenURL string, opts ...OAuth2Option) AuthOption {
	config := &OAuth2Config{
		TokenURL: tokenURL,
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(o AuthenticationOption) {
		// Initialize securitySchemes if it doesn't exist
		if o["securitySchemes"] == nil {
			o["securitySchemes"] = make(map[string]any)
		}

		if schemes, ok := o["securitySchemes"].(map[string]any); ok {
			schemes["oauth2"] = map[string]any{
				"flows": map[string]any{
					"clientCredentials": config,
				},
			}
		}
	}
}

// WithOAuth2Implicit configures OAuth2 Implicit flow (deprecated but supported for compatibility)
func WithOAuth2Implicit(authorizationURL string, opts ...OAuth2Option) AuthOption {
	config := &OAuth2Config{
		AuthorizationURL: authorizationURL,
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(o AuthenticationOption) {
		// Initialize securitySchemes if it doesn't exist
		if o["securitySchemes"] == nil {
			o["securitySchemes"] = make(map[string]any)
		}

		if schemes, ok := o["securitySchemes"].(map[string]any); ok {
			schemes["oauth2"] = map[string]any{
				"flows": map[string]any{
					"implicit": config,
				},
			}
		}
	}
}

// WithOAuth2Password configures OAuth2 Password flow (deprecated but supported for compatibility)
func WithOAuth2Password(tokenURL string, opts ...OAuth2Option) AuthOption {
	config := &OAuth2Config{
		TokenURL: tokenURL,
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(o AuthenticationOption) {
		// Initialize securitySchemes if it doesn't exist
		if o["securitySchemes"] == nil {
			o["securitySchemes"] = make(map[string]any)
		}

		if schemes, ok := o["securitySchemes"].(map[string]any); ok {
			schemes["oauth2"] = map[string]any{
				"flows": map[string]any{
					"password": config,
				},
			}
		}
	}
}

// Security Scheme Configuration

// WithSecurityScheme adds a named security scheme configuration
func WithSecurityScheme(name string, config SecuritySchemeConfig) AuthOption {
	return func(o AuthenticationOption) {
		// Initialize securitySchemes if it doesn't exist
		if o["securitySchemes"] == nil {
			o["securitySchemes"] = make(map[string]any)
		}

		if schemes, ok := o["securitySchemes"].(map[string]any); ok {
			schemes[name] = config
		}
	}
}

// APIKeyScheme creates an API key security scheme configuration
func APIKeyScheme(name string, location APIKeyLocation, value string) SecuritySchemeConfig {
	return SecuritySchemeConfig{
		"type":  "apiKey",
		"name":  name,
		"in":    string(location),
		"value": value,
	}
}

// BearerScheme creates a bearer token security scheme configuration
func BearerScheme(token string) SecuritySchemeConfig {
	return SecuritySchemeConfig{
		"type":   "http",
		"scheme": "bearer",
		"token":  token,
	}
}

// BasicScheme creates a basic auth security scheme configuration
func BasicScheme(username, password string) SecuritySchemeConfig {
	return SecuritySchemeConfig{
		"type":     "http",
		"scheme":   "basic",
		"username": username,
		"password": password,
	}
}

// OAuth2Scheme creates an OAuth2 security scheme configuration
func OAuth2Scheme(flow OAuth2FlowType, config OAuth2Config) SecuritySchemeConfig {
	return SecuritySchemeConfig{
		"type": "oauth2",
		"flows": map[string]any{
			string(flow): config,
		},
	}
}
