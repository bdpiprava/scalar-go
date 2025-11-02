package scalargo_test

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/model"
)

func Test_ShouldCallOverrideHandler_WhenProvider(t *testing.T) {
	var called bool
	content, err := scalargo.New(
		"./data/loader",
		scalargo.WithBaseFileName("pet-store.yml"),
		scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
			called = true
			require.Equal(t, "Swagger Petstore", spec.Info.Title)

			spec.Info.Title = "PetStore API"
			return spec
		}),
	)

	require.NoError(t, err)
	require.True(t, called)
	require.NotNil(t, content)
	require.Contains(t, content, "<title>PetStore API</title>")
}

func Test_NewV2(t *testing.T) {
	const specURL = "https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"
	testCases := []struct {
		name      string
		inputOpts []scalargo.Option
		asserter  asserter
		wantError string
	}{
		{
			name:      "should return error when no option is provided",
			inputOpts: []scalargo.Option{},
			asserter:  func(t *testing.T, got html) { require.Equal(t, html{}, got) },
			wantError: "one of SpecURL, SpecDirectory or SpecBytes must be configured",
		},
		{
			name:      "should render html containing script with spec URL when spec URL is configured",
			inputOpts: []scalargo.Option{scalargo.WithSpecURL(specURL)},
			asserter: func(t *testing.T, got html) {
				require.Empty(t, got.spec)
				require.Equal(t, "https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml", got.specURL)
			},
		},
		{
			name:      "should render html with inline spec when spec directory is configured",
			inputOpts: []scalargo.Option{scalargo.WithSpecDir("./data/loader"), scalargo.WithBaseFileName("pet-store.yml")},
			asserter: func(t *testing.T, got html) {
				require.Empty(t, got.specURL)
				require.True(t, strings.HasPrefix(got.spec, `{"openapi":"3.0.0","info":{"title":"Swagger Petstore",`))
			},
		},
		{
			name:      "should render html with inline spec when spec bytes is configured",
			inputOpts: []scalargo.Option{scalargo.WithSpecBytes([]byte(`{"openapi":"3.0.0","info":{"title":"Swagger Petstore"}}`))},
			asserter: func(t *testing.T, got html) {
				require.Empty(t, got.specURL)
				require.True(t, strings.HasPrefix(got.spec, `{"openapi":"3.0.0","info":{"title":"Swagger Petstore","version":""},"paths":{}`))
			},
		},
		{
			name: "should render html with authentication configuration",
			inputOpts: []scalargo.Option{
				scalargo.WithSpecURL(specURL),
				scalargo.WithAuthenticationOpts(
					scalargo.WithCustomSecurity(),
					scalargo.WithPreferredSecurityScheme("bearerAuth"),
					scalargo.WithHTTPBearerToken("this-is-a-token"),
				),
			},
			asserter: func(t *testing.T, got html) {
				require.Equal(t, map[string]any{
					"layout":         string(scalargo.LayoutModern),
					"theme":          string(scalargo.ThemeDefault),
					"metadata":       map[string]any{"title": "API Reference"},
					"authentication": `{"customSecurity":true,"http":{"bearer":{"token":"this-is-a-token"}},"preferredSecurityScheme":["bearerAuth"]}`,
				}, got.configuration)
			},
		},
		{
			name: "should render html with custom configuration",
			inputOpts: []scalargo.Option{
				scalargo.WithSpecURL(specURL),
				scalargo.WithTheme(scalargo.ThemeKepler),
				scalargo.WithHideAllClients(),
				scalargo.WithLayout(scalargo.LayoutClassic),
				scalargo.WithMetaDataOpts(scalargo.WithKeyValue("foo", "bar")),
			},
			asserter: func(t *testing.T, got html) {
				require.Equal(t, map[string]any{
					"hiddenClients": true,
					"layout":        string(scalargo.LayoutClassic),
					"theme":         string(scalargo.ThemeKepler),
					"metadata": map[string]any{
						"title": "API Reference",
						"foo":   "bar",
					},
				}, got.configuration)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotContent, gotErr := scalargo.NewV2(tc.inputOpts...)

			content := parseContent(gotContent)

			tc.asserter(t, content)
			if tc.wantError != "" {
				require.ErrorContains(t, gotErr, tc.wantError)
			} else {
				require.NoError(t, gotErr)
			}
		})
	}
}

type asserter func(t *testing.T, got html)

