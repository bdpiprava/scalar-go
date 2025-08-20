package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

// ExampleMultiFileSpec demonstrates loading OpenAPI specifications from multiple files
// @example Multi-File Specification
// @description This example shows how to load OpenAPI specs from multiple files organized in a structured directory layout, useful for large APIs where you want to split schemas, paths, and responses into separate files.
func ExampleMultiFileSpec() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
	)
}
