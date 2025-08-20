# Scalar-Go Examples

This directory contains practical, runnable examples demonstrating how to use the [scalar-go](https://github.com/bdpiprava/scalar-go) library to generate beautiful API documentation from OpenAPI specifications.

## 📚 Available Examples

### 1. Basic Usage (`basic/`)
**Use Case**: Simple API documentation generation from a single OpenAPI spec file.

```bash
cd basic && go run main.go
# Visit: http://localhost:8080
```

**What it demonstrates**:
- Loading OpenAPI spec from a directory
- Generating HTML documentation
- Serving through HTTP

### 2. Multi-File Specifications (`multi_file_spec/`)
**Use Case**: Loading segmented OpenAPI specs from organized directory structures.

```bash
cd multi_file_spec && go run main.go
# Visit: http://localhost:8081
```

**What it demonstrates**:
- Loading specs with separate schema, path, and response files
- Using `NewV2` API for explicit configuration
- Directory structure organization

### 3. Customization Options (`customization/`)
**Use Case**: Extensive UI customization including themes, layouts, and custom CSS.

```bash
cd customization && go run main.go
# Visit: http://localhost:8082/theme/moon (and other endpoints)
```

**What it demonstrates**:
- Available themes: `default`, `moon`, `purple`, `solarized`, etc.
- Layout options: `modern`, `classic`
- Visibility controls: sidebar, models, dark mode
- Custom CSS overrides

### 4. Spec Modification (`spec_modification/`)
**Use Case**: Dynamically modifying OpenAPI specifications before rendering.

```bash
cd spec_modification && go run main.go
# Visit: http://localhost:8083/basic (and other endpoints)
```

**What it demonstrates**:
- Runtime spec modifications using `WithSpecModifier`
- Dynamic server URL configuration
- Adding custom tags and information
- Working with documented paths

### 5. HTTP Server Integration (`http_server_integration/`)
**Use Case**: Professional HTTP server setup with multiple documentation endpoints.

```bash
cd http_server_integration && go run main.go
# Visit: http://localhost:8084/docs (and other endpoints)
```

**What it demonstrates**:
- Multiple API version documentation
- Health and status endpoints
- Proper HTTP server configuration
- Metadata customization

### 6. URL-Based Loading (`url_based_loading/`)
**Use Case**: Loading OpenAPI specifications from remote URLs.

```bash
cd url_based_loading && go run main.go
# Visit: http://localhost:8085/scalar-galaxy (and other endpoints)
```

**What it demonstrates**:
- Loading specs from CDNs or remote repositories
- External API documentation
- Custom styling for external specs
- Error handling for remote resources

## 🚀 Quick Start

1. **Clone the repository**:
   ```bash
   git clone https://github.com/bdpiprava/scalar-go.git
   cd scalar-go
   ```

2. **Choose an example** and run it:
   ```bash
   cd examples/basic
   go run main.go
   ```

3. **Open your browser** and visit the URL shown in the terminal.

## 📁 Project Structure

```
examples/
├── README.md                      # This file
├── basic/
│   └── main.go                    # Simple usage example
├── multi_file_spec/
│   └── main.go                    # Multi-file spec loading
├── customization/
│   └── main.go                    # UI customization options
├── spec_modification/
│   └── main.go                    # Dynamic spec modification
├── http_server_integration/
│   └── main.go                    # Professional server setup
└── url_based_loading/
    └── main.go                    # Remote URL spec loading
```

## 🔧 Key Features Demonstrated

### Themes Available
- `ThemeDefault` - Clean default theme
- `ThemeMoon` - Dark theme with blue accents  
- `ThemePurple` - Purple color scheme
- `ThemeSolarized` - Based on Solarized colors
- `ThemeBluePlanet`, `ThemeDeepSpace`, `ThemeSaturn`, `ThemeKepler`, `ThemeMars` - Space-themed variants

### Layout Options
- `LayoutModern` - Contemporary design (default)
- `LayoutClassic` - Traditional documentation layout

### Customization Options
- **Custom CSS**: Override default styles
- **Metadata**: Add title, description, and custom fields
- **Visibility**: Control sidebar, models, download button
- **Dark Mode**: Force dark mode or hide toggle
- **Search**: Custom hotkey configuration
- **Client Examples**: Hide specific code examples

### Spec Loading Methods
- **Directory**: Load from local file system
- **URL**: Load from remote endpoints
- **Multi-file**: Combine segmented specifications
- **Dynamic**: Modify specs at runtime

## 💡 Common Use Cases

1. **Internal API Documentation**: Use `basic/` or `multi_file_spec/` examples
2. **Public API Portal**: Use `customization/` with your branding
3. **Development Environment**: Use `http_server_integration/` with health checks
4. **Documentation Hub**: Use `url_based_loading/` to aggregate multiple APIs
5. **CI/CD Integration**: Use `spec_modification/` for environment-specific docs

## 🛠️ Prerequisites

- Go 1.19 or later
- OpenAPI specification files (YAML or JSON)
- Web browser for viewing documentation

## 📖 Further Reading

- [Scalar Documentation](https://github.com/scalar/scalar)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Main Library Documentation](../README.md)

## 🤝 Contributing

Found an issue or want to add more examples? Please contribute to the [main repository](https://github.com/bdpiprava/scalar-go).

---

**Note**: All examples use sample OpenAPI specifications included in the `data/` directory. You can replace these with your own API specifications.