var titleMatcher = regexp.MustCompile(".*<title>(.*)</title>.*")
var overrideCSSMatcher = regexp.MustCompile(".*<style>(.*)</style>.*")
var configurationMatcher = regexp.MustCompile(`.*data-configuration="(.*)".*`)
var specURLMatcher = regexp.MustCompile(`.*id="api-reference".*data-url="(.*?[^\\])".*`)
var specMatcher = regexp.MustCompile(`.*<script.*id="api-reference".*>(.*)</script>.*`)

type html struct {
	title         string
	specURL       string
	spec          string
	overrideCSS   string
	configuration map[string]any
}

func parseContent(content string) html {
	if content == "" {
		return html{}
	}

	configStr := strings.ReplaceAll(getFirstGroup(configurationMatcher, content), "&quot;", `"`)
	var config map[string]any
	if len(configStr) > 0 {
		err := json.Unmarshal([]byte(configStr), &config)
		if err != nil {
			config = map[string]any{}
		}
	} else {
		config = map[string]any{}
	}

	result := &html{
		title:         getFirstGroup(titleMatcher, content),
		overrideCSS:   getFirstGroup(overrideCSSMatcher, content),
		configuration: config,
		specURL:       getFirstGroup(specURLMatcher, content),
		spec:          getFirstGroup(specMatcher, content),
	}

	return *result
}

func getFirstGroup(matcher *regexp.Regexp, content string) string {
	matches := matcher.FindStringSubmatch(content)
	if len(matches) >= 1 {
		return matches[1]
	}
	return ""
}

func Test_XSS_Prevention(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		opts              []scalargo.Option
		maliciousContent  string
		shouldNotContain  []string // XSS payloads that should be escaped
		shouldContain     []string // Expected escaped versions
		description       string
	}{
		{
			name: "should escape XSS in title tag",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				scalargo.WithMetaDataOpts(
					scalargo.WithKeyValue("title", "</title><script>alert('XSS')</script><title>"),
				),
			},
			shouldNotContain: []string{
				"</title><script>alert('XSS')</script><title>",
				"<script>alert('XSS')</script>",
			},
			shouldContain: []string{
				"&lt;/title&gt;&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;&lt;title&gt;",
			},
			description: "Title should be HTML-escaped to prevent breaking out of title tag",
		},
		{
			name: "should sanitize XSS in CSS override",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				scalargo.WithOverrideCSS("</style><script>alert('CSS XSS')</script><style>body { color: red; }"),
			},
			shouldNotContain: []string{
				"<script>alert('CSS XSS')</script>",                                    // Script tag should be removed
				"</style><script>",                                                    // Must not allow breaking out of style
				"<style></style>",                                                     // Empty style tags from injection
			},
			shouldContain: []string{
				"body { color: red; }",                                               // Valid CSS preserved
				// Note: "alert('CSS XSS')" text remains but is harmless - it's invalid CSS that browsers ignore
			},
			description: "CSS should have HTML tags stripped to prevent breaking out of style tag",
		},
		{
			name: "should escape XSS in CDN URL",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				scalargo.WithCDN("\" onerror=\"alert('CDN XSS')"),
			},
			shouldNotContain: []string{
				"\" onerror=\"alert('CDN XSS')",
				"onerror=\"alert('CDN XSS')",
				"<script src=\"\" onerror=\"alert('CDN XSS')\">", // Must not allow attribute injection
			},
			description: "CDN URL should be escaped to prevent attribute injection",
		},
		{
			name: "should handle multiple XSS vectors simultaneously",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				scalargo.WithMetaDataOpts(
					scalargo.WithKeyValue("title", "<script>alert(1)</script>"),
				),
				scalargo.WithOverrideCSS("</style><img src=x onerror=alert(2)><style>body { color: blue; }"),
				scalargo.WithCDN("javascript:alert(3)"),
			},
			shouldNotContain: []string{
				"<script>alert(1)</script>",  // Title XSS should be escaped
				"</style><img src=x onerror=alert(2)><style>",  // CSS HTML tags should be removed
				"<img src=x onerror=alert(2)>",  // IMG tag should be removed
				"javascript:alert(3)",  // Dangerous JS URL should be sanitized by template
			},
			shouldContain: []string{
				"&lt;script&gt;alert(1)&lt;/script&gt;",  // Escaped title
				"body { color: blue; }",  // CSS preserved, tags removed
			},
			description: "Should sanitize all XSS vectors when multiple are present",
		},
		{
			name: "should escape HTML entities in title",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				scalargo.WithMetaDataOpts(
					scalargo.WithKeyValue("title", "<>&\"'"),
				),
			},
			shouldNotContain: []string{
				"<title><>&\"'</title>",
			},
			shouldContain: []string{
				"&lt;&gt;&amp;&#34;&#39;",
			},
			description: "HTML entities should be properly escaped",
		},
		{
			name: "should handle data URI XSS attempt in CDN",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				scalargo.WithCDN("data:text/javascript,alert('XSS')"),
			},
			shouldNotContain: []string{
				"<script src=\"data:text/javascript,alert('XSS')\"></script>",
			},
			description: "Data URIs should be escaped in CDN URL",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			html, err := scalargo.NewV2(tc.opts...)
			require.NoError(t, err)
			require.NotEmpty(t, html)

			// Verify malicious content is not present in raw form
			for _, malicious := range tc.shouldNotContain {
				require.NotContains(t, html, malicious,
					"HTML should not contain unescaped malicious content: %s", malicious)
			}

			// Verify escaped versions are present (if specified)
			for _, escaped := range tc.shouldContain {
				require.Contains(t, html, escaped,
					"HTML should contain properly escaped content: %s", escaped)
			}
		})
	}
}

