package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

// ExampleMultiFileSpec demonstrates loading OpenAPI specifications from multiple files
func ExampleMultiFileSpec() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
	)
}
