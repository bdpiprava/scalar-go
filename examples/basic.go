package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

const specDir = "./data/loader"
const specFileName = "api.yaml"

// ExampleBasicUsage demonstrates the most basic usage of scalar-go
// @example Basic Usage
// @description This example shows how to generate HTML documentation from a single OpenAPI specification file.
func ExampleBasicUsage() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
	)
}
