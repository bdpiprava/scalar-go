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
		scalargo.WithTheme(scalargo.ThemeDefault),
	)
}

// ExampleAllOptions demonstrates combining multiple customization options
func ExampleAllOptions() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemePurple),
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
	)
}