func Test_CSS_Sanitization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		maliciousCSS     string
		shouldNotContain []string // Dangerous patterns that must be removed
		shouldContain    []string // Safe CSS that should be preserved
		description      string
	}{
		{
			name:         "Bypass 1: Malformed HTML tags with spaces",
			maliciousCSS: "body { color: red; } < script>alert(1)</script > .test { margin: 0; }",
			shouldNotContain: []string{
				"< script>",
				"</script >",
				"<script>",
				"</script>",
			},
			shouldContain: []string{
				"body { color: red; }",
				".test { margin: 0; }",
			},
			description: "Should remove malformed HTML tags with extra spaces",
		},
		{
			name:         "Bypass 2: HTML entities",
			maliciousCSS: "body { color: red; } &lt;script&gt;alert(1)&lt;/script&gt;",
			shouldNotContain: []string{
				"&lt;",
				"&gt;",
				"&amp;",
				"&#",
			},
			shouldContain: []string{
				"body { color: red; }",
			},
			description: "Should remove HTML entities to prevent decoding attacks",
		},
		{
			name:         "Bypass 3: CSS expressions (old IE)",
			maliciousCSS: "body { width: expression(alert(1)); color: red; }",
			shouldNotContain: []string{
				"expression(",
				"expression(alert(1))",
			},
			shouldContain: []string{
				"color: red;",
			},
			description: "Should remove CSS expressions",
		},
		{
			name:         "Bypass 4: Data URIs with embedded scripts",
			maliciousCSS: "body { background: url('data:text/html,<script>alert(1)</script>'); color: blue; }",
			shouldNotContain: []string{
				"data:text/html",
				"data:",
				"<script>",
			},
			shouldContain: []string{
				"color: blue;",
				"background: url(",
			},
			description: "Should remove data: URIs from CSS",
		},
		{
			name:         "Bypass 5: JavaScript protocol in URLs",
			maliciousCSS: "a { background: url(javascript:alert(1)); color: green; }",
			shouldNotContain: []string{
				"javascript:",
				"javascript:alert",
			},
			shouldContain: []string{
				"color: green;",
			},
			description: "Should remove javascript: protocol from URLs",
		},
		{
			name:         "Bypass 6: vbscript protocol",
			maliciousCSS: "div { background: url(vbscript:msgbox); }",
			shouldNotContain: []string{
				"vbscript:",
			},
			description: "Should remove vbscript: protocol",
		},
		{
			name:         "Bypass 7: @import with external malicious CSS",
			maliciousCSS: "@import url('http://evil.com/xss.css'); body { color: red; }",
			shouldNotContain: []string{
				"@import",
				"http://evil.com",
			},
			shouldContain: []string{
				"body { color: red; }",
			},
			description: "Should remove @import rules",
		},
		{
			name:         "Bypass 8: Style-breaking sequences",
			maliciousCSS: "body { color: red; } </style><script>alert(1)</script><style> .test { margin: 0; }",
			shouldNotContain: []string{
				"</style>",
				"<script>",
			},
			shouldContain: []string{
				"body { color: red; }",
				".test { margin: 0; }",
			},
			description: "Should remove style-breaking tags",
		},
		{
			name:         "Bypass 9: Event handlers embedded in CSS",
			maliciousCSS: "div { color: red; } onclick=alert(1) .test { margin: 0; }",
			shouldNotContain: []string{
				"onclick=",
				"onerror=",
			},
			shouldContain: []string{
				"div { color: red; }",
				".test { margin: 0; }",
			},
			description: "Should remove event handler attributes",
		},
		{
			name:         "Bypass 10: Multiple attack vectors combined",
			maliciousCSS: "@import 'evil.css'; body { width: expression(alert(1)); background: url(javascript:void(0)); } </style><script>alert(2)</script><style> .test { color: url('data:text/html,<img src=x onerror=alert(3)>'); }",
			shouldNotContain: []string{
				"@import",
				"expression(",
				"javascript:",
				"data:",
				"</style>",
				"<script>",
				"onerror=",
			},
			shouldContain: []string{
				"body {",
				".test {",
			},
			description: "Should handle multiple simultaneous attack vectors",
		},
		{
			name:         "Bypass 11: Case variations",
			maliciousCSS: "body { background: URL(JAVASCRIPT:alert(1)); width: ExPrEsSiOn(alert(2)); } @IMPORT 'evil.css';",
			shouldNotContain: []string{
				"JAVASCRIPT:",
				"ExPrEsSiOn(",
				"@IMPORT",
			},
			description: "Should handle case-insensitive attack patterns",
		},
		{
			name:         "Safe CSS should be preserved",
			maliciousCSS: "body { color: #fff; margin: 0; padding: 10px; } .container { display: flex; background: url('/images/bg.png'); }",
			shouldNotContain: []string{
				// Nothing dangerous here
			},
			shouldContain: []string{
				"color: #fff;",
				"margin: 0;",
				"padding: 10px;",
				"display: flex;",
				"background: url('/images/bg.png');",
			},
			description: "Should preserve legitimate CSS including safe relative URLs",
		},
		{
			name:         "Empty CSS should remain empty",
			maliciousCSS: "",
			shouldNotContain: []string{},
			shouldContain:    []string{},
			description:      "Empty CSS should not cause issues",
		},
		{
			name:             "Whitespace-only CSS should be handled gracefully",
			maliciousCSS:     "   \n\t  ",
			shouldNotContain: []string{},
			shouldContain:    []string{},
			description:      "Whitespace-only CSS should not cause errors (may be stripped by template engine)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			html, err := scalargo.NewV2(
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				scalargo.WithOverrideCSS(tc.maliciousCSS),
			)

			require.NoError(t, err)
			require.NotEmpty(t, html)

			// Extract the CSS from the rendered HTML
			content := parseContent(html)

			// Verify dangerous patterns are not present
			for _, dangerous := range tc.shouldNotContain {
				require.NotContains(t, content.overrideCSS, dangerous,
					"CSS should not contain dangerous pattern: %s\nDescription: %s\nActual CSS: %s",
					dangerous, tc.description, content.overrideCSS)
			}

			// Verify safe CSS is preserved
			for _, safe := range tc.shouldContain {
				require.Contains(t, content.overrideCSS, safe,
					"CSS should preserve safe content: %s\nDescription: %s\nActual CSS: %s",
					safe, tc.description, content.overrideCSS)
			}
		})
	}
}

