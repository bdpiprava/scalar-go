package scalargo_test

import (
	"testing"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/stretchr/testify/assert"
)

func TestWithHiddenClients(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		options           []scalargo.Option
		wantHiddenClients any
	}{
		{
			name: "should set hidden clients with single client",
			options: []scalargo.Option{
				scalargo.WithHiddenClients("curl"),
			},
			wantHiddenClients: []string{"curl"},
		},
		{
			name: "should set hidden clients with multiple clients",
			options: []scalargo.Option{
				scalargo.WithHiddenClients("curl", "wget", "httpie"),
			},
			wantHiddenClients: []string{"curl", "wget", "httpie"},
		},
		{
			name: "should set hidden clients to nil slice when no clients provided",
			options: []scalargo.Option{
				scalargo.WithHiddenClients(),
			},
			wantHiddenClients: []string(nil),
		},
		{
			name: "should overwrite previous hidden clients when called twice",
			options: []scalargo.Option{
				scalargo.WithHiddenClients("curl"),
				scalargo.WithHiddenClients("wget", "httpie"),
			},
			wantHiddenClients: []string{"wget", "httpie"},
		},
		{
			name: "should ignore WithHiddenClients when WithHideAllClients is called first",
			options: []scalargo.Option{
				scalargo.WithHideAllClients(),
				scalargo.WithHiddenClients("curl", "wget"),
			},
			wantHiddenClients: true,
		},
		{
			name: "should be overridden by WithHideAllClients when called after",
			options: []scalargo.Option{
				scalargo.WithHiddenClients("curl", "wget"),
				scalargo.WithHideAllClients(),
			},
			wantHiddenClients: true,
		},
		{
			name: "should handle WithHiddenClients called between WithHideAllClients calls",
			options: []scalargo.Option{
				scalargo.WithHideAllClients(),
				scalargo.WithHiddenClients("curl"),
				scalargo.WithHideAllClients(),
			},
			wantHiddenClients: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{
				Configurations: make(map[string]any),
			}

			for _, opt := range tc.options {
				opt(opts)
			}

			got := opts.Configurations["hiddenClients"]
			assert.Equal(t, tc.wantHiddenClients, got)
		})
	}
}

func TestWithCDN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cdn     string
		wantCDN string
	}{
		{
			name:    "should set custom CDN URL",
			cdn:     "https://custom.cdn.com/scalar",
			wantCDN: "https://custom.cdn.com/scalar",
		},
		{
			name:    "should set empty CDN",
			cdn:     "",
			wantCDN: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{}
			scalargo.WithCDN(tc.cdn)(opts)

			assert.Equal(t, tc.wantCDN, opts.CDN)
		})
	}
}

func TestWithProxy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		proxy     string
		wantProxy any
	}{
		{
			name:      "should set proxy URL",
			proxy:     "https://proxy.example.com",
			wantProxy: "https://proxy.example.com",
		},
		{
			name:      "should set empty proxy",
			proxy:     "",
			wantProxy: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{
				Configurations: make(map[string]any),
			}
			scalargo.WithProxy(tc.proxy)(opts)

			assert.Equal(t, tc.wantProxy, opts.Configurations["proxy"])
		})
	}
}

func TestWithEditable(t *testing.T) {
	t.Parallel()

	opts := &scalargo.Options{
		Configurations: make(map[string]any),
	}
	scalargo.WithEditable()(opts)

	assert.Equal(t, true, opts.Configurations["isEditable"])
}

func TestWithSidebarVisibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		visible bool
		want    bool
	}{
		{
			name:    "should set sidebar visible",
			visible: true,
			want:    true,
		},
		{
			name:    "should set sidebar hidden",
			visible: false,
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{
				Configurations: make(map[string]any),
			}
			scalargo.WithSidebarVisibility(tc.visible)(opts)

			assert.Equal(t, tc.want, opts.Configurations["showSidebar"])
		})
	}
}

func TestWithHideModels(t *testing.T) {
	t.Parallel()

	opts := &scalargo.Options{
		Configurations: make(map[string]any),
	}
	scalargo.WithHideModels()(opts)

	assert.Equal(t, true, opts.Configurations["hideModels"])
}

func TestWithHideDownloadButton(t *testing.T) {
	t.Parallel()

	opts := &scalargo.Options{
		Configurations: make(map[string]any),
	}
	scalargo.WithHideDownloadButton()(opts)

	assert.Equal(t, true, opts.Configurations["hideDownloadButton"])
}

