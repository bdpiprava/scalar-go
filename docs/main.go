package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/examples"
)

// serverTimeout is the timeout for the server
const serverTimeout = 3 * time.Second

type Example struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Code         string `json:"code"`
	Output       string `json:"output"`
	Category     string `json:"category,omitempty"`
	CategoryIcon string `json:"categoryIcon,omitempty"`
	Icon         string `json:"icon,omitempty"`
}

const loadFromManyFiles = "./data/loader-multiple-files"

// exampleForSpecDir is an example of how to use the scalargo package to load the spec from multiple files
func exampleForSpecDir() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(loadFromManyFiles),
		scalargo.WithBaseFileName("api.yml"),
	)
}

// exampleForSpecURLAndMetadataUsage is an example of how to use the scalargo package to load the spec from a URL and add metadata
func exampleForSpecURLAndMetadataUsage() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Example"),
			scalargo.WithKeyValue("Description", "This is example description"),
		),
	)
}

// exampleForServersOverride is an example of how to use the scalargo package
// to load the spec from a URL and add metadata
func exampleForServersOverride() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		scalargo.WithServers(scalargo.ServerOverride{
			URL:         "http://localhost:8080",
			Description: "Example server",
		}),
	)
}

// exampleForOtherConfigs is an example of how to use the scalargo package to load the spec from a URL and add metadata
func exampleForOtherConfigs() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
		scalargo.WithHiddenClients("fetch", "httr"),
		scalargo.WithHideDarkModeToggle(),
		scalargo.WithOverrideCSS(`
			h1.section-header.tight {
				color: red;
			}
		`),
	)
}

type ExampleFn func() (string, error)

func handler(fn ExampleFn) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		content, err := fn()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(content))
	}
}

var generate = flag.Bool("generate", false, "Generate the static files")

// This is only for local testing
func main() {
	flag.Parse()

	if *generate {
		println("Generating static files")
		buildStatic()
		return
	}

	http.HandleFunc("/spec-dir", handler(exampleForSpecDir))
	http.HandleFunc("/spec-url", handler(exampleForSpecURLAndMetadataUsage))
	http.HandleFunc("/servers-override", handler(exampleForServersOverride))
	http.HandleFunc("/other-configs", handler(exampleForOtherConfigs))

	staticServer := http.FileServer(http.Dir("./build/docs"))
	// Serve static files from build directory
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buildStatic()
		staticServer.ServeHTTP(w, r)
	})

	println("Starting server at http://localhost:8090")
	server := &http.Server{
		Addr:              ":8090",
		ReadHeaderTimeout: serverTimeout,
	}

	_ = server.ListenAndServe()
}

func buildStatic() {
	tmpl, err := template.New("index.html").ParseFiles("./docs/template/index.html")
	if err != nil {
		panic(err)
	}

	f, err := os.Create("./build/docs/index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	exs := getExamples()
	for i, ex := range exs {
		ex.Name = fmt.Sprintf("%d. %s", i+1, ex.Name)
	}

	err = tmpl.Execute(f, exs)
	if err != nil {
		panic(err)
	}
}

func ignoreError(fn ExampleFn) string {
	content, err := fn()
	if err != nil {
		panic(err)
	}
	return content
}

func readFuncBodyIgnoreError(fn reflect.Value) string {
	body, _ := readFuncBody(fn)
	return fmt.Sprintf(`func example() (string, error)%s}`, body)
}

func readFuncBody(fn reflect.Value) (string, error) {
	p := fn.Pointer()
	fc := runtime.FuncForPC(p)
	filename, line := fc.FileLine(p)
	fset := token.NewFileSet()
	// parse file to AST tree
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}
	// walk and find the function block
	find := &FindBlockByLine{Fset: fset, Line: line}
	ast.Walk(find, node)

	if find.Block != nil {
		fp, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer fp.Close()
		_, _ = fp.Seek(int64(find.Block.Lbrace-1), 0)
		buf := make([]byte, int64(find.Block.Rbrace-find.Block.Lbrace))
		_, err = io.ReadFull(fp, buf)
		if err != nil {
			return "", err
		}

		return string(buf), nil
	}
	return "", nil
}

