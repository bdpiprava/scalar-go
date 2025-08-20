package examples

import (
	scalargo "github.com/bdpiprava/scalar-go"
)

const specDir = "./data/loader"
const specFileName = "pet-store.yml"

// ExampleBasicUsage demonstrates the most basic usage of scalar-go
func ExampleBasicUsage() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
	)
}
