package scalargo

// Theme as a type based on string for theme identification
type Theme string

const (
	ThemeDefault    Theme = "default"
	ThemeAlternate  Theme = "alternate"
	ThemeMoon       Theme = "moon"
	ThemePurple     Theme = "purple"
	ThemeSolarized  Theme = "solarized"
	ThemeBluePlanet Theme = "bluePlanet"
	ThemeDeepSpace  Theme = "deepSpace"
	ThemeSaturn     Theme = "saturn"
	ThemeKepler     Theme = "kepler"
	ThemeMars       Theme = "mars"
	ThemeNone       Theme = "none"
	ThemeNil        Theme = ""
)

// WithTheme sets the theme for the Scalar UI
func WithTheme(theme Theme) func(*Options) {
	return func(o *Options) {
		o.Theme = theme
	}
}
