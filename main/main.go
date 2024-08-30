package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"reflect"

	scalargo "github.com/bdpiprava/scalar-go"
)

type Example struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
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

// This is only for local testing
func main() {
	http.HandleFunc("/spec-dir", handler(exampleForSpecDir))
	http.HandleFunc("/spec-url", handler(exampleForSpecURLAndMetadataUsage))
	http.HandleFunc("/servers-override", handler(exampleForServersOverride))
	http.HandleFunc("/other-configs", handler(exampleForOtherConfigs))
	http.HandleFunc("/list", examples)
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

func examples(w http.ResponseWriter, request *http.Request) {
	examples := getExamples()

	data, err := json.Marshal(examples)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
}

func getExamples() []Example {
	return []Example{
		{
			Name:        "Reading Spec from Directory",
			Description: "This example shows how to read the spec from multiple files in a directory",
			Path:        "/spec-dir",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(exampleForSpecDir)),
			Output:      ignoreError(exampleForSpecDir),
		},
		{
			Name:        "Reading Spec from URL and Metadata Usage",
			Description: "This example shows how to read the spec from a URL and add metadata",
			Path:        "/spec-url",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(exampleForSpecURLAndMetadataUsage)),
			Output:      ignoreError(exampleForSpecURLAndMetadataUsage),
		},
		{
			Name:        "Servers Override",
			Description: "This example shows how to read the spec from a URL and override the servers",
			Path:        "/servers-override",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(exampleForServersOverride)),
			Output:      ignoreError(exampleForServersOverride),
		},
		{
			Name:        "Other Configs",
			Description: "This example shows how to read the spec from a URL and add other configurations",
			Path:        "/other-configs",
			Code:        readFuncBodyIgnoreError(reflect.ValueOf(exampleForOtherConfigs)),
			Output:      ignoreError(exampleForOtherConfigs),
		},
	}
}

func ignoreError(fn ExampleFn) string {
	content, _ := fn()
	return content
}
