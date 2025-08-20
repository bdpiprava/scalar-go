package scalargo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bdpiprava/scalar-go/loader"
	"github.com/bdpiprava/scalar-go/model"
)

// defaultTitle when title is not specified this default is used
const defaultTitle = "API Reference"

// New generates the HTML for the Scalar UI
func New(apiFilesDir string, opts ...Option) (string, error) {
	return NewV2(append(opts, WithSpecDir(apiFilesDir))...)
}

// NewV2 generate the HTML for the Scalar UI
func NewV2(opts ...Option) (string, error) {
	options := buildOptions(opts...)
	specScript, err := options.GetSpecScript()
	if err != nil {
		return "", err
	}

	return renderHTML(
		fmt.Sprintf("%v", options.Configurations[keyMetaData].(MetaData)["title"]),
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
	}

	for _, opt := range opts {
		opt(options)
	}
	return options
}

// renderHTML generte html from the provided options
func renderHTML(title, ccsOverride, specScript, cdn string) string {
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
        %s
        <script src="%s"></script>
      </body>
    </html>
  `, title, ccsOverride, specScript, cdn)
}

// GetSpecScript prepares and returns the spec script, prioritizing SpecURL, then SpecDirectory, then SpecBytes
func (o *Options) GetSpecScript() (string, error) {
	configAsBytes, err := json.Marshal(o.Configurations)
	if err != nil {
		return "", err
	}
	configJSON := strings.ReplaceAll(string(configAsBytes), `"`, `&quot;`)

	if strings.TrimSpace(o.SpecURL) != "" {
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

	metadata := o.Configurations[keyMetaData].(MetaData)
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
