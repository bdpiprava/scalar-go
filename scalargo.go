package scalargo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bdpiprava/scalar-go/loader"
	"github.com/bdpiprava/scalar-go/model"
)

const DefaultCDN = "https://cdn.jsdelivr.net/npm/@scalar/api-reference"

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

// Layout represents different layout options
type Layout string

const (
	LayoutModern  Layout = "modern"
	LayoutClassic Layout = "classic"
)

// SpecModifier is a function that can be used to override the spec
type SpecModifier func(spec *model.Spec) *model.Spec

type Options struct {
	Theme              Theme        `json:"theme,omitempty"`
	Layout             Layout       `json:"layout,omitempty"`
	Proxy              string       `json:"proxy,omitempty"`
	IsEditable         bool         `json:"isEditable,omitempty"`
	ShowSidebar        bool         `json:"showSidebar,omitempty"`
	HideModels         bool         `json:"hideModels,omitempty"`
	HideDownloadButton bool         `json:"hideDownloadButton,omitempty"`
	DarkMode           bool         `json:"darkMode,omitempty"`
	SearchHotKey       string       `json:"searchHotKey,omitempty"`
	MetaData           string       `json:"metaData,omitempty"`
	HiddenClients      []string     `json:"hiddenClients,omitempty"`
	Authentication     string       `json:"authentication,omitempty"`
	PathRouting        string       `json:"pathRouting,omitempty"`
	BaseServerURL      string       `json:"baseServerUrl,omitempty"`
	WithDefaultFonts   bool         `json:"withDefaultFonts,omitempty"`
	OverrideCSS        string       `json:"-"`
	BaseFileName       string       `json:"-"`
	CDN                string       `json:"-"`
	OverrideHandler    SpecModifier `json:"-"`
}

type Option func(*Options)

// WithTheme sets the theme for the Scalar UI
func WithTheme(theme Theme) func(*Options) {
	return func(o *Options) {
		o.Theme = theme
	}
}

// WithLayout sets the layout for the Scalar UI
func WithLayout(layout Layout) func(*Options) {
	return func(o *Options) {
		o.Layout = layout
	}
}

// WithCDN sets the CDN for the Scalar UI
func WithCDN(cdn string) func(*Options) {
	return func(o *Options) {
		o.CDN = cdn
	}
}

// WithProxy sets the proxy for the Scalar UI
func WithProxy(proxy string) func(*Options) {
	return func(o *Options) {
		o.Proxy = proxy
	}
}

// WithEditable sets the editable state for the Scalar UI
func WithEditable() func(*Options) {
	return func(o *Options) {
		o.IsEditable = true
	}
}

// WithSidebarVisibility sets the sidebar visibility for the Scalar UI
func WithSidebarVisibility(visible bool) func(*Options) {
	return func(o *Options) {
		o.ShowSidebar = visible
	}
}

// WithHideModels sets the models visibility for the Scalar UI
func WithHideModels() func(*Options) {
	return func(o *Options) {
		o.HideModels = true
	}
}

// WithHideDownloadButton sets the download button visibility for the Scalar UI
func WithHideDownloadButton() func(*Options) {
	return func(o *Options) {
		o.HideDownloadButton = true
	}
}

// WithDarkMode sets the dark mode for the Scalar UI
func WithDarkMode() func(*Options) {
	return func(o *Options) {
		o.DarkMode = true
	}
}

// WithSearchHotKey sets the search hot key for the Scalar UI
func WithSearchHotKey(searchHotKey string) func(*Options) {
	return func(o *Options) {
		o.SearchHotKey = searchHotKey
	}
}

// WithMetaData sets the metadata for the Scalar UI
func WithMetaData(metaData string) func(*Options) {
	return func(o *Options) {
		o.MetaData = metaData
	}
}

// WithHiddenClients sets the hidden clients for the Scalar UI
func WithHiddenClients(hiddenClients []string) func(*Options) {
	return func(o *Options) {
		o.HiddenClients = hiddenClients
	}
}

// WithOverrideCSS sets the override CSS for the Scalar UI
func WithOverrideCSS(overrideCSS string) func(*Options) {
	return func(o *Options) {
		o.OverrideCSS = overrideCSS
	}
}

// WithAuthentication sets the authentication method for the Scalar UI
func WithAuthentication(authentication string) func(*Options) {
	return func(o *Options) {
		o.Authentication = authentication
	}
}

// WithPathRouting sets the path routing for the Scalar UI
func WithPathRouting(pathRouting string) func(*Options) {
	return func(o *Options) {
		o.PathRouting = pathRouting
	}
}

// WithBaseServerURL sets the base server URL for the Scalar UI
func WithBaseServerURL(baseServerURL string) func(*Options) {
	return func(o *Options) {
		o.BaseServerURL = baseServerURL
	}
}

// WithDefaultFonts sets the default fonts usage for the Scalar UI
func WithDefaultFonts() func(*Options) {
	return func(o *Options) {
		o.WithDefaultFonts = true
	}
}

// WithBaseFileName sets the base file name for the Scalar UI
func WithBaseFileName(baseFileName string) func(*Options) {
	return func(o *Options) {
		o.BaseFileName = baseFileName
	}
}

// WithSpecModifier allows to modify the spec before rendering
func WithSpecModifier(handler SpecModifier) func(*Options) {
	return func(o *Options) {
		o.OverrideHandler = handler
	}
}

// New generates the HTML for the Scalar UI
func New(apiFilesDir string, opts ...Option) (string, error) {
	options := &Options{
		CDN:          DefaultCDN,
		Layout:       LayoutModern,
		Theme:        ThemeDefault,
		BaseFileName: "api.yaml",
	}

	for _, opt := range opts {
		opt(options)
	}

	spec, err := loader.LoadWithName(apiFilesDir, options.BaseFileName)
	if err != nil {
		return "", err
	}

	if options.OverrideHandler != nil {
		spec = options.OverrideHandler(spec)
	}

	content, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}

	dataConfig, err := json.Marshal(options)
	if err != nil {
		return "", err
	}
	escapedJSON := strings.ReplaceAll(string(dataConfig), `"`, `&quot;`)

	return fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
      <head>
        <title>%s</title>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <style>%s</style>
      </head>
      <body>
        <script id="api-reference" type="application/json" data-configuration="%s">%s</script>
        <script src="%s"></script>
      </body>
    </html>
  `, spec.Info.Title, options.OverrideCSS, escapedJSON, string(content), options.CDN), nil
}
