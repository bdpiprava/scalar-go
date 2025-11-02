# 🚀 Scalar-Go

[![Go Reference](https://pkg.go.dev/badge/github.com/bdpiprava/scalar-go.svg)](https://pkg.go.dev/github.com/bdpiprava/scalar-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/bdpiprava/scalar-go)](https://goreportcard.com/report/github.com/bdpiprava/scalar-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> 🎯 **Transform your OpenAPI specs into beautiful, interactive documentation with just a few lines of Go code!**

Scalar-Go is the official Go integration for the powerful [Scalar API Documentation](https://github.com/scalar/scalar)
platform. Whether you're building internal tools, public APIs, or microservices, Scalar-Go makes it incredibly easy to
generate stunning, interactive API documentation that your developers will actually want to use.

## ✨ Why Choose Scalar-Go?

- **🎨 Beautiful by Default**: Professional-looking docs with multiple themes and layouts
- **⚡ Lightning Fast**: Generate documentation in milliseconds, not minutes
- **🔧 Incredibly Flexible**: Support for files, URLs, embedded specs, and runtime modifications
- **🌍 Universal**: Works with any OpenAPI 3.x specification
- **📱 Mobile-First**: Responsive design that looks great on all devices
- **🎭 Highly Customizable**: Custom CSS, themes, and UI options
- **📊 Production Ready**: Used by teams worldwide for mission-critical documentation

## 🛠️ Installation

```bash
go get github.com/bdpiprava/scalar-go
```

## 🎯 Quick Start

Get your API documentation up and running in under 30 seconds:

```go
package main

import (
	"fmt"
	"net/http"
	scalargo "github.com/bdpiprava/scalar-go"
)

func main() {
	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		// Generate beautiful docs from your OpenAPI spec
		html, err := scalargo.NewV2(
			scalargo.WithSpecDir("./api"), // or WithSpecURL, WithSpecBytes
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprint(w, html)
	})

	fmt.Println("📚 API Docs available at: http://localhost:8080/docs")
	http.ListenAndServe(":8080", nil)
}
```

> **📚 Interactive Demo**: Visit [https://bdpiprava.github.io/scalar-go](https://bdpiprava.github.io/scalar-go) to see
> all examples in action!

## 🔥 Core Features

### 📁 **Multiple Spec Sources** (Priority Order)

Load your OpenAPI specifications from anywhere:

1. **🌐 Remote URLs** (highest priority) - Perfect for CI/CD and external specs
2. **📂 Local Directories** - Great for development and file-based workflows
3. **💾 Embedded Bytes** (lowest priority) - Ideal for self-contained deployments

> **💡 Pro Tip**: Mix and match sources! Scalar-Go automatically picks the best available source.

### 🌐 **Remote URL Loading**

Perfect for loading specs from GitHub, CDNs, or your API servers:

```go
// Load from GitHub, CDN, or any public URL
html, err := scalargo.NewV2(
    scalargo.WithSpecURL("https://petstore3.swagger.io/api/v3/openapi.json"),
    scalargo.WithMetaDataOpts(
       scalargo.WithTitle("🐾 Pet Store API"),
        scalargo.WithKeyValue("description", "The most comprehensive pet store API"),
    ),
)
	

// Load your company's API spec from private repos
html, err := scalargo.NewV2(
    scalargo.WithSpecURL("https://api.yourcompany.com/openapi.yaml"),
    scalargo.WithTheme(scalargo.ThemeMoon), // Dark theme
)
```

### 📂 **Directory-Based Loading**

Great for local development and organized spec files:

```go
// Load from directory with default filename (api.yaml)
html, err := scalargo.NewV2(
    scalargo.WithSpecDir("./docs/api"),
)

// Specify custom filename
html, err := scalargo.NewV2(
    scalargo.WithSpecDir("./specs"),
    scalargo.WithBaseFileName("petstore.yaml"), // or .json
)

// Legacy support (still works, but use NewV2 for new projects)
html, err := scalargo.New("/path/to/specs/") // ⚠️ Consider migrating to NewV2
```

### 🗂️ **Multi-File Specifications**

Perfect for large APIs with organized file structures. Scalar-Go automatically combines segmented files:

```text
📁 /api-specs/
├── 📄 api.yaml           # Main specification file
├── 📁 schemas/           # Data models and schemas
│   ├── 📄 User.yaml
│   ├── 📄 Pet.yaml
│   └── 📄 Order.yaml
├── 📁 paths/             # API endpoints
│   ├── 📄 users.yaml
│   ├── 📄 pets.yaml
│   └── 📄 orders.yaml
└── 📁 responses/         # Reusable responses
    └── 📄 Error.yaml
```

```go
// Scalar-Go intelligently combines all files
html, err := scalargo.NewV2(
    scalargo.WithSpecDir("./api-specs"),
    scalargo.WithTheme(scalargo.ThemeDefault),
)
// ✨ Automatically merges schemas/, paths/, and responses/ into main spec
```

### 💾 **Embedded Specifications**

Build self-contained applications with embedded specs - perfect for containers and serverless:

```go
package main

import (
	_ "embed" // Enable embed functionality
	scalargo "github.com/bdpiprava/scalar-go"
)

//go:embed openapi.yaml
var apiSpec []byte

//go:embed company-logo.css
var customCSS string

func generateDocs() (string, error) {
	return scalargo.NewV2(
		// 🚀 Zero external dependencies!
		scalargo.WithSpecBytes(apiSpec),
		scalargo.WithOverrideCSS(customCSS),
		scalargo.WithTheme(scalargo.ThemePurple),
	)
}

// Or create specs programmatically
func dynamicSpec() (string, error) {
	spec := []byte(`
openapi: 3.0.0
info:
  title: "🚀 Dynamic API"
  version: "1.0.0"
  description: "Generated at runtime!"
paths:
  /health:
    get:
      summary: Health Check
      responses:
        '200':
          description: OK
`)

	return scalargo.NewV2(
		scalargo.WithSpecBytes(spec),
		scalargo.WithDarkMode(), // 🌙 Dark mode by default
	)
}
```

> **🎯 Use Cases**: Docker containers, AWS Lambda, single-binary deployments, offline documentation

## 🎨 Customization Showcase

Make your documentation uniquely yours with extensive customization options:

### 🌈 **Stunning Themes**

Choose from professionally designed themes:

```go
// 🌟 Available Themes
scalargo.ThemeDefault // Clean, professional
scalargo.ThemeMoon    // Dark with blue accents  
scalargo.ThemePurple     // Vibrant purple vibes
scalargo.ThemeSolarized  // Easy on the eyes
scalargo.ThemeBluePlanet // Space-age blue
scalargo.ThemeDeepSpace // Deep cosmic theme
scalargo.ThemeSaturn    // Ringed planet aesthetics
scalargo.ThemeKepler     // Exoplanet explorer
scalargo.ThemeMars       // Red planet inspired

// Apply any theme
html, err := scalargo.NewV2(
    scalargo.WithSpecURL("https://api.example.com/openapi.json"),
    scalargo.WithTheme(scalargo.ThemeMoon), // 🌙
)
```

### 📐 **Layout Options**

```go
// Modern (default) - Contemporary, spacious design
scalargo.WithLayout(scalargo.LayoutModern)

// Classic - Traditional documentation feel
scalargo.WithLayout(scalargo.LayoutClassic)
```

### 🎛️ **UI Controls**

```go
html, err := scalargo.NewV2(
scalargo.WithSpecDir("./api"),

// Visibility Controls
scalargo.WithSidebarVisibility(false), // Hide sidebar for focus
scalargo.WithHideModels(),             // Hide schema models
scalargo.WithHideDownloadButton(), // Remove download option

// Dark Mode Options
scalargo.WithDarkMode(), // Default to dark mode
scalargo.WithForceDarkMode(), // Lock to dark mode
scalargo.WithHideDarkModeToggle(), // Remove theme switcher

// Advanced Options
scalargo.WithSearchHotKey("ctrl+k"), // Custom search shortcut
scalargo.WithHiddenClients("curl", "php"), // Hide specific code examples
scalargo.WithHideAllClients(), // Hide all code examples
)
```

### 🎨 **Custom Styling**

```go
customCSS := `
/* Your brand colors and fonts */
:root {
  --scalar-color-1: #ff6b6b;
  --scalar-color-2: #4ecdc4;
  --scalar-font: 'Inter', sans-serif;
}

/* Custom component styling */
.section-header {
  background: linear-gradient(45deg, #ff6b6b, #4ecdc4);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
`

html, err := scalargo.NewV2(
    scalargo.WithSpecDir("./api"),
    scalargo.WithOverrideCSS(customCSS),
    scalargo.WithTheme(scalargo.ThemeDefault),
)
```

### 📊 **Metadata & Branding**

```go
html, err := scalargo.NewV2(
    scalargo.WithSpecURL("https://api.company.com/openapi.yaml"),
    scalargo.WithMetaDataOpts(
        scalargo.WithTitle("🚀 CompanyName API Hub"),
        scalargo.WithKeyValue("description", "The definitive API reference"),
        scalargo.WithKeyValue("logo", "https://company.com/logo.png"),
    ),
)
```

### 🔐 **Authentication Configuration**

Scalar-Go provides comprehensive authentication support for modern APIs, including API Keys, HTTP Basic/Bearer, and OAuth2 flows.

#### **Enhanced API Key Authentication**

```go
// Simple API Key (header-based, backward compatible)
scalargo.WithAuthenticationOpts(
    scalargo.WithAPIKey("your-api-key"),
)

// Custom header name
scalargo.WithAuthenticationOpts(
    scalargo.WithAPIKey("your-api-key", scalargo.WithAPIKeyName("X-API-Key")),
)

// Query parameter-based API Key
scalargo.WithAuthenticationOpts(
    scalargo.WithAPIKeyQuery("api_key", "your-api-key"),
)

// Cookie-based API Key
scalargo.WithAuthenticationOpts(
    scalargo.WithAPIKeyCookie("session_token", "your-token"),
)
```

#### **HTTP Basic & Bearer Authentication**

```go
// HTTP Basic Auth
scalargo.WithAuthenticationOpts(
    scalargo.WithHTTPBasicAuth("username", "password"),
)

// HTTP Bearer Token
scalargo.WithAuthenticationOpts(
    scalargo.WithHTTPBearerToken("your-bearer-token"),
)
```

#### **OAuth2 Authentication**

Full OAuth2 support with all standard flows and PKCE.

**Authorization Code Flow (Recommended)**
```go
scalargo.WithAuthenticationOpts(
    scalargo.WithOAuth2AuthorizationCode(
        "https://auth.example.com/oauth/authorize",
        "https://auth.example.com/oauth/token",
        scalargo.WithOAuth2ClientID("my-client-id"),
        scalargo.WithOAuth2RedirectURI("https://myapp.com/callback"),
        scalargo.WithOAuth2PKCE(scalargo.PKCES256), // SHA-256 PKCE (recommended)
        scalargo.WithOAuth2Scopes("read:api", "write:api"),
    ),
)
```

**Client Credentials Flow**
```go
scalargo.WithAuthenticationOpts(
    scalargo.WithOAuth2ClientCredentials(
        "https://auth.example.com/oauth/token",
        scalargo.WithOAuth2ClientID("service-account"),
        scalargo.WithOAuth2ClientSecret("super-secret"),
    ),
)
```

**Advanced OAuth2 Customization**
```go
scalargo.WithAuthenticationOpts(
    scalargo.WithOAuth2AuthorizationCode(
        "https://auth.example.com/oauth/authorize",
        "https://auth.example.com/oauth/token",
        scalargo.WithOAuth2ClientID("my-client"),
        scalargo.WithOAuth2CustomToken("custom_access_token"), // Custom token field name
        scalargo.WithOAuth2AdditionalAuthParams(map[string]string{
            "audience": "https://api.example.com",
        }),
        scalargo.WithOAuth2AdditionalTokenParams(map[string]string{
            "resource": "https://resource.example.com",
        }),
        scalargo.WithOAuth2CredentialsLocation(scalargo.OAuth2CredentialsHeader),
    ),
)
```

#### **Multiple Security Schemes**

Configure multiple authentication methods for your API:

```go
scalargo.WithAuthenticationOpts(
    // Define multiple security schemes
    scalargo.WithSecurityScheme("api_key",
        scalargo.APIKeyScheme("X-API-Key", scalargo.APIKeyLocationHeader, "default-key"),
    ),
    scalargo.WithSecurityScheme("bearer_auth",
        scalargo.BearerScheme("default-token"),
    ),
    scalargo.WithSecurityScheme("oauth2",
        scalargo.OAuth2Scheme(
            scalargo.OAuth2FlowAuthorizationCode,
            scalargo.OAuth2Config{
                AuthorizationURL: "https://auth.example.com/authorize",
                TokenURL:         "https://auth.example.com/token",
                ClientID:         "my-client",
                UsePKCE:          scalargo.PKCES256,
                SelectedScopes:   []string{"read:api", "write:api"},
            },
        ),
    ),
    // Set preferred security scheme
    scalargo.WithPreferredSecurityScheme("bearer_auth"),
)
```

**PKCE Modes:**
- `scalargo.PKCES256` - SHA-256 PKCE (recommended for production)
- `scalargo.PKCEPlain` - Plain PKCE
- `scalargo.PKCEDisabled` - Disable PKCE

**OAuth2 Credentials Location:**
- `scalargo.OAuth2CredentialsHeader` - Send credentials in Authorization header (default)
- `scalargo.OAuth2CredentialsBody` - Send credentials in request body

## 🚀 Real-World Examples

### 🏢 **Enterprise API Documentation**

```go
package main

import (
	"fmt"
	"net/http"
	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/model"
)

func main() {
	// Multiple API versions with custom branding
	http.HandleFunc("/docs/v1", generateV1Docs)
	http.HandleFunc("/docs/v2", generateV2Docs)
	http.HandleFunc("/docs/internal", generateInternalDocs)

	fmt.Println("🏢 Enterprise API Hub:")
	fmt.Println("   📚 Public API v1:  http://localhost:8080/docs/v1")
	fmt.Println("   🚀 Public API v2:  http://localhost:8080/docs/v2")
	fmt.Println("   🔒 Internal APIs:  http://localhost:8080/docs/internal")

	http.ListenAndServe(":8080", nil)
}

func generateV1Docs(w http.ResponseWriter, r *http.Request) {
	html, err := scalargo.NewV2(
		scalargo.WithSpecDir("./specs/v1"),
		scalargo.WithTheme(scalargo.ThemeDefault),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("🏢 Company API v1.0"),
			scalargo.WithKeyValue("description", "Stable production API"),
		),
		// Add environment-specific server URLs
		scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
			spec.Servers = []model.Server{
				{URL: "https://api.company.com/v1", Description: "Production"},
				{URL: "https://staging-api.company.com/v1", Description: "Staging"},
			}
			return spec
		}),
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("Documentation error: %v", err), 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func generateV2Docs(w http.ResponseWriter, r *http.Request) {
	html, err := scalargo.NewV2(
		scalargo.WithSpecURL("https://github.com/company/api-specs/raw/main/v2.yaml"),
		scalargo.WithTheme(scalargo.ThemeMoon), // Modern dark theme
		scalargo.WithDarkMode(),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("🚀 Company API v2.0 (Beta)"),
			scalargo.WithKeyValue("description", "Next-generation API with GraphQL support"),
		),
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("Documentation error: %v", err), 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func generateInternalDocs(w http.ResponseWriter, r *http.Request) {
	// Internal documentation with restricted styling
	html, err := scalargo.NewV2(
		scalargo.WithSpecDir("./internal-specs"),
		scalargo.WithTheme(scalargo.ThemeSolarized),
		scalargo.WithHideDownloadButton(),                // Prevent spec downloads
		scalargo.WithHiddenClients("curl", "javascript"), // Hide external clients
		scalargo.WithOverrideCSS(`
            .section-header::before {
                content: "🔒 INTERNAL ";
                color: #ff6b6b;
                font-weight: bold;
            }
        `),
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("Documentation error: %v", err), 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}
```

### 🐳 **Containerized Microservice Documentation**

```go
package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/model"
)

//go:embed openapi.yaml
var apiSpec []byte

//go:embed assets/custom.css
var brandingCSS string

func main() {
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"status": "healthy", "service": "api-docs"}`)
	})

	// Self-contained documentation
	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		html, err := scalargo.NewV2(
			scalargo.WithSpecBytes(apiSpec), // 🚀 No external files needed!
			scalargo.WithTheme(scalargo.ThemeBluePlanet),
			scalargo.WithOverrideCSS(brandingCSS),

			// Dynamic environment-based configuration
			scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
				env := os.Getenv("ENVIRONMENT")
				if env == "" {
					env = "development"
				}

				// Dynamic title with environment
				spec.Info.Title = fmt.Sprintf("%s (%s)", spec.Info.Title, env)

				// Environment-specific servers
				switch env {
				case "production":
					spec.Servers = []model.Server{
						{URL: "https://api.company.com", Description: "Production API"},
					}
				case "staging":
					spec.Servers = []model.Server{
						{URL: "https://staging.company.com", Description: "Staging API"},
					}
				default:
					spec.Servers = []model.Server{
						{URL: "http://localhost:8080", Description: "Development API"},
					}
				}

				return spec
			}),
		)

		if err != nil {
			http.Error(w, fmt.Sprintf("Documentation generation failed: %v", err), 500)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
		fmt.Fprint(w, html)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🐳 Containerized API Documentation\n")
	fmt.Printf("   📚 Documentation: http://localhost:%s/docs\n", port)
	fmt.Printf("   💚 Health Check:  http://localhost:%s/health\n", port)

	http.ListenAndServe(":"+port, nil)
}
```

### 🔄 **Dynamic Runtime Modification**

```go
func advancedSpecModification() (string, error) {
    return scalargo.NewV2(
        scalargo.WithSpecDir("./api"),
        scalargo.WithSpecModifier(func (spec *model.Spec) *model.Spec {
            // Add build information
            buildTime := time.Now().Format("2006-01-02 15:04:05")
            spec.Info.Description = fmt.Sprintf("%s\n\n**Last Updated:** %s", *spec.Info.Description, buildTime)
            
            // Add API statistics
            paths := spec.DocumentedPaths()
            spec.Info.Description = fmt.Sprintf("%s\n**Total Endpoints:** %d", *spec.Info.Description, len(paths))
            
            // Add custom tags
            spec.Tags = append(spec.Tags, model.Tag{
                Name:        "build-info",
                Description: "Automatically generated build information",
            })
            
            return spec
        }),
        scalargo.WithTheme(scalargo.ThemePurple),
    )
}
```

## 🎯 Specification Source Priority

Scalar-Go intelligently handles multiple spec sources with a clear priority system:

```go
// 🎯 Priority Demonstration
html, err := scalargo.NewV2(
    // 🥇 1st Priority: Remote URL (if provided)
    scalargo.WithSpecURL("https://api.example.com/openapi.yaml"),
    
    // 🥈 2nd Priority: Local Directory (fallback if URL fails)
    scalargo.WithSpecDir("./backup-specs"),
    
    // 🥉 3rd Priority: Embedded Bytes (ultimate fallback)
    scalargo.WithSpecBytes(embeddedSpec),
    
    // ✨ These always apply regardless of source
    scalargo.WithTheme(scalargo.ThemeMoon),
    scalargo.WithDarkMode(),
)
```

**🧠 Smart Behavior:**

- **URL Available?** → Load from URL, ignore directory and bytes
- **URL Failed?** → Try directory, ignore bytes
- **Directory Failed?** → Use embedded bytes
- **All Failed?** → Return helpful error message

> **💡 Pro Tip**: Use this for robust deployments! URL for latest specs, directory for local dev, bytes for offline
> fallback.

## 📖 Comprehensive Examples

Explore real-world implementations in our [examples directory](./examples/):

- **🚀 [Basic Usage](./examples/basic.go)** - Get started in 5 minutes
- **🗂️ [Multi-File Specs](./examples/multi_file_spec.go)** - Organize large APIs
- **🎨 [Customization](./examples/customization.go)** - Themes, layouts, and styling
- **🔧 [Spec Modification](./examples/spec_modification.go)** - Runtime modifications
- **🌐 [URL Loading](./examples/url_based_loading.go)** - Remote spec loading
- **🏗️ [HTTP Integration](./examples/http_server_integration.go)** - Production server setup

> **📚 Interactive Demo**: Visit [https://bdpiprava.github.io/scalar-go](https://bdpiprava.github.io/scalar-go) to see
> all examples in action!

## 🚀 Advanced Use Cases

### 📊 **Multi-Tenant API Documentation**

```go
func tenantSpecificDocs(tenantID string) (string, error) {
    return scalargo.NewV2(
        scalargo.WithSpecURL(fmt.Sprintf("https://specs.company.com/%s/api.yaml", tenantID)),
        scalargo.WithSpecModifier(func (spec *model.Spec) *model.Spec {
            spec.Info.Title = fmt.Sprintf("%s - %s API", strings.Title(tenantID), spec.Info.Title)
            spec.Servers = []model.Server{
                {
                    URL: fmt.Sprintf("https://%s.api.company.com", tenantID)
                },
            }
            return spec
        }),
        scalargo.WithTheme(getThemeForTenant(tenantID)),
    )
}
```

### 🔄 **CI/CD Integration**

```go
// Perfect for automated documentation updates
func generateDocsForBranch(branch string) {
    html, _ := scalargo.NewV2(
        scalargo.WithSpecURL(fmt.Sprintf("https://raw.githubusercontent.com/company/api-specs/%s/openapi.yaml", branch)),
        scalargo.WithMetaDataOpts(
            scalargo.WithTitle(fmt.Sprintf("API Docs - %s branch", branch)),
            scalargo.WithKeyValue("build", os.Getenv("BUILD_NUMBER")),
        ),
    )
// Deploy to branch-specific documentation site
}
```

### 🎨 **White-Label Documentation**

```go
func whitelabelDocs(customerConfig CustomerConfig) (string, error) {
    customCSS := fmt.Sprintf(`
        :root {
            --primary-color: %s;
            --logo-url: url('%s');
        }
        .navbar-brand::before {
            content: '';
            background-image: var(--logo-url);
        }
    `, customerConfig.PrimaryColor, customerConfig.LogoURL)

    return scalargo.NewV2(
        scalargo.WithSpecBytes(customerConfig.APISpec),
        scalargo.WithOverrideCSS(customCSS),
        scalargo.WithMetaDataOpts(
            scalargo.WithTitle(customerConfig.CompanyName + " API"),
        ),
    )
}
```

## 🤝 Contributing

We ❤️ contributions! Here's how you can help:

1. **🐛 Found a Bug?** [Open an issue](https://github.com/bdpiprava/scalar-go/issues)
2. **💡 Have an Idea?** [Start a discussion](https://github.com/bdpiprava/scalar-go/discussions)
3. **🔧 Want to Contribute?** Fork, code, and submit a PR!

### Development Setup

```bash
git clone https://github.com/bdpiprava/scalar-go.git
cd scalar-go
go mod tidy
go test ./...
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Credits & Acknowledgments

- **[Scalar Team](https://github.com/scalar/scalar)** - For creating the amazing Scalar documentation platform that
  powers this library
- **[MarceloPetrucio](https://github.com/MarceloPetrucio/go-scalar-api-reference)** - For the original Go integration
  that inspired this project
- **[OpenAPI Initiative](https://www.openapis.org/)** - For the OpenAPI specification standard
- **Go Community** - For the fantastic ecosystem and tooling

---

<div align="center">

**Made with ❤️ for the Go community**

[⭐ Star this repo](https://github.com/bdpiprava/scalar-go) • [📖 Documentation](https://bdpiprava.github.io/scalar-go) • [🐛 Report Issues](https://github.com/bdpiprava/scalar-go/issues) • [💬 Discussions](https://github.com/bdpiprava/scalar-go/discussions)

</div>
