## Overview 🌐

The ScalarGo package serves as a provider for the [Scalar](https://github.com/scalar/scalar) project. It offers easy to
integrate functions for documenting APIs in HTML format, with a focus on JSON data handling and web presentation
customization. This includes functions to serialize options into JSON, manage HTML escaping, and dynamically handle
different types of specification content.

## Features 🚀

Supports reading API specification from a directory with multiple files, a single file.

### Reading from single file

```go
// When file is located in directory /example/docs/ and filename is api.yaml(default lookup name)
content, err := scalargo.New("/example/docs/")

// When file is located in directory /example/docs/ and filename is different from default lookup name e.g. petstore.yaml
content, err := scalargo.New(
"/example/docs/",
scalargo.WithBaseFileName("petstore.yaml"),
)
```

### Reading from segmented files

The package supports reading segmented API specification files over schemas and paths. The segmented files are combined
into a single specification file before generating the API reference.

Expected directory structure:

```text
/example/docs/
    ├── api.yaml            // main file
    ├── schemas/            // directory for schema files
    │   ├── pet.yaml
    │   ├── user.yaml
    │   └── order.yaml
    ├── paths/              // directory for path files
    │   ├── pet.yaml
    │   ├── user.yaml
    │   └── order.yaml
    ├── responses/          // directory for response files
    └── └── Error.yaml
```

```go
// When segmented files are located in directory /example/docs/ following the expected directory structure
content, err := scalargo.New("/example/docs/")
```

## Customization Options ⚙️

The package allows extensive customization of the generated API reference through the `Options`

supporting scalar built-in options:

- **Theme**:  Select theme for scalar UI from
  available [themes](https://github.com/scalar/scalar/blob/main/documentation/themes.md).
- **Layout**: Chose between modern and classic layout designs
- **ShowSidebar**: Show or hide the sidebar in the API reference.
- **HideModels**: Hide the models section in the API reference.
- **HideDownloadButton**: Hide the download button in the API reference.
- **DarkMode**: Default dark mode for the API reference.
- **SearchHotKey**: Set a hotkey for the search functionality.
- **MetaData**: Set metadata for the API reference.
- **HiddenClients**: Hide clients in the API reference.

and customer options for easy of documenting APIs:

- **OverrideCSS**: A custom CSS style to override the default scalar style.
- **BaseFileName**: The base file name if it is not `api.yml`.
- **CDN**: URL of the CDN to load additional scripts or styles.

## Usage 📚

```go
package main

import (
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
    "github.com/bdpiprava/scalar-go/model"
)

func main() {
	apiDir := "path/to/api/directory"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := scalargo.New(
			apiDir,
			scalargo.WithBaseFileName("api.yml"),
            scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
			  // Customise the spec here
              spec.Info.Title = "PetStore API"
              return spec
            }),
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(content))
	})
	http.ListenAndServe(":8090", nil)
}
```

## Credits 🙏

- [Scalar](https://github.com/scalar/scalar) - The project that inspired this package.
- [Go Scalar API Reference](https://github.com/MarceloPetrucio/go-scalar-api-reference) - The package that inspired this
  package.