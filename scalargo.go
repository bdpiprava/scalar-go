package scalargo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"github.com/bdpiprava/scalar-go/loader"
	"github.com/bdpiprava/scalar-go/model"
	"github.com/bdpiprava/scalar-go/sanitizer"
)

// defaultTitle when title is not specified this default is used
const defaultTitle = "API Reference"

// validateURL validates that a URL uses a safe scheme (http or https only)
// Returns an error if the URL is invalid or uses a dangerous scheme
func validateURL(rawURL, fieldName string) error {
	if strings.TrimSpace(rawURL) == "" {
		return nil // Empty URLs are allowed (will use defaults or be omitted)
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid %s: %w", fieldName, err)
	}

	// Check that scheme is http or https only
	// This prevents javascript:, data:, file:, vbscript:, and other dangerous schemes
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("invalid %s: scheme must be http or https, got %q", fieldName, parsedURL.Scheme)
	}

	return nil
}

// New generates the HTML for the Scalar UI
func New(apiFilesDir string, opts ...Option) (string, error) {
	return NewV2(append(opts, WithSpecDir(apiFilesDir))...)
}

// NewV2 generate the HTML for the Scalar UI
func NewV2(opts ...Option) (string, error) {
	options := buildOptions(opts...)

	// Validate CDN URL to prevent XSS via dangerous URL schemes
	if err := validateURL(options.CDN, "CDN"); err != nil {
		return "", err
	}

	specScript, err := options.GetSpecScript()
	if err != nil {
		return "", err
	}

	title := extractTitle(options.Configurations, defaultTitle)

	return renderHTML(
		title,
		options.OverrideCSS,
		specScript,
		options.CDN,
	), nil
}

// buildOptions build Options from applying OptionFn to defaults
func buildOptions(opts ...Option) *Options {
	options := &Options{
		Configurations: map[string]any{
			keyTheme:  ThemeDefault,
			keyLayout: LayoutModern,
			keyMetaData: MetaData{
				"title": "API Reference",
			},
		},

		CDN:          DefaultCDN,
		BaseFileName: "api.yaml",
		RenderMode:   RenderModeDataAttribute, // Default to data-attribute for backward compatibility
	}

	for _, opt := range opts {
		opt(options)
	}
	return options
}

var htmlTemplate = template.Must(template.New("scalar").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>{{.Title}}</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>{{.CSS}}</style>
  </head>
  <body>
    {{.SpecScript}}
    <script src="{{.CDN}}"></script>
  </body>
</html>`))

// renderHTML generates HTML from the provided options with proper escaping to prevent XSS
func renderHTML(title, cssOverride, specScript, cdn string) string {
	var buf bytes.Buffer

	// Sanitize CSS to remove HTML tags while preserving CSS content
	sanitizedCSS := sanitizer.CSS(cssOverride)

	// Execute template with proper type conversions for context-aware escaping
	err := htmlTemplate.Execute(&buf, map[string]interface{}{
		"Title":      title,                      // Auto-escaped for HTML context
		"CSS":        template.CSS(sanitizedCSS), // CSS-safe after sanitization
		"SpecScript": template.HTML(specScript),  // Already validated JSON, safe to render as HTML
		"CDN":        cdn,                        // Auto-escaped for attribute context
	})

	if err != nil {
		// Template execution should never fail with our static template
		// If it does, return a safe error page instead of panicking
		return fmt.Sprintf("<!DOCTYPE html><html><head><title>Error</title></head><body>Template error: %s</body></html>",
			template.HTMLEscapeString(err.Error()))
	}

	return buf.String()
}

// GetSpecScript prepares and returns the spec script, prioritizing SpecURL, then SpecDirectory, then SpecBytes
func (o *Options) GetSpecScript() (string, error) {
	configAsBytes, err := json.Marshal(o.Configurations)
	if err != nil {
		return "", err
	}
	configJSON := strings.ReplaceAll(string(configAsBytes), `"`, `&quot;`)

	if strings.TrimSpace(o.SpecURL) != "" {
		// Validate SpecURL to prevent XSS via dangerous URL schemes
		if err := validateURL(o.SpecURL, "SpecURL"); err != nil {
			return "", err
		}

		return fmt.Sprintf(
			`<script id="api-reference" data-url="%s" data-configuration="%s"></script>`,
			o.SpecURL,
			configJSON,
		), nil
	}

	var spec *model.Spec
	switch {
	case o.SpecDirectory != "":
		spec, err = loader.LoadFromDir(o.SpecDirectory, o.BaseFileName)
		if err != nil {
			return "", err
		}
	case o.SpecBytes != nil:
		spec, err = loader.LoadFromBytes(o.SpecBytes)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("one of SpecURL, SpecDirectory or SpecBytes must be configured")
	}

	if o.SpecModifier != nil {
		spec = o.SpecModifier(spec)
	}

	metadata, ok := o.Configurations[keyMetaData].(MetaData)
	if !ok {
		metadata = MetaData{}
		o.Configurations[keyMetaData] = metadata
	}

	if title, ok := metadata["title"]; !ok || title == defaultTitle {
		metadata["title"] = spec.Info.Title
	}

	content, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		`<script id="api-reference" type="application/json" data-configuration="%s">%s</script>`,
		configJSON,
		string(content),
	), nil
}

// extractTitle safely extracts the title from metadata with fallback to default
func extractTitle(configurations map[string]any, fallback string) string {
	if metadata, ok := configurations[keyMetaData].(MetaData); ok {
		if titleVal, exists := metadata["title"]; exists {
			return fmt.Sprintf("%v", titleVal)
		}
	}
	return fallback
}
