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