// Helper function to safely call NewV2 and validate result without panics
func assertNewV2NoPanic(t *testing.T, opts []scalargo.Option, wantTitle string, wantErr bool) {
	t.Helper()

	var html string
	var err error
	require.NotPanics(t, func() {
		html, err = scalargo.NewV2(opts...)
	})

	if wantErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)
		require.Contains(t, html, "<title>"+wantTitle+"</title>")
	}
}

// Helper function to safely call GetSpecScript and validate metadata without panics
func assertGetSpecScriptNoPanic(t *testing.T, opts *scalargo.Options, wantErr bool, validateTitle bool, expectedTitle string) {
	t.Helper()

	var script string
	var err error
	require.NotPanics(t, func() {
		script, err = opts.GetSpecScript()
	})

	if wantErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)
		require.NotEmpty(t, script)

		if validateTitle {
			// Verify the metadata was properly set/fixed
			metadata, ok := opts.Configurations["metadata"].(scalargo.MetaData)
			require.True(t, ok, "metadata should be of type MetaData")
			require.Equal(t, expectedTitle, metadata["title"])
		}
	}
}

func Test_NewV2_WithCorruptedMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		opts      []scalargo.Option
		wantTitle string
		wantErr   bool
	}{
		{
			name: "should handle missing metadata gracefully",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				func(o *scalargo.Options) {
					// Intentionally remove metadata after defaults are set
					delete(o.Configurations, "metadata")
				},
			},
			wantTitle: "API Reference", // Should fall back to default
			wantErr:   false,
		},
		{
			name: "should handle wrong type for metadata",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				func(o *scalargo.Options) {
					// Set metadata to wrong type after defaults
					o.Configurations["metadata"] = "not-a-map"
				},
			},
			wantTitle: "API Reference", // Should fall back to default
			wantErr:   false,
		},
		{
			name: "should handle metadata without title",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				scalargo.WithMetaDataOpts(scalargo.WithKeyValue("description", "Some description")),
				func(o *scalargo.Options) {
					// Remove title from metadata
					if metadata, ok := o.Configurations["metadata"].(scalargo.MetaData); ok {
						delete(metadata, "title")
					}
				},
			},
			wantTitle: "API Reference", // Should fall back to default
			wantErr:   false,
		},
		{
			name: "should handle numeric value for metadata",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				func(o *scalargo.Options) {
					o.Configurations["metadata"] = 12345
				},
			},
			wantTitle: "API Reference",
			wantErr:   false,
		},
		{
			name: "should handle nil metadata",
			opts: []scalargo.Option{
				scalargo.WithSpecURL("https://example.com/api.yaml"),
				func(o *scalargo.Options) {
					o.Configurations["metadata"] = nil
				},
			},
			wantTitle: "API Reference",
			wantErr:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assertNewV2NoPanic(t, tc.opts, tc.wantTitle, tc.wantErr)
		})
	}
}

