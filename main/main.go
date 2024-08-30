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

	scalargo "github.com/bdpiprava/scalar-go"
)

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

// exampleForServersOverride is an example of how to use the scalargo package to load the spec from a URL and add metadata
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
	return func(w http.ResponseWriter, request *http.Request) {
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
	_ = http.ListenAndServe(":8090", nil)
}

func buildStatic() {
	tmpl, err := template.New("index.html").ParseFiles("./main/template/index.html")
	if err != nil {
		panic(err)
	}

	f, err := os.Create("./main/static/index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = tmpl.Execute(f, getExamples())
	if err != nil {
		panic(err)
	}
}

func getExamples() []Example {
	return []Example{
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
