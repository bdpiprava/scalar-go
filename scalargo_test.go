package scalargo_test

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/model"
	"github.com/stretchr/testify/assert"
)

func Test_ShouldCallOverrideHandler_WhenProvider(t *testing.T) {
	var called bool
	content, err := scalargo.New(
		"./data/loader",
		scalargo.WithBaseFileName("pet-store.yml"),
		scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
			called = true
			assert.Equal(t, "Swagger Petstore", spec.Info.Title)

			spec.Info.Title = "PetStore API"
			return spec
		}),
	)

	assert.NoError(t, err)
	assert.True(t, called)
	assert.NotNil(t, content)
	assert.Contains(t, content, "<title>PetStore API</title>")
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
			asserter:  func(t *testing.T, got html) { assert.Equal(t, html{}, got) },
			wantError: "SpecURL or SpecDirectory must be configured",
		},
		{
			name:      "should render html containing script with spec URL when spec URL is configured",
			inputOpts: []scalargo.Option{scalargo.WithSpecURL(specURL)},
			asserter: func(t *testing.T, got html) {
				assert.Len(t, got.spec, 0)
				assert.Equal(t, "https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml", got.specURL)
			},
		},
		{
			name:      "should render html with inline spec when spec directory is configured",
			inputOpts: []scalargo.Option{scalargo.WithSpecDir("./data/loader"), scalargo.WithBaseFileName("pet-store.yml")},
			asserter: func(t *testing.T, got html) {
				assert.Len(t, got.specURL, 0)
				assert.True(t, strings.HasPrefix(got.spec, `{"openapi":"3.0.0","info":{"title":"Swagger Petstore",`))
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
				assert.Equal(t, map[string]any{
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
				assert.ErrorContains(t, gotErr, tc.wantError)
			} else {
				assert.NoError(t, gotErr)
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
