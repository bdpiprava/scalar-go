package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

// ExampleThemeDefault demonstrates the default theme
// @example Default Theme
// @description This example shows the default theme with clean, modern styling for professional API documentation.
func ExampleThemeDefault() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeDefault),
	)
}

// ExampleThemeMoon demonstrates the moon theme
// @example Moon Theme
// @description This example shows the moon theme with dark styling and blue accents, perfect for modern dark-mode preferences.
func ExampleThemeMoon() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeMoon),
	)
}

// ExampleThemePurple demonstrates the purple theme
// @example Purple Theme
// @description This example shows the purple theme with vibrant purple color scheme for distinctive and creative API documentation.
func ExampleThemePurple() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemePurple),
	)
}

// ExampleThemeSolarized demonstrates the solarized theme
// @example Solarized Theme
// @description This example shows the solarized theme based on the popular Solarized color scheme, offering excellent readability and reduced eye strain.
func ExampleThemeSolarized() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithTheme(scalargo.ThemeSolarized),
	)
}

// Layout examples - self-contained spec functions

// ExampleLayoutModern demonstrates the modern layout
// @example Modern Layout
// @description This example shows the modern layout with contemporary design elements and enhanced user experience for API documentation.
func ExampleLayoutModern() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithLayout(scalargo.LayoutModern),
	)
}

// ExampleLayoutClassic demonstrates the classic layout
// @example Classic Layout
// @description This example shows the classic layout with traditional documentation design, familiar to users of conventional API documentation tools.
func ExampleLayoutClassic() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithLayout(scalargo.LayoutClassic),
	)
}

// Visibility examples - self-contained spec functions

// ExampleHideSidebar demonstrates hiding the sidebar for a cleaner look
// @example Hide Sidebar
// @description This example shows how to hide the sidebar to create a cleaner, more focused documentation layout with more space for content.
func ExampleHideSidebar() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithSidebarVisibility(false),
	)
}

// ExampleHideModels demonstrates hiding the models section to focus on endpoints
// @example Hide Models
// @description This example shows how to hide the models section to focus purely on API endpoints, useful for endpoint-centric documentation.
func ExampleHideModels() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithHideModels(),
	)
}

// ExampleDarkMode demonstrates enabling dark mode by default
// @example Dark Mode
// @description This example shows how to enable dark mode by default, providing a modern dark interface that's easier on the eyes.
func ExampleDarkMode() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithDarkMode(),
	)
}

// Advanced examples - self-contained spec functions

// ExampleCustomCSS demonstrates custom CSS styling for branded documentation
// @example Custom CSS
// @description This example shows how to apply custom CSS overrides to create branded documentation with custom colors, fonts, and styling elements.
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
// @example All Options Combined
// @description This example shows how to combine multiple customization options including theme, layout, UI controls, custom CSS, and client hiding for comprehensive documentation branding.
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