func Test_GetSpecScript_WithCorruptedMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupOpts     func() *scalargo.Options
		wantErr       bool
		validateTitle bool
		expectedTitle string
	}{
		{
			name: "should handle missing metadata when using spec directory",
			setupOpts: func() *scalargo.Options {
				opts := &scalargo.Options{
					Configurations: make(map[string]any),
					SpecDirectory:  "./data/loader",
					BaseFileName:   "pet-store.yml",
					CDN:            scalargo.DefaultCDN,
				}
				// Intentionally don't set metadata
				return opts
			},
			wantErr:       false,
			validateTitle: true,
			expectedTitle: "Swagger Petstore", // Should use spec title
		},
		{
			name: "should handle wrong type for metadata when using spec directory",
			setupOpts: func() *scalargo.Options {
				opts := &scalargo.Options{
					Configurations: make(map[string]any),
					SpecDirectory:  "./data/loader",
					BaseFileName:   "pet-store.yml",
					CDN:            scalargo.DefaultCDN,
				}
				// Set metadata to wrong type
				opts.Configurations["metadata"] = 12345
				return opts
			},
			wantErr:       false,
			validateTitle: true,
			expectedTitle: "Swagger Petstore", // Should create new metadata and use spec title
		},
		{
			name: "should handle nil metadata when using spec directory",
			setupOpts: func() *scalargo.Options {
				opts := &scalargo.Options{
					Configurations: make(map[string]any),
					SpecDirectory:  "./data/loader",
					BaseFileName:   "pet-store.yml",
					CDN:            scalargo.DefaultCDN,
				}
				opts.Configurations["metadata"] = nil
				return opts
			},
			wantErr:       false,
			validateTitle: true,
			expectedTitle: "Swagger Petstore",
		},
		{
			name: "should handle string metadata when using spec directory",
			setupOpts: func() *scalargo.Options {
				opts := &scalargo.Options{
					Configurations: make(map[string]any),
					SpecDirectory:  "./data/loader",
					BaseFileName:   "pet-store.yml",
					CDN:            scalargo.DefaultCDN,
				}
				opts.Configurations["metadata"] = "invalid-string"
				return opts
			},
			wantErr:       false,
			validateTitle: true,
			expectedTitle: "Swagger Petstore",
		},
		{
			name: "should handle metadata with default title when using spec directory",
			setupOpts: func() *scalargo.Options {
				opts := &scalargo.Options{
					Configurations: map[string]any{
						"metadata": scalargo.MetaData{
							"title": "API Reference", // Default title
						},
					},
					SpecDirectory: "./data/loader",
					BaseFileName:  "pet-store.yml",
					CDN:           scalargo.DefaultCDN,
				}
				return opts
			},
			wantErr:       false,
			validateTitle: true,
			expectedTitle: "Swagger Petstore", // Should replace default with spec title
		},
		{
			name: "should preserve custom title when using spec directory",
			setupOpts: func() *scalargo.Options {
				opts := &scalargo.Options{
					Configurations: map[string]any{
						"metadata": scalargo.MetaData{
							"title": "Custom Title",
						},
					},
					SpecDirectory: "./data/loader",
					BaseFileName:  "pet-store.yml",
					CDN:           scalargo.DefaultCDN,
				}
				return opts
			},
			wantErr:       false,
			validateTitle: true,
			expectedTitle: "Custom Title", // Should preserve custom title
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			opts := tc.setupOpts()
			assertGetSpecScriptNoPanic(t, opts, tc.wantErr, tc.validateTitle, tc.expectedTitle)
		})
	}
}
