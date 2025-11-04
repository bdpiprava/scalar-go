package scalargo

import (
	"strings"
	"testing"

	"github.com/bdpiprava/scalar-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithRenderMode(t *testing.T) {
	t.Run("should default to JavaScript API mode", func(t *testing.T) {
		options := buildOptions()
		assert.Equal(t, RenderModeJavaScriptAPI, options.RenderMode)
	})

	t.Run("should allow setting data-attribute mode for backward compatibility", func(t *testing.T) {
		options := buildOptions(WithRenderMode(RenderModeDataAttribute))
		assert.Equal(t, RenderModeDataAttribute, options.RenderMode)
	})

	t.Run("should explicitly set JavaScript API mode", func(t *testing.T) {
		options := buildOptions(WithRenderMode(RenderModeJavaScriptAPI))
		assert.Equal(t, RenderModeJavaScriptAPI, options.RenderMode)
	})
}

func TestRenderMode_DataAttribute(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeDataAttribute),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `<script id="api-reference"`)
	assert.Contains(t, html, `data-url="https://example.com/openapi.json"`)
	assert.Contains(t, html, `data-configuration=`)
	assert.NotContains(t, html, `<div id="api-reference"></div>`)
	assert.NotContains(t, html, `Scalar.createApiReference`)
}

func TestRenderMode_JavaScriptAPI_WithURL(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `<div id="app"></div>`)
	assert.Contains(t, html, `Scalar.createApiReference('#app',`)
	assert.Contains(t, html, `"url":"https://example.com/openapi.json"`)
	assert.NotContains(t, html, `<script id="api-reference"`)
	assert.NotContains(t, html, `data-url=`)
}

func TestRenderMode_JavaScriptAPI_WithDirectory(t *testing.T) {
	html, err := NewV2(
		WithSpecDir("./data/loader"),
		WithBaseFileName("pet-store.yml"),
		WithRenderMode(RenderModeJavaScriptAPI),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `<div id="app"></div>`)
	assert.Contains(t, html, `Scalar.createApiReference('#app',`)
	assert.Contains(t, html, `"content":"{`)     // Verify content field with JSON
	assert.Contains(t, html, `Swagger Petstore`) // Verify spec title is in the HTML
	assert.NotContains(t, html, `<script id="api-reference"`)
}

func TestWithCustomHeadJS(t *testing.T) {
	customJS := `console.log('Custom head script');`

	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithCustomHeadJS(customJS),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `<script>console.log('Custom head script');</script>`)
	// Should be in head, before CDN
	headIdx := strings.Index(html, customJS)
	cdnIdx := strings.Index(html, DefaultCDN)
	assert.Less(t, headIdx, cdnIdx, "Custom head JS should appear before CDN script")
}

func TestWithCustomBodyJS(t *testing.T) {
	customJS := `console.log('Custom body script');`

	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithCustomBodyJS(customJS),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `<script>console.log('Custom body script');</script>`)
	// Should be in body, after CDN
	cdnIdx := strings.Index(html, DefaultCDN)
	bodyJSIdx := strings.Index(html, customJS)
	assert.Greater(t, bodyJSIdx, cdnIdx, "Custom body JS should appear after CDN script")
}

func TestWithCustomJS_BothHeadAndBody(t *testing.T) {
	headJS := `console.log('head');`
	bodyJS := `console.log('body');`

	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithCustomHeadJS(headJS),
		WithCustomBodyJS(bodyJS),
	)

	require.NoError(t, err)
	assert.Contains(t, html, headJS)
	assert.Contains(t, html, bodyJS)

	headIdx := strings.Index(html, headJS)
	bodyIdx := strings.Index(html, bodyJS)
	assert.Less(t, headIdx, bodyIdx, "Head JS should appear before body JS")
}

func TestCustomJS_WithJavaScriptAPIMode(t *testing.T) {
	headJS := `const config = { debug: true };`
	bodyJS := `document.addEventListener('scalar:loaded', () => console.log('Loaded'));`

	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithCustomHeadJS(headJS),
		WithCustomBodyJS(bodyJS),
	)

	require.NoError(t, err)
	assert.Contains(t, html, headJS)
	assert.Contains(t, html, bodyJS)
	assert.Contains(t, html, `Scalar.createApiReference`)

	// Verify order: head JS -> CDN -> init script -> body JS
	headIdx := strings.Index(html, headJS)
	cdnIdx := strings.Index(html, DefaultCDN)
	initIdx := strings.Index(html, `Scalar.createApiReference`)
	bodyIdx := strings.Index(html, bodyJS)

	assert.Less(t, headIdx, cdnIdx, "Head JS before CDN")
	assert.Less(t, cdnIdx, initIdx, "CDN before init script")
	assert.Less(t, initIdx, bodyIdx, "Init script before body JS")
}

func TestWithHideSearch(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithHideSearch(true),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `"hideSearch":true`)
}

func TestWithShowOperationID(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithShowOperationID(true),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `"showOperationId":true`)
}

func TestWithDefaultHTTPClient(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithDefaultHTTPClient("node", "undici"),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `"defaultHttpClient":{"targetKey":"node","clientKey":"undici"}`)
}

func TestWithTagsSorter(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithTagsSorter(SorterAlpha),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `"tagsSorter":"alpha"`)
}

