## Overview 🌐

The ScalarGo package serves as a provider for the [Scalar](https://github.com/scalar/scalar) project. It offers easy to
integrate functions for documenting APIs in HTML format, with a focus on JSON data handling and web presentation
customization. This includes functions to serialize options into JSON, manage HTML escaping, and dynamically handle
different types of specification content.

## Features 🚀

Supports reading API specification from multiple sources with the following precedence order:
1. **URL** (highest priority) - Loads spec from a remote URL
2. **Directory** - Loads from a directory with multiple files or single file
3. **Bytes** (lowest priority) - Loads from in-memory bytes (YAML or JSON)

> **⚠️ Deprecation Notice**: The `Load` and `LoadWithName` functions in the `loader` package are deprecated. Use `NewV2` with the appropriate `With*` options instead.

### Reading from URL

```go
// Load spec from a remote URL
content, err := scalargo.NewV2(
    scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
)
```

### Reading from single file

See [Documentation](https://bdpiprava.github.io/scalar-go) for more details.

```go
// When file is located in directory /example/docs/ and filename is api.yaml(default lookup name)
content, err := scalargo.New("/example/docs/")

// Using NewV2 (recommended approach)
content, err := scalargo.NewV2(
    scalargo.WithSpecDir("/example/docs/"),
)

// When file is located in directory /example/docs/ and filename is different from default lookup name e.g. petstore.yaml
content, err := scalargo.NewV2(
    scalargo.WithSpecDir("/example/docs/"),
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
content, err := scalargo.NewV2(
    scalargo.WithSpecDir("/example/docs/"),
)
```

### Reading from bytes (self-contained builds)

`WithSpecBytes` enables **self-contained builds** by embedding the API specification directly in your binary. This is particularly useful for:
- Deploying applications without external file dependencies
- Creating portable executables
- Ensuring the API spec is always available

```go
//go:embed api.yaml
var specBytes []byte

// Load spec from embedded bytes (YAML or JSON format supported)
content, err := scalargo.NewV2(
    scalargo.WithSpecBytes(specBytes),
)

// You can also load from any byte slice
yamlSpec := []byte(`
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
paths: {}
`)

content, err := scalargo.NewV2(
    scalargo.WithSpecBytes(yamlSpec),
)
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

### Basic Example

```go
package main

import (
    "net/http"

    scalargo "github.com/bdpiprava/scalar-go"
    "github.com/bdpiprava/scalar-go/model"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Using NewV2 (recommended approach)
        content, err := scalargo.NewV2(
            scalargo.WithSpecDir("path/to/api/directory"),
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

### Self-Contained Build Example

```go
package main

import (
    _ "embed"
    "net/http"

    scalargo "github.com/bdpiprava/scalar-go"
)

//go:embed api.yaml
var apiSpec []byte

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Self-contained build using embedded bytes
        content, err := scalargo.NewV2(
            scalargo.WithSpecBytes(apiSpec),
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

### Precedence Order

When multiple spec sources are provided, ScalarGo follows this precedence order:

1. **`WithSpecURL`** (highest priority) - Remote URL takes precedence
2. **`WithSpecDir`** - Directory-based specs are used if no URL is provided  
3. **`WithSpecBytes`** (lowest priority) - Byte-based specs are used as fallback

```go
// This will use the URL, ignoring the directory and bytes
content, err := scalargo.NewV2(
    scalargo.WithSpecURL("https://example.com/api.yaml"),  // <- This wins
    scalargo.WithSpecDir("/path/to/api/"),                 // <- Ignored
    scalargo.WithSpecBytes(embeddedBytes),                 // <- Ignored
)
```

See the [examples](./main/main.go) for more details.

## Credits 🙏

- [Scalar](https://github.com/scalar/scalar) - The project that inspired this package.
- [Go Scalar API Reference](https://github.com/MarceloPetrucio/go-scalar-api-reference) - The package that inspired this package.
