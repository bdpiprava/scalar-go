package scalargo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bdpiprava/scalar-go/loader"
)

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