// FindBlockByLine is a ast.Visitor implementation that finds a block by line.
type FindBlockByLine struct {
	Fset  *token.FileSet
	Line  int
	Block *ast.BlockStmt
}

// Visit implements the ast.Visitor interface.
func (f *FindBlockByLine) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	if blockStmt, ok := node.(*ast.BlockStmt); ok {
		stmtStartingPosition := blockStmt.Pos()
		stmtLine := f.Fset.Position(stmtStartingPosition).Line
		if stmtLine == f.Line {
			f.Block = blockStmt
			return nil
		}
	}
	return f
}

func getExamples() []*Example {
	return []*Example{
		// ============================================================
		// Getting Started
		// ============================================================
		{
			Name:         "Basic Usage",
			Description:  `Generate HTML documentation from a single OpenAPI spec file`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleBasicUsage)),
			Output:       ignoreError(examples.ExampleBasicUsage),
			Category:     "Getting Started",
			CategoryIcon: "fa-rocket",
			Icon:         "fa-play",
		},
		{
			Name:         "Multi-File Specification",
			Description:  `Load OpenAPI specs from multiple files in structured directories`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleMultiFileSpec)),
			Output:       ignoreError(examples.ExampleMultiFileSpec),
			Category:     "Getting Started",
			CategoryIcon: "fa-rocket",
			Icon:         "fa-folder-open",
		},
		{
			Name:         "Spec from Bytes",
			Description:  `Load spec from embedded bytes without external files`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleForSpecBytes)),
			Output:       ignoreError(examples.ExampleForSpecBytes),
			Category:     "Getting Started",
			CategoryIcon: "fa-rocket",
			Icon:         "fa-file-code",
		},

		// ============================================================
		// Authentication - API Key
		// ============================================================
		{
			Name:         "Simple API Key",
			Description:  `Basic API key authentication with header-based auth`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIKeySimple)),
			Output:       ignoreError(examples.ExampleAPIKeySimple),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-key",
		},
		{
			Name:         "Custom Header API Key",
			Description:  `API key authentication with custom header name`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIKeyCustomHeader)),
			Output:       ignoreError(examples.ExampleAPIKeyCustomHeader),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-heading",
		},
		{
			Name:         "Query Parameter API Key",
			Description:  `API key passed as query parameter`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIKeyQueryParameter)),
			Output:       ignoreError(examples.ExampleAPIKeyQueryParameter),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-question-circle",
		},
		{
			Name:         "Cookie API Key",
			Description:  `API key passed as cookie value`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIKeyCookie)),
			Output:       ignoreError(examples.ExampleAPIKeyCookie),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-cookie-bite",
		},

		// ============================================================
		// Authentication - HTTP
		// ============================================================
		{
			Name:         "HTTP Basic Auth",
			Description:  `HTTP Basic authentication with username and password`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHTTPBasicAuth)),
			Output:       ignoreError(examples.ExampleHTTPBasicAuth),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-user-lock",
		},
		{
			Name:         "HTTP Bearer Token",
			Description:  `HTTP Bearer token authentication`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHTTPBearerToken)),
			Output:       ignoreError(examples.ExampleHTTPBearerToken),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-shield-alt",
		},

		// ============================================================
		// Authentication - OAuth2
		// ============================================================
		{
			Name:         "OAuth2 Authorization Code",
			Description:  `OAuth2 Authorization Code flow with PKCE`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleOAuth2AuthorizationCode)),
			Output:       ignoreError(examples.ExampleOAuth2AuthorizationCode),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-code-branch",
		},
		{
			Name:         "OAuth2 Client Credentials",
			Description:  `OAuth2 Client Credentials for service-to-service auth`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleOAuth2ClientCredentials)),
			Output:       ignoreError(examples.ExampleOAuth2ClientCredentials),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-server",
		},
		{
			Name:         "OAuth2 Advanced Options",
			Description:  `OAuth2 with custom parameters and advanced configuration`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleOAuth2WithAdvancedOptions)),
			Output:       ignoreError(examples.ExampleOAuth2WithAdvancedOptions),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-sliders-h",
		},
		{
			Name:         "OAuth2 PKCE",
			Description:  `OAuth2 with different PKCE modes (SHA-256, Plain, Disabled)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleOAuth2PKCE)),
			Output:       ignoreError(examples.ExampleOAuth2PKCE),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-fingerprint",
		},

		// ============================================================
		// Authentication - Multiple & Real-World
		// ============================================================
		{
			Name:         "Multiple Security Schemes",
			Description:  `Configure multiple authentication methods (API Key, Bearer, Basic)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleMultipleSecuritySchemes)),
			Output:       ignoreError(examples.ExampleMultipleSecuritySchemes),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-layer-group",
		},
		{
			Name:         "Multiple with OAuth2",
			Description:  `Multiple auth methods including OAuth2`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleMultipleSecuritySchemesWithOAuth2)),
			Output:       ignoreError(examples.ExampleMultipleSecuritySchemesWithOAuth2),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-users-cog",
		},
		{
			Name:         "GitHub-Style Auth",
			Description:  `GitHub-style auth with PAT, OAuth App, and GitHub App tokens`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleGitHubStyleAuth)),
			Output:       ignoreError(examples.ExampleGitHubStyleAuth),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-github",
		},
		{
			Name:         "Stripe-Style Auth",
			Description:  `Stripe-style authentication with test/live key separation`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleStripeStyleAuth)),
			Output:       ignoreError(examples.ExampleStripeStyleAuth),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-credit-card",
		},
		{
			Name:         "Auth0-Style OAuth2",
			Description:  `Auth0-style OAuth2 configuration with audience`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAuth0StyleOAuth2)),
			Output:       ignoreError(examples.ExampleAuth0StyleOAuth2),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-id-badge",
		},
		{
			Name:         "Production API Full Auth",
			Description:  `Production-ready API with comprehensive auth setup`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleProductionAPIWithFullAuth)),
			Output:       ignoreError(examples.ExampleProductionAPIWithFullAuth),
			Category:     "Authentication",
			CategoryIcon: "fa-lock",
			Icon:         "fa-industry",
		},

		// ============================================================
		// Spec Modification
		// ============================================================
		{
			Name:         "Basic Modification",
			Description:  `Dynamically modify API title, description, and version`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleBasicModification)),
			Output:       ignoreError(examples.ExampleBasicModification),
			Category:     "Spec Modification",
			CategoryIcon: "fa-edit",
			Icon:         "fa-pencil-alt",
		},
		{
			Name:         "Server Modification",
			Description:  `Add dynamic server URLs based on environment`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleServerModification)),
			Output:       ignoreError(examples.ExampleServerModification),
			Category:     "Spec Modification",
			CategoryIcon: "fa-edit",
			Icon:         "fa-server",
		},
		{
			Name:         "Dynamic Information",
			Description:  `Add dynamic information and tags at runtime`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDynamicInfo)),
			Output:       ignoreError(examples.ExampleDynamicInfo),
			Category:     "Spec Modification",
			CategoryIcon: "fa-edit",
			Icon:         "fa-sync-alt",
		},
		{
			Name:         "Path Analysis",
			Description:  `Analyze and display API path statistics`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExamplePathModification)),
			Output:       ignoreError(examples.ExamplePathModification),
			Category:     "Spec Modification",
			CategoryIcon: "fa-edit",
			Icon:         "fa-route",
		},

		// ============================================================
		// HTTP Server Integration
		// ============================================================
		{
			Name:         "Static Documentation",
			Description:  `Static documentation with proper caching for production`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleStaticDocumentation)),
			Output:       ignoreError(examples.ExampleStaticDocumentation),
			Category:     "Server Integration",
			CategoryIcon: "fa-network-wired",
			Icon:         "fa-file-alt",
		},
		{
			Name:         "Dynamic Documentation",
			Description:  `Generate documentation with dynamic metadata`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDynamicDocumentation)),
			Output:       ignoreError(examples.ExampleDynamicDocumentation),
			Category:     "Server Integration",
			CategoryIcon: "fa-network-wired",
			Icon:         "fa-sync",
		},
		{
			Name:         "URL-Based Documentation",
			Description:  `Load specification from external URL`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleURLBasedDocumentation)),
			Output:       ignoreError(examples.ExampleURLBasedDocumentation),
			Category:     "Server Integration",
			CategoryIcon: "fa-network-wired",
			Icon:         "fa-link",
		},
		{
			Name:         "API Version 1",
			Description:  `Serve documentation for API v1 with metadata`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIV1)),
			Output:       ignoreError(examples.ExampleAPIV1),
			Category:     "Server Integration",
			CategoryIcon: "fa-network-wired",
			Icon:         "fa-tag",
		},
		{
			Name:         "API Version 2",
			Description:  `Serve documentation for API v2 with enhanced theming`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIV2)),
			Output:       ignoreError(examples.ExampleAPIV2),
			Category:     "Server Integration",
			CategoryIcon: "fa-network-wired",
			Icon:         "fa-tags",
		},
		{
			Name:         "API V1 with Auth",
			Description:  `API v1 documentation with basic authentication`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIV1WithAuth)),
			Output:       ignoreError(examples.ExampleAPIV1WithAuth),
			Category:     "Server Integration",
			CategoryIcon: "fa-network-wired",
			Icon:         "fa-user-shield",
		},
		{
			Name:         "API V2 with OAuth2",
			Description:  `API v2 documentation with OAuth2 authentication`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIV2WithOAuth2)),
			Output:       ignoreError(examples.ExampleAPIV2WithOAuth2),
			Category:     "Server Integration",
			CategoryIcon: "fa-network-wired",
			Icon:         "fa-lock",
		},
		{
			Name:         "Production API",
			Description:  `Production API with environment-specific auth endpoints`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleProductionAPI)),
			Output:       ignoreError(examples.ExampleProductionAPI),
			Category:     "Server Integration",
			CategoryIcon: "fa-network-wired",
			Icon:         "fa-industry",
		},

		// ============================================================
		// External APIs
		// ============================================================
		{
			Name:         "Scalar Galaxy API",
			Description:  `Load Scalar Galaxy API spec from CDN`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleScalarGalaxy)),
			Output:       ignoreError(examples.ExampleScalarGalaxy),
			Category:     "External APIs",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-space-shuttle",
		},
		{
			Name:         "Petstore API",
			Description:  `Classic Petstore OpenAPI specification`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExamplePetstore)),
			Output:       ignoreError(examples.ExamplePetstore),
			Category:     "External APIs",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-paw",
		},
		{
			Name:         "GitHub API",
			Description:  `Complete GitHub REST API documentation`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleGitHubAPI)),
			Output:       ignoreError(examples.ExampleGitHubAPI),
			Category:     "External APIs",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-github",
		},
		{
			Name:         "OpenAI API Demo",
			Description:  `External API with custom titles and theming`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleOpenAIAPI)),
			Output:       ignoreError(examples.ExampleOpenAIAPI),
			Category:     "External APIs",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-brain",
		},
		{
			Name:         "Customized External API",
			Description:  `External spec with comprehensive branding customization`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleCustomizedExternal)),
			Output:       ignoreError(examples.ExampleCustomizedExternal),
			Category:     "External APIs",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-paint-brush",
		},

		// ============================================================
		// Themes
		// ============================================================
		{
			Name:         "Default Theme",
			Description:  `Default theme with clean, modern styling`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeDefault)),
			Output:       ignoreError(examples.ExampleThemeDefault),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-circle",
		},
		{
			Name:         "Alternate Theme",
			Description:  `Alternative theme with distinct styling`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeAlternate)),
			Output:       ignoreError(examples.ExampleThemeAlternate),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-adjust",
		},
		{
			Name:         "Moon Theme",
			Description:  `Dark theme with blue accents for modern interfaces`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeMoon)),
			Output:       ignoreError(examples.ExampleThemeMoon),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-moon",
		},
		{
			Name:         "Purple Theme",
			Description:  `Vibrant purple color scheme for distinctive docs`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemePurple)),
			Output:       ignoreError(examples.ExampleThemePurple),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-crown",
		},
		{
			Name:         "Solarized Theme",
			Description:  `Solarized theme for excellent readability`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeSolarized)),
			Output:       ignoreError(examples.ExampleThemeSolarized),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-sun",
		},
		{
			Name:         "Blue Planet Theme",
			Description:  `Oceanic blue theme with calming aesthetics`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeBluePlanet)),
			Output:       ignoreError(examples.ExampleThemeBluePlanet),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-globe",
		},
		{
			Name:         "Deep Space Theme",
			Description:  `Dark cosmic theme with stellar design elements`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeDeepSpace)),
			Output:       ignoreError(examples.ExampleThemeDeepSpace),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-satellite",
		},
		{
			Name:         "Saturn Theme",
			Description:  `Planetary theme with sophisticated color scheme`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeSaturn)),
			Output:       ignoreError(examples.ExampleThemeSaturn),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-ring",
		},
		{
			Name:         "Kepler Theme",
			Description:  `Astronomical theme inspired by space exploration`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeKepler)),
			Output:       ignoreError(examples.ExampleThemeKepler),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-meteor",
		},
		{
			Name:         "Mars Theme",
			Description:  `Red planet theme with warm, earthy tones`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeMars)),
			Output:       ignoreError(examples.ExampleThemeMars),
			Category:     "Themes",
			CategoryIcon: "fa-palette",
			Icon:         "fa-circle",
		},

		// ============================================================
		// Layouts
		// ============================================================
		{
			Name:         "Modern Layout",
			Description:  `Modern layout with contemporary design`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleLayoutModern)),
			Output:       ignoreError(examples.ExampleLayoutModern),
			Category:     "Layouts",
			CategoryIcon: "fa-th-large",
			Icon:         "fa-layer-group",
		},
		{
			Name:         "Classic Layout",
			Description:  `Classic layout with traditional design`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleLayoutClassic)),
			Output:       ignoreError(examples.ExampleLayoutClassic),
			Category:     "Layouts",
			CategoryIcon: "fa-th-large",
			Icon:         "fa-columns",
		},

		// ============================================================
		// UI Options
		// ============================================================
		{
			Name:         "Hide Sidebar",
			Description:  `Cleaner layout with hidden sidebar`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideSidebar)),
			Output:       ignoreError(examples.ExampleHideSidebar),
			Category:     "UI Options",
			CategoryIcon: "fa-eye",
			Icon:         "fa-eye-slash",
		},
		{
			Name:         "Hide Models",
			Description:  `Focus on endpoints by hiding models section`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideModels)),
			Output:       ignoreError(examples.ExampleHideModels),
			Category:     "UI Options",
			CategoryIcon: "fa-eye",
			Icon:         "fa-border-none",
		},
		{
			Name:         "Dark Mode",
			Description:  `Enable dark mode by default`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDarkMode)),
			Output:       ignoreError(examples.ExampleDarkMode),
			Category:     "UI Options",
			CategoryIcon: "fa-eye",
			Icon:         "fa-adjust",
		},

		// ============================================================
		// Client Options
		// ============================================================
		{
			Name:         "Hide All Clients",
			Description:  `Hide all client code examples from documentation`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideAllClients)),
			Output:       ignoreError(examples.ExampleHideAllClients),
			Category:     "Client Options",
			CategoryIcon: "fa-code",
			Icon:         "fa-eye-slash",
		},
		{
			Name:         "Show Only Curl & Fetch",
			Description:  `Display only curl and fetch client examples`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleShowOnlyCurlAndFetch)),
			Output:       ignoreError(examples.ExampleShowOnlyCurlAndFetch),
			Category:     "Client Options",
			CategoryIcon: "fa-code",
			Icon:         "fa-terminal",
		},
		{
			Name:         "Show All Clients",
			Description:  `Display all available client code examples (default)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleShowAllClients)),
			Output:       ignoreError(examples.ExampleShowAllClients),
			Category:     "Client Options",
			CategoryIcon: "fa-code",
			Icon:         "fa-list",
		},

		// ============================================================
		// Advanced Customization
		// ============================================================
		{
			Name:         "Custom CSS",
			Description:  `Apply custom CSS for branded documentation`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleCustomCSS)),
			Output:       ignoreError(examples.ExampleCustomCSS),
			Category:     "Advanced",
			CategoryIcon: "fa-cogs",
			Icon:         "fa-paint-brush",
		},
		{
			Name:         "All Options Combined",
			Description:  `Comprehensive example combining multiple options`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAllOptions)),
			Output:       ignoreError(examples.ExampleAllOptions),
			Category:     "Advanced",
			CategoryIcon: "fa-cogs",
			Icon:         "fa-star",
		},

		// ============================================================
		// JavaScript API Mode
		// ============================================================
		{
			Name:         "JavaScript API Mode",
			Description:  `Use JavaScript API rendering mode for enhanced control`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleJavaScriptAPIMode)),
			Output:       ignoreError(examples.ExampleJavaScriptAPIMode),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-code",
		},
		{
			Name:         "Custom JavaScript Injection",
			Description:  `Inject custom JavaScript in head and body for advanced initialization`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleJavaScriptAPIWithCustomJS)),
			Output:       ignoreError(examples.ExampleJavaScriptAPIWithCustomJS),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-file-code",
		},
		{
			Name:         "Display Configuration",
			Description:  `Configure search visibility and operation ID display`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideSearchAndShowOperationID)),
			Output:       ignoreError(examples.ExampleHideSearchAndShowOperationID),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-sliders-h",
		},
		{
			Name:         "Default HTTP Client",
			Description:  `Set default HTTP client for code examples`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDefaultHTTPClient)),
			Output:       ignoreError(examples.ExampleDefaultHTTPClient),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-network-wired",
		},
		{
			Name:         "Sorting Options",
			Description:  `Configure sorting for tags, operations, and schema properties`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleSortingOptions)),
			Output:       ignoreError(examples.ExampleSortingOptions),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-sort",
		},
		{
			Name:         "Persist Authentication",
			Description:  `Enable authentication persistence in localStorage`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExamplePersistAuth)),
			Output:       ignoreError(examples.ExamplePersistAuth),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-save",
		},
		{
			Name:         "Multiple API Sources",
			Description:  `Configure multiple OpenAPI document sources with version switcher`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleMultipleSources)),
			Output:       ignoreError(examples.ExampleMultipleSources),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-layer-group",
		},
		{
			Name:         "Custom CSS via Config",
			Description:  `Apply custom CSS through configuration object`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleCustomCSSInConfig)),
			Output:       ignoreError(examples.ExampleCustomCSSInConfig),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-palette",
		},
		{
			Name:         "Advanced Configuration",
			Description:  `Comprehensive JavaScript API configuration with multiple features`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAdvancedConfiguration)),
			Output:       ignoreError(examples.ExampleAdvancedConfiguration),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-tools",
		},
		{
			Name:         "Backward Compatibility",
			Description:  `Verify existing code works without changes (data-attribute mode)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleBackwardCompatibility)),
			Output:       ignoreError(examples.ExampleBackwardCompatibility),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-check-circle",
		},
		{
			Name:         "Toolbar Always Visible",
			Description:  `Enable developer toolbar in all environments for debugging`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleToolbarAlwaysVisible)),
			Output:       ignoreError(examples.ExampleToolbarAlwaysVisible),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-wrench",
		},
		{
			Name:         "Toolbar Localhost Only",
			Description:  `Show developer toolbar only on localhost (Scalar default behavior)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleToolbarLocalhostOnly)),
			Output:       ignoreError(examples.ExampleToolbarLocalhostOnly),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-laptop-code",
		},
		{
			Name:         "Toolbar Hidden",
			Description:  `Hide developer toolbar completely (this library's default)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleToolbarNeverVisible)),
			Output:       ignoreError(examples.ExampleToolbarNeverVisible),
			Category:     "JavaScript API",
			CategoryIcon: "fa-js",
			Icon:         "fa-eye-slash",
		},
	}
}
