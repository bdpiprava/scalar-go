package main

import (
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
)

// This is only for local testing
func main() {
	apiDir := "./data/loader-multiple-files"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := scalargo.NewV2(
			scalargo.WithSpecDir(apiDir),
			scalargo.WithBaseFileName("api.yml"),
			scalargo.WithMetaDataOpts(
				scalargo.WithTitle("Example"),
				scalargo.WithKeyValue("ogTitle", "This is example"),
			),
			scalargo.WithSpecURL("https://cdn.jsdelivr.net/npm/@scalar/galaxy/dist/latest.yaml"),
			scalargo.WithServers(scalargo.Server{
				URL:         "http://localhost:8080",
				Description: "Example server",
			}),
			scalargo.WithHiddenClients("fetch", "httr"),
			scalargo.WithHideDarkModeToggle(),
			scalargo.WithOverrideCSS(`
				h1.section-header.tight {
					color: red;
				}
			`),
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write([]byte(content))
	})

	println("Starting server at http://localhost:8090")
	_ = http.ListenAndServe(":8090", nil)
}
