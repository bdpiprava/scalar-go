package main

import (
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
)

// This is only for local testing
func main() {
	apiDir := "./data/loader-multiple-files"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := scalargo.New(
			apiDir,
			scalargo.WithBaseFileName("api.yml"),
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write([]byte(content))
	})

	_ = http.ListenAndServe(":8090", nil)
}
