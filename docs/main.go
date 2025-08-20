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
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Output      string `json:"output"`
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
		scalargo.WithServers(scalargo.Server{
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
	http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		buildStatic()
		http.FileServer(http.Dir("./main/static")).ServeHTTP(w, request)
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

	f, err := os.Create("./docs/static/index.html")
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

func getExamples() []*Example {
	return []*Example{
		{
			Name:        "Reading Spec from Directory",
			Description: "This example shows how to read the spec from multiple files in a directory",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(exampleForSpecDir)),
			Output:      ignoreError(exampleForSpecDir),
		},
		{
			Name:        "Reading Spec from URL and Metadata Usage",
			Description: "This example shows how to read the spec from a URL and add metadata",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(exampleForSpecURLAndMetadataUsage)),
			Output:      ignoreError(exampleForSpecURLAndMetadataUsage),
		},
		{
			Name:        "Servers Override",
			Description: "This example shows how to read the spec from a URL and override the servers",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(exampleForServersOverride)),
			Output:      ignoreError(exampleForServersOverride),
		},
		{
			Name:        "Other Configs",
			Description: "This example shows how to read the spec from a URL and add other configurations",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(exampleForOtherConfigs)),
			Output:      ignoreError(exampleForOtherConfigs),
		},

		// Basic Examples
		{
			Name:        "Basic Usage",
			Description: "This example shows how to generate HTML documentation from a single OpenAPI specification file.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleBasicUsage)),
			Output:      ignoreError(examples.ExampleBasicUsage),
		},

		// Multi-file Spec Examples
		{
			Name:        "Multi-File Specification",
			Description: "This example shows how to load OpenAPI specs from multiple files organized in a structured directory layout, useful for large APIs where you want to split schemas, paths, and responses into separate files.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleMultiFileSpec)),
			Output:      ignoreError(examples.ExampleMultiFileSpec),
		},

		// Spec Modification Examples
		{
			Name:        "Basic Modification",
			Description: "This example shows how to dynamically modify API title, description, and version information at runtime.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleBasicModification)),
			Output:      ignoreError(examples.ExampleBasicModification),
		},
		{
			Name:        "Server Modification",
			Description: "This example shows how to add dynamic server URLs based on environment or request, useful for multi-environment deployments.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleServerModification)),
			Output:      ignoreError(examples.ExampleServerModification),
		},
		{
			Name:        "Dynamic Information",
			Description: "This example shows how to add dynamic information and tags based on current state, including runtime-generated content and custom tags.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDynamicInfo)),
			Output:      ignoreError(examples.ExampleDynamicInfo),
		},
		{
			Name:        "Path Analysis",
			Description: "This example shows how to analyze and display information about documented API paths, including endpoint statistics and path listings.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExamplePathModification)),
			Output:      ignoreError(examples.ExamplePathModification),
		},

		// HTTP Server Integration Examples
		{
			Name:        "Static Documentation",
			Description: "This example shows how to create simple static documentation with proper headers and caching for production use.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleStaticDocumentation)),
			Output:      ignoreError(examples.ExampleStaticDocumentation),
		},
		{
			Name:        "Dynamic Documentation",
			Description: "This example shows how to generate documentation with custom metadata based on request or environment, including dynamic titles and timestamps.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDynamicDocumentation)),
			Output:      ignoreError(examples.ExampleDynamicDocumentation),
		},
		{
			Name:        "URL-Based Documentation",
			Description: "This example shows how to load specification from an external URL, useful for documentation services or when specs are hosted elsewhere.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleURLBasedDocumentation)),
			Output:      ignoreError(examples.ExampleURLBasedDocumentation),
		},
		{
			Name:        "API Version 1",
			Description: "This example shows how to serve documentation for multiple API versions with version-specific metadata and theming.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIV1)),
			Output:      ignoreError(examples.ExampleAPIV1),
		},
		{
			Name:        "API Version 2",
			Description: "This example shows how to serve documentation for a newer API version using external specifications with enhanced theming and dark mode.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAPIV2)),
			Output:      ignoreError(examples.ExampleAPIV2),
		},

		// URL-based Loading Examples
		{
			Name:        "Scalar Galaxy API",
			Description: "This example shows how to load the Scalar Galaxy API specification directly from a CDN URL.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleScalarGalaxy)),
			Output:      ignoreError(examples.ExampleScalarGalaxy),
		},
		{
			Name:        "Petstore API",
			Description: "This example shows how to load the classic Petstore OpenAPI specification from the official repository.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExamplePetstore)),
			Output:      ignoreError(examples.ExamplePetstore),
		},
		{
			Name:        "GitHub API",
			Description: "This example shows how to load and display the complete GitHub REST API documentation using their public OpenAPI specification.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleGitHubAPI)),
			Output:      ignoreError(examples.ExampleGitHubAPI),
		},
		{
			Name:        "OpenAI API Demo",
			Description: "This example shows how to load external API documentation with custom titles, descriptions, and alternative theming (using Scalar Galaxy as demo).",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleOpenAIAPI)),
			Output:      ignoreError(examples.ExampleOpenAIAPI),
		},
		{
			Name:        "Customized External API",
			Description: "This example shows how to load external specifications with comprehensive examples. including custom CSS, theming, and UI options for branded documentation.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleCustomizedExternal)),
			Output:      ignoreError(examples.ExampleCustomizedExternal),
		},

		// Theme Examples
		{
			Name:        "Default Theme",
			Description: "This example shows the default theme with clean, modern styling for professional API documentation.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeDefault)),
			Output:      ignoreError(examples.ExampleThemeDefault),
		},
		{
			Name:        "Moon Theme",
			Description: "This example shows the moon theme with dark styling and blue accents, perfect for modern dark-mode preferences.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeMoon)),
			Output:      ignoreError(examples.ExampleThemeMoon),
		},
		{
			Name:        "Purple Theme",
			Description: "This example shows the purple theme with vibrant purple color scheme for distinctive and creative API documentation.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemePurple)),
			Output:      ignoreError(examples.ExampleThemePurple),
		},
		{
			Name:        "Solarized Theme",
			Description: "This example shows the solarized theme based on the popular Solarized color scheme, offering excellent readability and reduced eye strain.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleThemeSolarized)),
			Output:      ignoreError(examples.ExampleThemeSolarized),
		},

		// Layout Examples
		{
			Name:        "Modern Layout",
			Description: "This example shows the modern layout with contemporary design elements and enhanced user experience for API documentation.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleLayoutModern)),
			Output:      ignoreError(examples.ExampleLayoutModern),
		},
		{
			Name:        "Classic Layout",
			Description: "This example shows the classic layout with traditional documentation design, familiar to users of conventional API documentation tools.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleLayoutClassic)),
			Output:      ignoreError(examples.ExampleLayoutClassic),
		},

		// Visibility Examples
		{
			Name:        "Hide Sidebar",
			Description: "This example shows how to hide the sidebar to create a cleaner, more focused documentation layout with more space for content.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideSidebar)),
			Output:      ignoreError(examples.ExampleHideSidebar),
		},
		{
			Name:        "Hide Models",
			Description: "This example shows how to hide the models section to focus purely on API endpoints, useful for endpoint-centric documentation.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleHideModels)),
			Output:      ignoreError(examples.ExampleHideModels),
		},
		{
			Name:        "Dark Mode",
			Description: "This example shows how to enable dark mode by default, providing a modern dark interface that's easier on the eyes.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleDarkMode)),
			Output:      ignoreError(examples.ExampleDarkMode),
		},

		// Advanced Customization Examples
		{
			Name:        "Custom CSS",
			Description: "This example shows how to apply custom CSS overrides to create branded documentation with custom colors, fonts, and styling elements.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleCustomCSS)),
			Output:      ignoreError(examples.ExampleCustomCSS),
		},
		{
			Name:        "All Options Combined",
			Description: "This example shows how to combine multiple examples. options including theme, layout, UI controls, custom CSS, and client hiding for comprehensive documentation branding.",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(examples.ExampleAllOptions)),
			Output:      ignoreError(examples.ExampleAllOptions),
		},
	}
}

func ignoreError(fn ExampleFn) string {
	content, _ := fn()
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