func TestWithOperationsSorter(t *testing.T) {
	t.Run("alpha sorter", func(t *testing.T) {
		html, err := NewV2(
			WithSpecURL("https://example.com/openapi.json"),
			WithRenderMode(RenderModeJavaScriptAPI),
			WithOperationsSorter(SorterAlpha),
		)

		require.NoError(t, err)
		assert.Contains(t, html, `"operationsSorter":"alpha"`)
	})

	t.Run("method sorter", func(t *testing.T) {
		html, err := NewV2(
			WithSpecURL("https://example.com/openapi.json"),
			WithRenderMode(RenderModeJavaScriptAPI),
			WithOperationsSorter(SorterMethod),
		)

		require.NoError(t, err)
		assert.Contains(t, html, `"operationsSorter":"method"`)
	})
}

func TestWithOperationTitleSource(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithOperationTitleSource(OperationTitleSourcePath),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `"operationTitleSource":"path"`)
}

func TestWithOrderSchemaPropertiesBy(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithOrderSchemaPropertiesBy(SchemaPropertiesOrderAlpha),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `"orderSchemaPropertiesBy":"alpha"`)
}

func TestWithPersistAuth(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithPersistAuth(true),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `"persistAuth":true`)
}

func TestWithCustomCSS(t *testing.T) {
	customCSS := `.scalar-card { border: 2px solid red; }`

	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithCustomCSS(customCSS),
	)

	require.NoError(t, err)
	// customCss goes in config, not in <style> tag
	assert.Contains(t, html, `"customCss":".scalar-card { border: 2px solid red; }"`)
}

func TestWithMultipleSources(t *testing.T) {
	sources := []DocumentSource{
		{
			Title:   "API v1",
			Slug:    "v1",
			URL:     "https://example.com/v1/openapi.json",
			Default: true,
		},
		{
			Title: "API v2",
			Slug:  "v2",
			URL:   "https://example.com/v2/openapi.json",
		},
	}

	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithMultipleSources(sources...),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `"sources":[`)
	assert.Contains(t, html, `"title":"API v1"`)
	assert.Contains(t, html, `"slug":"v1"`)
	assert.Contains(t, html, `"url":"https://example.com/v1/openapi.json"`)
	assert.Contains(t, html, `"default":true`)
	assert.Contains(t, html, `"title":"API v2"`)
}

func TestCombinedNewOptions(t *testing.T) {
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithTheme(ThemePurple),
		WithLayout(LayoutModern),
		WithHideSearch(true),
		WithShowOperationID(true),
		WithDefaultHTTPClient("node", "undici"),
		WithTagsSorter(SorterAlpha),
		WithOperationsSorter(SorterMethod),
		WithPersistAuth(true),
		WithCustomHeadJS(`console.log('init');`),
		WithCustomBodyJS(`console.log('ready');`),
	)

	require.NoError(t, err)

	// Verify all options are present
	assert.Contains(t, html, `Scalar.createApiReference`)
	assert.Contains(t, html, `"theme":"purple"`)
	assert.Contains(t, html, `"layout":"modern"`)
	assert.Contains(t, html, `"hideSearch":true`)
	assert.Contains(t, html, `"showOperationId":true`)
	assert.Contains(t, html, `"defaultHttpClient":{"targetKey":"node","clientKey":"undici"}`)
	assert.Contains(t, html, `"tagsSorter":"alpha"`)
	assert.Contains(t, html, `"operationsSorter":"method"`)
	assert.Contains(t, html, `"persistAuth":true`)
	assert.Contains(t, html, `console.log('init')`)
	assert.Contains(t, html, `console.log('ready')`)
}

func TestBuildInitScript_WithSpecModifier(t *testing.T) {
	html, err := NewV2(
		WithSpecDir("./data/loader"),
		WithBaseFileName("pet-store.yml"),
		WithRenderMode(RenderModeJavaScriptAPI),
		WithSpecModifier(func(spec *model.Spec) *model.Spec {
			spec.Info.Title = "Modified Title"
			return spec
		}),
	)

	require.NoError(t, err)
	assert.Contains(t, html, `"title":"Modified Title"`)
	assert.Contains(t, html, `Scalar.createApiReference`)
}

func TestBackwardCompatibility_DefaultMode(t *testing.T) {
	// When no render mode is specified, should use JavaScript API mode (new default)
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
	)

	require.NoError(t, err)
	// Should use JavaScript API mode (new default)
	assert.Contains(t, html, `<div id="app"></div>`)
	assert.Contains(t, html, `Scalar.createApiReference('#app',`)
	assert.NotContains(t, html, `<script id="api-reference"`)
	assert.NotContains(t, html, `data-url=`)
}

func TestBackwardCompatibility_ExplicitDataAttribute(t *testing.T) {
	// Explicitly requesting data-attribute mode should still work
	html, err := NewV2(
		WithSpecURL("https://example.com/openapi.json"),
		WithRenderMode(RenderModeDataAttribute),
	)

	require.NoError(t, err)
	// Should use data-attribute mode (legacy)
	assert.Contains(t, html, `<script id="api-reference"`)
	assert.Contains(t, html, `data-url=`)
	assert.NotContains(t, html, `Scalar.createApiReference`)
}
