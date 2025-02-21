package scalargo

type AuthenticationOption map[string]any

type AuthOption func(AuthenticationOption)

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
		o["http"] = map[string]any{
			"basic": map[string]any{
				"username": username,
				"password": password,
			},
		}
	}
}

// WithHTTPBearerToken sets the HTTP Bearer Token options
func WithHTTPBearerToken(token string) AuthOption {
	return func(o AuthenticationOption) {
		o["http"] = map[string]any{
			"bearer": map[string]any{
				"token": token,
			},
		}
	}
}

// WithAPIKey sets the API Key options
func WithAPIKey(token string) AuthOption {
	return func(o AuthenticationOption) {
		o["apiKey"] = map[string]any{
			"token": token,
		}
	}
}
