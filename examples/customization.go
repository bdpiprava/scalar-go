package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

// ExampleThemeDefault demonstrates the default theme
func ExampleThemeDefault() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeDefault),
	)
}

// ExampleThemeMoon demonstrates the moon theme
func ExampleThemeMoon() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeMoon),
	)
}

// ExampleThemePurple demonstrates the purple theme
func ExampleThemePurple() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemePurple),
	)
}

// ExampleThemeSolarized demonstrates the solarized theme
func ExampleThemeSolarized() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeSolarized),
	)
}

// ExampleThemeAlternate demonstrates the alternate theme
func ExampleThemeAlternate() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeAlternate),
	)
}

// ExampleThemeBluePlanet demonstrates the blue planet theme
func ExampleThemeBluePlanet() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeBluePlanet),
	)
}

// ExampleThemeDeepSpace demonstrates the deep space theme
func ExampleThemeDeepSpace() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeDeepSpace),
	)
}

// ExampleThemeSaturn demonstrates the saturn theme
func ExampleThemeSaturn() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeSaturn),
	)
}

// ExampleThemeKepler demonstrates the kepler theme
func ExampleThemeKepler() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeKepler),
	)
}

// ExampleThemeMars demonstrates the mars theme
func ExampleThemeMars() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeMars),
	)
}

// Layout examples - self-contained spec functions

// ExampleLayoutModern demonstrates the modern layout
func ExampleLayoutModern() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithLayout(scalargo.LayoutModern),
	)
}

// ExampleLayoutClassic demonstrates the classic layout
func ExampleLayoutClassic() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithLayout(scalargo.LayoutClassic),
	)
}

// Visibility examples - self-contained spec functions

// ExampleHideSidebar demonstrates hiding the sidebar for a cleaner look
func ExampleHideSidebar() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithSidebarVisibility(false),
	)
}

// ExampleHideModels demonstrates hiding the models section to focus on endpoints
func ExampleHideModels() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithHideModels(),
	)
}

// ExampleDarkMode demonstrates enabling dark mode by default
func ExampleDarkMode() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithDarkMode(),
	)
}

// ExampleDarkModeOptions demonstrates comprehensive dark mode configuration
func ExampleDarkModeOptions() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		// WithDarkMode: Enable dark mode by default (user can still toggle)
		scalargo.WithDarkMode(),
		// WithHideDarkModeToggle: Hide the dark mode toggle button
		// scalargo.WithHideDarkModeToggle(),
		// WithForceDarkMode: Force dark mode to specific state, preventing user changes
		// scalargo.WithForceDarkMode(),
	)
}

// ExampleHideDownloadButton demonstrates hiding the OpenAPI spec download button
func ExampleHideDownloadButton() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithHideDownloadButton(),
	)
}

// Advanced examples - self-contained spec functions

// ExampleCustomCSS demonstrates custom CSS styling for branded documentation
func ExampleCustomCSS() (string, error) {
	customCSS := `
		/* Custom styling for the API documentation */
		.section-header {
			color: #e74c3c !important;
			font-weight: bold !important;
		}

		.api-client__request {
			background-color: #f8f9fa !important;
			border-left: 4px solid #007bff !important;
		}

		/* Custom button styling */
		button {
			border-radius: 8px !important;
		}
	`

	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithOverrideCSS(customCSS),
	)
}

// Client Visibility Examples

// ExampleHideAllClients demonstrates hiding all client examples
func ExampleHideAllClients() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithHideAllClients(),
	)
}

// ExampleShowOnlyCurlAndFetch demonstrates showing only curl and fetch clients
func ExampleShowOnlyCurlAndFetch() (string, error) {
	// To show only specific clients, hide all others
	hiddenClients := []string{
		scalargo.ClientAsyncHTTP,
		scalargo.ClientAxios,
		scalargo.ClientCljHTTP,
		scalargo.ClientCoHTTP,
		// Keep: scalargo.ClientCurl,
		// Keep: scalargo.ClientFetch,
		scalargo.ClientGuzzle,
		scalargo.ClientHTTP1,
		scalargo.ClientHTTPClient,
		scalargo.ClientHTTP2,
		scalargo.ClientHttpie,
		scalargo.ClientHttr,
		scalargo.ClientJquery,
		scalargo.ClientLibCurl,
		scalargo.ClientNative,
		scalargo.ClientNetHTTP,
		scalargo.ClientNsurlSession,
		scalargo.ClientOkHTTP,
		scalargo.ClientPython3,
		scalargo.ClientRequest,
		scalargo.ClientRequests,
		scalargo.ClientRestMethod,
		scalargo.ClientRestSharp,
		scalargo.ClientUndici,
		scalargo.ClientUnirest,
		scalargo.ClientWebRequest,
		scalargo.ClientWget,
		scalargo.ClientXhr,
		scalargo.ClientHTTP,
		scalargo.ClientOfetch,
		scalargo.ClientHTTPXSync,
		scalargo.ClientHTTPXAsync,
		scalargo.ClientReqWest,
		scalargo.ClientOkhttp,
	}

	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithHiddenClients(hiddenClients...),
	)
}

// ExampleShowAllClients demonstrates showing all available client examples (default behavior)
func ExampleShowAllClients() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		// No WithHiddenClients or WithHideAllClients means all clients are visible
	)
}

// ExampleAllOptions demonstrates combining multiple customization options including authentication
func ExampleAllOptions() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithLayout(scalargo.LayoutModern),
		scalargo.WithDarkMode(),
		scalargo.WithHideDownloadButton(),
		scalargo.WithHiddenClients("fetch", "curl"),
		scalargo.WithOverrideCSS(`
			.scalar-api-reference {
				--scalar-color-1: #2d3748;
				--scalar-color-2: #4a5568;
			}
		`),
		// Add authentication configuration
		scalargo.WithAuthenticationOpts(
			scalargo.WithAPIKey("demo-api-key"),
			scalargo.WithHTTPBearerToken("demo-bearer-token"),
		),
	)
}
