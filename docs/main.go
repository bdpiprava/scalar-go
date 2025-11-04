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

type ConfigDescription struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type TemplateData struct {
	Examples           []*Example                   `json:"examples"`
	ConfigDescriptions map[string]ConfigDescription `json:"configDescriptions"`
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

	templateData := TemplateData{
		Examples:           exs,
		ConfigDescriptions: getConfigDescriptions(),
	}

	err = tmpl.Execute(f, templateData)
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

func getConfigDescriptions() map[string]ConfigDescription {
	return map[string]ConfigDescription{
		"WithTheme": {
			Name:        "theme",
			Type:        "string",
			Description: "Color scheme for the reference (default, alternate, moon, purple, solarized, bluePlanet, deepSpace, saturn, kepler, mars).",
		},
		"WithLayout": {
			Name:        "layout",
			Type:        "'modern' | 'classic'",
			Description: "Layout style for the reference. Modern provides a contemporary design, while classic offers a traditional look.",
		},
		"WithDarkMode": {
			Name:        "darkMode",
			Type:        "boolean",
			Description: "Initial dark mode state. Set to true to enable dark mode by default.",
		},
		"WithHideModels": {
			Name:        "hideModels",
			Type:        "boolean",
			Description: "Hides models in sidebar, search, and content to focus on API endpoints.",
		},
		"WithHideSearch": {
			Name:        "hideSearch",
			Type:        "boolean",
			Description: "Hides the sidebar search bar when set to true.",
		},
		"WithHideDownloadButton": {
			Name:        "hideDownloadButton",
			Type:        "boolean",
			Description: "Hides the download OpenAPI spec button from the interface.",
		},
		"WithSidebarVisibility": {
			Name:        "showSidebar",
			Type:        "boolean",
			Description: "Controls sidebar visibility. Set to false to hide the sidebar for a cleaner layout.",
		},
		"WithHiddenClients": {
			Name:        "hiddenClients",
			Type:        "array | boolean",
			Description: "Controls which HTTP clients are hidden from code examples. Pass an array of client names to hide specific ones, or true to hide all.",
		},
		"WithHideAllClients": {
			Name:        "hiddenClients",
			Type:        "boolean",
			Description: "Hides all HTTP client code examples from the documentation.",
		},
		"WithOverrideCSS": {
			Name:        "customCss",
			Type:        "string",
			Description: "Custom CSS applied directly to the component for branded documentation styling.",
		},
		"WithCustomCSS": {
			Name:        "customCss",
			Type:        "string",
			Description: "Custom CSS applied via configuration object for styling the reference.",
		},
		"WithProxy": {
			Name:        "proxyUrl",
			Type:        "string",
			Description: "Proxy URL for cross-origin API requests.",
		},
		"WithBaseServerURL": {
			Name:        "baseServerURL",
			Type:        "string",
			Description: "Prefix all relative servers with this base URL.",
		},
		"WithServers": {
			Name:        "servers",
			Type:        "Server[]",
			Description: "Overrides servers from the OpenAPI document with custom server configurations.",
		},
		"WithAuthenticationOpts": {
			Name:        "authentication",
			Type:        "AuthenticationConfiguration",
			Description: "Prefills credentials for users. Supports API Key, HTTP Bearer/Basic, and OAuth2 flows.",
		},
		"WithShowOperationID": {
			Name:        "showOperationId",
			Type:        "boolean",
			Description: "Displays operation IDs in the UI for easier reference.",
		},
		"WithDefaultHTTPClient": {
			Name:        "defaultHttpClient",
			Type:        "HttpClientState",
			Description: "Sets the default HTTP client for code examples (e.g., curl, fetch, axios).",
		},
		"WithTagsSorter": {
			Name:        "tagsSorter",
			Type:        "'alpha' | function",
			Description: "Sorts tags alphanumerically or with custom function.",
		},
		"WithOperationsSorter": {
			Name:        "operationsSorter",
			Type:        "'alpha' | 'method' | function",
			Description: "Sorts operations by alpha, method, or custom function.",
		},
		"WithOperationTitleSource": {
			Name:        "operationTitleSource",
			Type:        "'summary' | 'path'",
			Description: "Uses operation summary or path for sidebar display.",
		},
		"WithOrderSchemaPropertiesBy": {
			Name:        "orderSchemaPropertiesBy",
			Type:        "'alpha' | 'preserve'",
			Description: "Sorts properties alphabetically or preserves original order from spec.",
		},
		"WithPersistAuth": {
			Name:        "persistAuth",
			Type:        "boolean",
			Description: "Persists authentication credentials in local storage across page reloads.",
		},
		"WithMultipleSources": {
			Name:        "sources",
			Type:        "DocumentSource[]",
			Description: "Configure multiple OpenAPI document sources with version switcher.",
		},
		"WithCustomHeadJS": {
			Name:        "customHeadJS",
			Type:        "string",
			Description: "Custom JavaScript injected in <head> before Scalar initialization.",
		},
		"WithCustomBodyJS": {
			Name:        "customBodyJS",
			Type:        "string",
			Description: "Custom JavaScript injected in <body> after Scalar initialization.",
		},
		"WithRenderMode": {
			Name:        "renderMode",
			Type:        "'javascript-api' | 'data-attribute'",
			Description: "Rendering mode for Scalar. JavaScript API (default) uses Scalar.createApiReference(), data-attribute (legacy) uses data attributes on script tag.",
		},
		"WithShowToolbar": {
			Name:        "showToolbar",
			Type:        "'always' | 'localhost' | 'never'",
			Description: "Controls developer tools visibility. Set to \"always\" to show in all environments, \"localhost\" for local development only, or \"never\" to hide completely.",
		},
		"WithHideDarkModeToggle": {
			Name:        "hideDarkModeToggle",
			Type:        "boolean",
			Description: "Hides the dark mode toggle button from the interface.",
		},
		"WithForceDarkMode": {
			Name:        "forceDarkModeState",
			Type:        "'dark' | 'light'",
			Description: "Forces dark mode to a specific state, preventing users from changing it.",
		},
		"WithSearchHotKey": {
			Name:        "searchHotKey",
			Type:        "string",
			Description: "Key used with CMD/CTRL to open search (default: \"k\").",
		},
		"WithMetaDataOpts": {
			Name:        "metaData",
			Type:        "object",
			Description: "Configures meta information including title, description, OG tags, and Twitter card data.",
		},
		"WithTitle": {
			Name:        "metaData.title",
			Type:        "string",
			Description: "Sets the page title in metadata. Commonly used for browser tabs, SEO, and social media sharing.",
		},
		"WithKeyValue": {
			Name:        "metaData.<key>",
			Type:        "any",
			Description: "Adds custom key-value pairs to metadata. Use for description, OG tags, Twitter cards, or any custom metadata fields.",
		},
		"WithSpecModifier": {
			Name: "specModifier",
			Type: "func(*model.Spec) *model.Spec",
			Description: `Allows runtime modification of the OpenAPI specification before rendering. Use this to 
dynamically add servers, modify info, add tags, or transform the spec structure.`,
		},
		// Authentication Options
		"WithAPIKey": {
			Name:        "apiKey",
			Type:        "object",
			Description: "Configures API key authentication. Supports header, query parameter, or cookie-based keys with customizable parameter names.",
		},
		"WithAPIKeyName": {
			Name:        "apiKey.name",
			Type:        "string",
			Description: "Sets the parameter name for the API key (e.g., 'X-API-Key', 'api_key').",
		},
		"WithAPIKeyLocation": {
			Name:        "apiKey.in",
			Type:        "'header' | 'query' | 'cookie'",
			Description: "Specifies where the API key is transmitted (header, query parameter, or cookie).",
		},
		"WithAPIKeyHeader": {
			Name:        "apiKey",
			Type:        "object",
			Description: "Shorthand for header-based API key authentication with custom header name.",
		},
		"WithAPIKeyQuery": {
			Name:        "apiKey",
			Type:        "object",
			Description: "Shorthand for query parameter-based API key authentication with custom parameter name.",
		},
		"WithAPIKeyCookie": {
			Name:        "apiKey",
			Type:        "object",
			Description: "Shorthand for cookie-based API key authentication with custom cookie name.",
		},
		"WithHTTPBasicAuth": {
			Name:        "securitySchemes.httpBasic",
			Type:        "object",
			Description: "Configures HTTP Basic authentication with username and password.",
		},
		"WithHTTPBearerToken": {
			Name:        "securitySchemes.httpBearer",
			Type:        "object",
			Description: "Configures HTTP Bearer token authentication (commonly used for JWTs).",
		},
		"WithOAuth2AuthorizationCode": {
			Name:        "securitySchemes.oauth2.flows.authorizationCode",
			Type:        "object",
			Description: "Configures OAuth2 Authorization Code flow with PKCE support. Best for web applications where users grant permission.",
		},
		"WithOAuth2ClientCredentials": {
			Name:        "securitySchemes.oauth2.flows.clientCredentials",
			Type:        "object",
			Description: "Configures OAuth2 Client Credentials flow for service-to-service authentication.",
		},
		"WithOAuth2Implicit": {
			Name:        "securitySchemes.oauth2.flows.implicit",
			Type:        "object",
			Description: "Configures OAuth2 Implicit flow (deprecated but supported for backward compatibility).",
		},
		"WithOAuth2Password": {
			Name:        "securitySchemes.oauth2.flows.password",
			Type:        "object",
			Description: "Configures OAuth2 Password flow (deprecated but supported for backward compatibility).",
		},
		"WithOAuth2ClientID": {
			Name:        "oauth2.clientId",
			Type:        "string",
			Description: "Sets the OAuth2 application client ID.",
		},
		"WithOAuth2ClientSecret": {
			Name:        "oauth2.clientSecret",
			Type:        "string",
			Description: "Sets the OAuth2 application client secret (for Client Credentials flow).",
		},
		"WithOAuth2RedirectURI": {
			Name:        "oauth2.redirectUri",
			Type:        "string",
			Description: "Sets the OAuth2 callback/redirect URI where authorization codes are sent.",
		},
		"WithOAuth2PKCE": {
			Name:        "oauth2.usePkce",
			Type:        "'S256' | 'plain' | 'no'",
			Description: "Enables PKCE (Proof Key for Code Exchange) with specified mode. S256 (SHA-256) recommended for security.",
		},
		"WithOAuth2Scopes": {
			Name:        "oauth2.selectedScopes",
			Type:        "string[]",
			Description: "Pre-selects OAuth2 permission scopes for authorization requests.",
		},
		"WithOAuth2CustomToken": {
			Name:        "oauth2.tokenName",
			Type:        "string",
			Description: "Sets a custom token field name if the OAuth2 provider uses non-standard token response format.",
		},
		"WithOAuth2AdditionalAuthParams": {
			Name:        "oauth2.securityQuery",
			Type:        "object",
			Description: "Adds custom query parameters to the OAuth2 authorization endpoint (e.g., audience, prompt).",
		},
		"WithOAuth2AdditionalTokenParams": {
			Name:        "oauth2.securityBody",
			Type:        "object",
			Description: "Adds custom body parameters to the OAuth2 token endpoint (e.g., resource).",
		},
		"WithOAuth2CredentialsLocation": {
			Name:        "oauth2.credentialsLocation",
			Type:        "'header' | 'body'",
			Description: "Specifies where client credentials are sent during token exchange (header or body).",
		},
		"WithSecurityScheme": {
			Name:        "securitySchemes.<name>",
			Type:        "object",
			Description: "Adds a custom named security scheme. Allows defining multiple authentication methods with unique identifiers.",
		},
		"WithPreferredSecurityScheme": {
			Name:        "preferredSecurityScheme",
			Type:        "string | string[]",
			Description: "Sets the default authentication method shown to users when multiple security schemes are available.",
		},
		"WithCustomSecurity": {
			Name:        "customSecurity",
			Type:        "boolean",
			Description: "Enables custom security scheme configuration, allowing full control over authentication setup.",
		},
		"APIKeyScheme": {
			Name:        "securityScheme",
			Type:        "SecuritySchemeConfig",
			Description: "Creates an API key security scheme configuration for use with WithSecurityScheme.",
		},
		"BearerScheme": {
			Name:        "securityScheme",
			Type:        "SecuritySchemeConfig",
			Description: "Creates a bearer token security scheme configuration for use with WithSecurityScheme.",
		},
		"BasicScheme": {
			Name:        "securityScheme",
			Type:        "SecuritySchemeConfig",
			Description: "Creates a basic auth security scheme configuration for use with WithSecurityScheme.",
		},
		"OAuth2Scheme": {
			Name:        "securityScheme",
			Type:        "SecuritySchemeConfig",
			Description: "Creates an OAuth2 security scheme configuration for use with WithSecurityScheme.",
		},
	}
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
		// Specification Modification
		// ============================================================
		{
			Name:         "Basic Modification",
			Description:  `Dynamically modify API title, description, and version`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleBasicModification)),
			Output:       ignoreError(examples.ExampleBasicModification),
			Category:     "Specification Modification",
			CategoryIcon: "fa-edit",
			Icon:         "fa-pencil-alt",
		},
		{
			Name:         "Server Modification",
			Description:  `Add dynamic server URLs based on environment`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleServerModification)),
			Output:       ignoreError(examples.ExampleServerModification),
			Category:     "Specification Modification",
			CategoryIcon: "fa-edit",
			Icon:         "fa-server",
		},
		{
			Name:         "Dynamic Information",
			Description:  `Add dynamic information and tags at runtime`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDynamicInfo)),
			Output:       ignoreError(examples.ExampleDynamicInfo),
			Category:     "Specification Modification",
			CategoryIcon: "fa-edit",
			Icon:         "fa-sync-alt",
		},
		{
			Name:         "Path Analysis",
			Description:  `Analyze and display API path statistics`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExamplePathModification)),
			Output:       ignoreError(examples.ExamplePathModification),
			Category:     "Specification Modification",
			CategoryIcon: "fa-edit",
			Icon:         "fa-route",
		},

		// ============================================================
		// Server Integration
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
		// External API Examples
		// ============================================================
		{
			Name:         "Scalar Galaxy API",
			Description:  `Load Scalar Galaxy API spec from CDN`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleScalarGalaxy)),
			Output:       ignoreError(examples.ExampleScalarGalaxy),
			Category:     "External API Examples",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-space-shuttle",
		},
		{
			Name:         "Petstore API",
			Description:  `Classic Petstore OpenAPI specification`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExamplePetstore)),
			Output:       ignoreError(examples.ExamplePetstore),
			Category:     "External API Examples",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-paw",
		},
		{
			Name:         "GitHub API",
			Description:  `Complete GitHub REST API documentation`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleGitHubAPI)),
			Output:       ignoreError(examples.ExampleGitHubAPI),
			Category:     "External API Examples",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-github",
		},
		{
			Name:         "OpenAI API Demo",
			Description:  `External API with custom titles and theming`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleOpenAIAPI)),
			Output:       ignoreError(examples.ExampleOpenAIAPI),
			Category:     "External API Examples",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-brain",
		},
		{
			Name:         "Customized External API",
			Description:  `External spec with comprehensive branding customization`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleCustomizedExternal)),
			Output:       ignoreError(examples.ExampleCustomizedExternal),
			Category:     "External API Examples",
			CategoryIcon: "fa-cloud-download-alt",
			Icon:         "fa-paint-brush",
		},

		// ============================================================
		// Theming & Appearance
		// ============================================================
		{
			Name:         "Default Theme",
			Description:  `Default theme with clean, modern styling`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeDefault)),
			Output:       ignoreError(examples.ExampleThemeDefault),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-circle",
		},
		{
			Name:         "Alternate Theme",
			Description:  `Alternative theme with distinct styling`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeAlternate)),
			Output:       ignoreError(examples.ExampleThemeAlternate),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-adjust",
		},
		{
			Name:         "Moon Theme",
			Description:  `Dark theme with blue accents for modern interfaces`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeMoon)),
			Output:       ignoreError(examples.ExampleThemeMoon),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-moon",
		},
		{
			Name:         "Purple Theme",
			Description:  `Vibrant purple color scheme for distinctive docs`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemePurple)),
			Output:       ignoreError(examples.ExampleThemePurple),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-crown",
		},
		{
			Name:         "Solarized Theme",
			Description:  `Solarized theme for excellent readability`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeSolarized)),
			Output:       ignoreError(examples.ExampleThemeSolarized),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-sun",
		},
		{
			Name:         "Blue Planet Theme",
			Description:  `Oceanic blue theme with calming aesthetics`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeBluePlanet)),
			Output:       ignoreError(examples.ExampleThemeBluePlanet),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-globe",
		},
		{
			Name:         "Deep Space Theme",
			Description:  `Dark cosmic theme with stellar design elements`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeDeepSpace)),
			Output:       ignoreError(examples.ExampleThemeDeepSpace),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-satellite",
		},
		{
			Name:         "Saturn Theme",
			Description:  `Planetary theme with sophisticated color scheme`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeSaturn)),
			Output:       ignoreError(examples.ExampleThemeSaturn),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-ring",
		},
		{
			Name:         "Kepler Theme",
			Description:  `Astronomical theme inspired by space exploration`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeKepler)),
			Output:       ignoreError(examples.ExampleThemeKepler),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-meteor",
		},
		{
			Name:         "Mars Theme",
			Description:  `Red planet theme with warm, earthy tones`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeMars)),
			Output:       ignoreError(examples.ExampleThemeMars),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-circle",
		},
		{
			Name:         "Modern Layout",
			Description:  `Modern layout with contemporary design`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleLayoutModern)),
			Output:       ignoreError(examples.ExampleLayoutModern),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-layer-group",
		},
		{
			Name:         "Classic Layout",
			Description:  `Classic layout with traditional design`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleLayoutClassic)),
			Output:       ignoreError(examples.ExampleLayoutClassic),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-columns",
		},
		{
			Name:         "Dark Mode Options",
			Description:  `Configure dark mode with multiple options`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDarkModeOptions)),
			Output:       ignoreError(examples.ExampleDarkModeOptions),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-moon",
		},
		{
			Name:         "Custom CSS (Style Tag)",
			Description:  `Apply custom CSS via style tag for branded documentation`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleCustomCSS)),
			Output:       ignoreError(examples.ExampleCustomCSS),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-paint-brush",
		},
		{
			Name:         "Custom CSS (Config)",
			Description:  `Apply custom CSS through configuration object`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleCustomCSSInConfig)),
			Output:       ignoreError(examples.ExampleCustomCSSInConfig),
			Category:     "Theming & Appearance",
			CategoryIcon: "fa-palette",
			Icon:         "fa-palette",
		},

		// ============================================================
		// UI Visibility
		// ============================================================
		{
			Name:         "Hide Sidebar",
			Description:  `Cleaner layout with hidden sidebar`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideSidebar)),
			Output:       ignoreError(examples.ExampleHideSidebar),
			Category:     "UI Visibility",
			CategoryIcon: "fa-eye",
			Icon:         "fa-eye-slash",
		},
		{
			Name:         "Hide Models",
			Description:  `Focus on endpoints by hiding models section`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideModels)),
			Output:       ignoreError(examples.ExampleHideModels),
			Category:     "UI Visibility",
			CategoryIcon: "fa-eye",
			Icon:         "fa-border-none",
		},
		{
			Name:         "Hide Search",
			Description:  `Hide the search functionality from the sidebar`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideSearch)),
			Output:       ignoreError(examples.ExampleHideSearch),
			Category:     "UI Visibility",
			CategoryIcon: "fa-eye",
			Icon:         "fa-search-minus",
		},
		{
			Name:         "Hide Download Button",
			Description:  `Hide the OpenAPI spec download button`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideDownloadButton)),
			Output:       ignoreError(examples.ExampleHideDownloadButton),
			Category:     "UI Visibility",
			CategoryIcon: "fa-eye",
			Icon:         "fa-download",
		},
		{
			Name:         "Show Operation ID",
			Description:  `Display operation IDs in the UI for easier reference`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleShowOperationID)),
			Output:       ignoreError(examples.ExampleShowOperationID),
			Category:     "UI Visibility",
			CategoryIcon: "fa-eye",
			Icon:         "fa-id-badge",
		},
		{
			Name:         "Toolbar Visibility",
			Description:  `Control developer toolbar visibility (always/localhost/never)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleToolbarVisibility)),
			Output:       ignoreError(examples.ExampleToolbarVisibility),
			Category:     "UI Visibility",
			CategoryIcon: "fa-eye",
			Icon:         "fa-tools",
		},

		// ============================================================
		// Code Example Display
		// ============================================================
		{
			Name:         "Hide All Clients",
			Description:  `Hide all client code examples from documentation`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideAllClients)),
			Output:       ignoreError(examples.ExampleHideAllClients),
			Category:     "Code Example Display",
			CategoryIcon: "fa-code",
			Icon:         "fa-eye-slash",
		},
		{
			Name:         "Show Only Curl & Fetch",
			Description:  `Display only curl and fetch client examples`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleShowOnlyCurlAndFetch)),
			Output:       ignoreError(examples.ExampleShowOnlyCurlAndFetch),
			Category:     "Code Example Display",
			CategoryIcon: "fa-code",
			Icon:         "fa-terminal",
		},
		{
			Name:         "Show All Clients",
			Description:  `Display all available client code examples (default)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleShowAllClients)),
			Output:       ignoreError(examples.ExampleShowAllClients),
			Category:     "Code Example Display",
			CategoryIcon: "fa-code",
			Icon:         "fa-list",
		},
		{
			Name:         "Default HTTP Client",
			Description:  `Set default HTTP client for code examples`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDefaultHTTPClient)),
			Output:       ignoreError(examples.ExampleDefaultHTTPClient),
			Category:     "Code Example Display",
			CategoryIcon: "fa-code",
			Icon:         "fa-network-wired",
		},

		// ============================================================
		// Display Options
		// ============================================================
		{
			Name:         "Tags Sorter",
			Description:  `Sort tags alphabetically in the sidebar`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleTagsSorter)),
			Output:       ignoreError(examples.ExampleTagsSorter),
			Category:     "Display Options",
			CategoryIcon: "fa-sliders-h",
			Icon:         "fa-sort-alpha-down",
		},
		{
			Name:         "Operations Sorter",
			Description:  `Sort operations by method or alphabetically`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleOperationsSorter)),
			Output:       ignoreError(examples.ExampleOperationsSorter),
			Category:     "Display Options",
			CategoryIcon: "fa-sliders-h",
			Icon:         "fa-sort",
		},
		{
			Name:         "Operation Title Source",
			Description:  `Choose operation title source (summary or path)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleOperationTitleSource)),
			Output:       ignoreError(examples.ExampleOperationTitleSource),
			Category:     "Display Options",
			CategoryIcon: "fa-sliders-h",
			Icon:         "fa-heading",
		},
		{
			Name:         "Schema Properties Order",
			Description:  `Control schema property order (alphabetical or preserve)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleSchemaPropertiesOrder)),
			Output:       ignoreError(examples.ExampleSchemaPropertiesOrder),
			Category:     "Display Options",
			CategoryIcon: "fa-sliders-h",
			Icon:         "fa-list-ol",
		},

		// ============================================================
		// Specification Loading
		// ============================================================
		{
			Name:         "Multiple API Sources",
			Description:  `Configure multiple OpenAPI document sources with version switcher`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleMultipleSources)),
			Output:       ignoreError(examples.ExampleMultipleSources),
			Category:     "Specification Loading",
			CategoryIcon: "fa-file-code",
			Icon:         "fa-layer-group",
		},

		// ============================================================
		// Advanced Features
		// ============================================================
		{
			Name:         "Custom JavaScript Injection",
			Description:  `Inject custom JavaScript in head and body for advanced initialization`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleJavaScriptAPIWithCustomJS)),
			Output:       ignoreError(examples.ExampleJavaScriptAPIWithCustomJS),
			Category:     "Advanced Features",
			CategoryIcon: "fa-cogs",
			Icon:         "fa-file-code",
		},
		{
			Name:         "Persist Authentication",
			Description:  `Enable authentication persistence in localStorage`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExamplePersistAuth)),
			Output:       ignoreError(examples.ExamplePersistAuth),
			Category:     "Advanced Features",
			CategoryIcon: "fa-cogs",
			Icon:         "fa-save",
		},
		{
			Name:         "Advanced Configuration",
			Description:  `Comprehensive configuration with multiple features`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAdvancedConfiguration)),
			Output:       ignoreError(examples.ExampleAdvancedConfiguration),
			Category:     "Advanced Features",
			CategoryIcon: "fa-cogs",
			Icon:         "fa-tools",
		},
		{
			Name:         "All Options Combined",
			Description:  `Comprehensive example combining multiple options`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAllOptions)),
			Output:       ignoreError(examples.ExampleAllOptions),
			Category:     "Advanced Features",
			CategoryIcon: "fa-cogs",
			Icon:         "fa-star",
		},

		// ============================================================
		// Rendering Modes
		// ============================================================
		{
			Name:         "JavaScript API Mode",
			Description:  `Use JavaScript API rendering mode (default, recommended)`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleJavaScriptAPIMode)),
			Output:       ignoreError(examples.ExampleJavaScriptAPIMode),
			Category:     "Rendering Modes",
			CategoryIcon: "fa-exchange-alt",
			Icon:         "fa-code",
		},
		{
			Name:         "Data Attribute Mode",
			Description:  `Use legacy data-attribute rendering mode for backward compatibility`,
			Code:         readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDataAttributeMode)),
			Output:       ignoreError(examples.ExampleDataAttributeMode),
			Category:     "Rendering Modes",
			CategoryIcon: "fa-exchange-alt",
			Icon:         "fa-clock",
		},
	}
}
