package scalargo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bdpiprava/scalar-go/model"
)

// DefaultCDN default CDN for api-reference
const DefaultCDN = "https://cdn.jsdelivr.net/npm/@scalar/api-reference"

const (
	keyTheme              = "theme"
	keyLayout             = "layout"
	keyProxy              = "proxy"
	keyIsEditable         = "isEditable"
	keyShowSidebar        = "showSidebar"
	keyHideModels         = "hideModels"
	keyHideDownloadButton = "hideDownloadButton"
	keyDarkMode           = "darkMode"
	keyFroceDarkMode      = "forceDarkModeState"
	keyHideDarkModeToggle = "hideDarkModeToggle"
	keySearchHotKey       = "searchHotKey"
	keyHiddenClients      = "hiddenClients"
	keyAuthentication     = "authentication"
	keyPathRouting        = "pathRouting"
	keyBaseServerURL      = "baseServerUrl"
	keyWithDefaultFonts   = "withDefaultFonts"
	keyServers            = "servers"
	keyMetaData           = "metadata"
)

// SpecModifier is a function that can be used to override the spec
type SpecModifier func(spec *model.Spec) *model.Spec

type Options struct {
	Configurations map[string]any
	OverrideCSS    string
	BaseFileName   string
	CDN            string
	SpecModifier   SpecModifier
	SpecDirectory  string
	SpecURL        string
	SpecBytes      []byte
}

type Option func(*Options)

// WithCDN sets the CDN for the Scalar UI
func WithCDN(cdn string) func(*Options) {
	return func(o *Options) {
		o.CDN = cdn
	}
}

// WithProxy sets the proxy for the Scalar UI
func WithProxy(proxy string) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyProxy] = proxy
	}
}

// WithEditable sets the editable state for the Scalar UI
func WithEditable() func(*Options) {
	return func(o *Options) {
		o.Configurations[keyIsEditable] = true
	}
}

// WithSidebarVisibility sets the sidebar visibility for the Scalar UI
func WithSidebarVisibility(visible bool) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyShowSidebar] = visible
	}
}

// WithHideModels sets the models visibility for the Scalar UI
func WithHideModels() func(*Options) {
	return func(o *Options) {
		o.Configurations[keyHideModels] = true
	}
}

// WithHideDownloadButton hide to download OpenAPI spec button
func WithHideDownloadButton() func(*Options) {
	return func(o *Options) {
		o.Configurations[keyHideDownloadButton] = true
	}
}

// WithDarkMode set the dark mode as default
func WithDarkMode() func(*Options) {
	return func(o *Options) {
		o.Configurations[keyDarkMode] = true
	}
}

// WithForceDarkMode makes it always this state no matter what
func WithForceDarkMode() func(*Options) {
	return func(o *Options) {
		o.Configurations[keyFroceDarkMode] = true
	}
}

// WithHideDarkModeToggle hides the dark mode toggle button
func WithHideDarkModeToggle() func(*Options) {
	return func(o *Options) {
		o.Configurations[keyHideDarkModeToggle] = true
	}
}

// WithSearchHotKey sets the search hot key for the Scalar UI
func WithSearchHotKey(searchHotKey string) func(*Options) {
	return func(o *Options) {
		o.Configurations[keySearchHotKey] = searchHotKey
	}
}

// WithHiddenClients hide the set clients
func WithHiddenClients(hiddenClients ...string) func(*Options) {
	return func(o *Options) {
		value := o.Configurations[keyHiddenClients]

		// WithHideAllClients() takes precedence over this
		if strings.ToLower(fmt.Sprintf("%v", value)) != "true" {
			o.Configurations[keyHiddenClients] = hiddenClients
		}
	}
}

// WithHideAllClients sets the hidden clients for the Scalar UI
func WithHideAllClients() func(*Options) {
	return func(o *Options) {
		o.Configurations[keyHiddenClients] = true
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
		o.Configurations[keyAuthentication] = authentication
	}
}

// WithPathRouting sets the path routing for the Scalar UI
func WithPathRouting(pathRouting string) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyPathRouting] = pathRouting
	}
}

// WithBaseServerURL sets the base server URL for the Scalar UI
func WithBaseServerURL(baseServerURL string) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyBaseServerURL] = baseServerURL
	}
}

// WithDefaultFonts sets the default fonts usage for the Scalar UI
func WithDefaultFonts() func(*Options) {
	return func(o *Options) {
		o.Configurations[keyWithDefaultFonts] = true
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
		o.SpecModifier = handler
	}
}

// WithSpecDir read spec from directory
func WithSpecDir(specDir string) func(*Options) {
	return func(o *Options) {
		o.SpecDirectory = specDir
	}
}

// WithSpecURL set the spec URL in the doc
func WithSpecURL(specURL string) func(*Options) {
	return func(o *Options) {
		o.SpecURL = specURL
	}
}

// WithSpecBytes loads the spec from the provided bytes in either YAML or JSON format
func WithSpecBytes(specBytes []byte) func(*Options) {
	return func(o *Options) {
		o.SpecBytes = specBytes
	}
}

// WithAuthenticationOpts sets the authentication method for the Scalar UI
func WithAuthenticationOpts(opts ...AuthOption) func(*Options) {
	auth := make(AuthenticationOption)
	for _, opt := range opts {
		opt(auth)
	}

	return func(o *Options) {
		content, err := json.Marshal(auth)
		if err == nil {
			o.Configurations[keyAuthentication] = string(content)
		}
	}
}
