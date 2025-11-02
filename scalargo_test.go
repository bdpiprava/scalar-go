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