func TestWithDarkMode(t *testing.T) {
	t.Parallel()

	opts := &scalargo.Options{
		Configurations: make(map[string]any),
	}
	scalargo.WithDarkMode()(opts)

	assert.Equal(t, true, opts.Configurations["darkMode"])
}

func TestWithForceDarkMode(t *testing.T) {
	t.Parallel()

	opts := &scalargo.Options{
		Configurations: make(map[string]any),
	}
	scalargo.WithForceDarkMode()(opts)

	assert.Equal(t, true, opts.Configurations["forceDarkModeState"])
}

func TestWithHideDarkModeToggle(t *testing.T) {
	t.Parallel()

	opts := &scalargo.Options{
		Configurations: make(map[string]any),
	}
	scalargo.WithHideDarkModeToggle()(opts)

	assert.Equal(t, true, opts.Configurations["hideDarkModeToggle"])
}

func TestWithSearchHotKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		hotKey string
		want   any
	}{
		{
			name:   "should set custom search hotkey",
			hotKey: "ctrl+k",
			want:   "ctrl+k",
		},
		{
			name:   "should set cmd+k hotkey",
			hotKey: "cmd+k",
			want:   "cmd+k",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{
				Configurations: make(map[string]any),
			}
			scalargo.WithSearchHotKey(tc.hotKey)(opts)

			assert.Equal(t, tc.want, opts.Configurations["searchHotKey"])
		})
	}
}

func TestWithHideAllClients(t *testing.T) {
	t.Parallel()

	opts := &scalargo.Options{
		Configurations: make(map[string]any),
	}
	scalargo.WithHideAllClients()(opts)

	assert.Equal(t, true, opts.Configurations["hiddenClients"])
}

func TestWithOverrideCSS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		css  string
		want string
	}{
		{
			name: "should set custom CSS",
			css:  ".custom { color: red; }",
			want: ".custom { color: red; }",
		},
		{
			name: "should set empty CSS",
			css:  "",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{}
			scalargo.WithOverrideCSS(tc.css)(opts)

			assert.Equal(t, tc.want, opts.OverrideCSS)
		})
	}
}

func TestWithAuthentication(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		authentication string
		want           any
	}{
		{
			name:           "should set authentication string",
			authentication: "bearer",
			want:           "bearer",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{
				Configurations: make(map[string]any),
			}
			scalargo.WithAuthentication(tc.authentication)(opts)

			assert.Equal(t, tc.want, opts.Configurations["authentication"])
		})
	}
}

func TestWithPathRouting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		pathRouting string
		want        any
	}{
		{
			name:        "should set path routing",
			pathRouting: "/api-docs",
			want:        "/api-docs",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{
				Configurations: make(map[string]any),
			}
			scalargo.WithPathRouting(tc.pathRouting)(opts)

			assert.Equal(t, tc.want, opts.Configurations["pathRouting"])
		})
	}
}

func TestWithBaseServerURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  string
		want any
	}{
		{
			name: "should set base server URL",
			url:  "https://api.example.com",
			want: "https://api.example.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{
				Configurations: make(map[string]any),
			}
			scalargo.WithBaseServerURL(tc.url)(opts)

			assert.Equal(t, tc.want, opts.Configurations["baseServerUrl"])
		})
	}
}

func TestWithDefaultFonts(t *testing.T) {
	t.Parallel()

	opts := &scalargo.Options{
		Configurations: make(map[string]any),
	}
	scalargo.WithDefaultFonts()(opts)

	assert.Equal(t, true, opts.Configurations["withDefaultFonts"])
}

func TestWithBaseFileName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		baseFileName string
		want         string
	}{
		{
			name:         "should set custom base filename",
			baseFileName: "openapi.yaml",
			want:         "openapi.yaml",
		},
		{
			name:         "should set json filename",
			baseFileName: "api.json",
			want:         "api.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{}
			scalargo.WithBaseFileName(tc.baseFileName)(opts)

			assert.Equal(t, tc.want, opts.BaseFileName)
		})
	}
}

func TestWithSpecModifier(t *testing.T) {
	t.Parallel()

	// SpecModifier is func(*model.Spec) *model.Spec
	// The integration tests in scalargo_test.go already verify this works correctly
	// Just verify the option function doesn't panic
	opts := &scalargo.Options{}
	assert.NotPanics(t, func() {
		scalargo.WithSpecModifier(nil)(opts)
	})
	assert.Nil(t, opts.SpecModifier)
}

func TestWithSpecDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		specDir string
		want    string
	}{
		{
			name:    "should set spec directory",
			specDir: "./specs",
			want:    "./specs",
		},
		{
			name:    "should set absolute path",
			specDir: "/usr/local/specs",
			want:    "/usr/local/specs",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{}
			scalargo.WithSpecDir(tc.specDir)(opts)

			assert.Equal(t, tc.want, opts.SpecDirectory)
		})
	}
}

func TestWithSpecURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		specURL string
		want    string
	}{
		{
			name:    "should set spec URL",
			specURL: "https://api.example.com/openapi.yaml",
			want:    "https://api.example.com/openapi.yaml",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{}
			scalargo.WithSpecURL(tc.specURL)(opts)

			assert.Equal(t, tc.want, opts.SpecURL)
		})
	}
}

func TestWithSpecBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		specBytes []byte
		want      []byte
	}{
		{
			name:      "should set spec bytes",
			specBytes: []byte("openapi: 3.0.0"),
			want:      []byte("openapi: 3.0.0"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{}
			scalargo.WithSpecBytes(tc.specBytes)(opts)

			assert.Equal(t, tc.want, opts.SpecBytes)
		})
	}
}

func TestWithHTTPBasicAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		username string
		password string
	}{
		{
			name:     "should set HTTP Basic Auth credentials",
			username: "admin",
			password: "secret123",
		},
		{
			name:     "should set HTTP Basic Auth with empty username",
			username: "",
			password: "password",
		},
		{
			name:     "should set HTTP Basic Auth with empty password",
			username: "user",
			password: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			auth := make(scalargo.AuthenticationOption)
			scalargo.WithHTTPBasicAuth(tc.username, tc.password)(auth)

			schemes, ok := auth["securitySchemes"].(map[string]any)
			assert.True(t, ok, "securitySchemes should exist and be a map")

			httpBasic, ok := schemes["httpBasic"].(map[string]any)
			assert.True(t, ok, "httpBasic scheme should exist")
			assert.Equal(t, tc.username, httpBasic["username"])
			assert.Equal(t, tc.password, httpBasic["password"])
		})
	}
}

func TestWithHTTPBasicAuth_InitializesSecuritySchemes(t *testing.T) {
	t.Parallel()

	auth := make(scalargo.AuthenticationOption)
	assert.Nil(t, auth["securitySchemes"], "securitySchemes should be nil initially")

	scalargo.WithHTTPBasicAuth("user", "pass")(auth)

	assert.NotNil(t, auth["securitySchemes"], "securitySchemes should be initialized")
}

func TestWithAPIKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "should set API key token",
			token: "sk-1234567890abcdef",
		},
		{
			name:  "should set empty API key token",
			token: "",
		},
		{
			name:  "should set API key with special characters",
			token: "key!@#$%^&*()",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			auth := make(scalargo.AuthenticationOption)
			scalargo.WithAPIKey(tc.token)(auth)

			apiKey, ok := auth["apiKey"].(map[string]any)
			assert.True(t, ok, "apiKey should exist and be a map")
			assert.Equal(t, tc.token, apiKey["token"])
		})
	}
}

func TestWithServers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		servers []scalargo.ServerOverride
		want    []scalargo.ServerOverride
	}{
		{
			name: "should set single server override",
			servers: []scalargo.ServerOverride{
				{URL: "https://api.example.com", Description: "Production server"},
			},
			want: []scalargo.ServerOverride{
				{URL: "https://api.example.com", Description: "Production server"},
			},
		},
		{
			name: "should set multiple server overrides",
			servers: []scalargo.ServerOverride{
				{URL: "https://api.example.com", Description: "Production"},
				{URL: "https://staging.example.com", Description: "Staging"},
				{URL: "http://localhost:8080", Description: "Development"},
			},
			want: []scalargo.ServerOverride{
				{URL: "https://api.example.com", Description: "Production"},
				{URL: "https://staging.example.com", Description: "Staging"},
				{URL: "http://localhost:8080", Description: "Development"},
			},
		},
		{
			name:    "should set empty servers slice",
			servers: []scalargo.ServerOverride{},
			want:    []scalargo.ServerOverride{},
		},
		{
			name: "should set server with empty description",
			servers: []scalargo.ServerOverride{
				{URL: "https://api.example.com", Description: ""},
			},
			want: []scalargo.ServerOverride{
				{URL: "https://api.example.com", Description: ""},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{
				Configurations: make(map[string]any),
			}
			scalargo.WithServers(tc.servers...)(opts)

			got := opts.Configurations["servers"]
			assert.Equal(t, tc.want, got)
		})
	}
}
