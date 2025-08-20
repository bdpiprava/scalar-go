package examples

import (
	"fmt"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/model"
)

// ExampleBasicModification demonstrates basic spec information changes
// @example Basic Modification
// @description This example shows how to dynamically modify API title, description, and version information at runtime.
func ExampleBasicModification() (string, error) {
	description := "This API specification has been dynamically modified to show custom title and description."
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
			spec.Info.Title = "Modified Pet Store API"
			spec.Info.Description = &description
			spec.Info.Version = "2.0.0-modified"
			return spec
		}),
	)
}

// ExampleServerModification demonstrates dynamic server URL modification
// @example Server Modification
// @description This example shows how to add dynamic server URLs based on environment or request, useful for multi-environment deployments.
func ExampleServerModification() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
			spec.Servers = []model.Server{
				{
					URL:         "https://api.example.com/v1",
					Description: toPrt("Production server"),
				},
				{
					URL:         "https://staging-api.example.com/v1",
					Description: toPrt("Staging server"),
				},
				{
					URL:         "http://localhost:8080/api/v1",
					Description: toPrt("Local development server"),
				},
			}
			return spec
		}),
	)
}

// ExampleDynamicInfo demonstrates runtime information updates
// @example Dynamic Information
// @description This example shows how to add dynamic information and tags based on current state, including runtime-generated content and custom tags.
func ExampleDynamicInfo() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
			originalTitle := spec.Info.Title
			spec.Info.Title = fmt.Sprintf("%s (Generated at Runtime)", originalTitle)

			spec.Info.Description = toPrt(fmt.Sprintf(
				"%s\n\n**Note:** This documentation was generated dynamically and includes runtime modifications.",
				spec.Info.Description,
			))

			spec.Tags = append(spec.Tags, model.Tag{
				Name:        "runtime-info",
				Description: "This tag was added dynamically during spec modification",
			})

			return spec
		}),
	)
}

// ExamplePathModification demonstrates working with API paths
// @example Path Analysis
// @description This example shows how to analyze and display information about documented API paths, including endpoint statistics and path listings.
func ExamplePathModification() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
			documentedPaths := spec.DocumentedPaths()

			pathInfo := fmt.Sprintf("\n\n**API Statistics:**\n- Total endpoints: %d\n", len(documentedPaths))

			pathInfo += "\n**Available Endpoints:**\n"
			for _, path := range documentedPaths {
				pathInfo += fmt.Sprintf("- %s %s\n", path.Method, path.Path)
			}

			spec.Info.Description = toPrt(*spec.Info.Description + pathInfo)

			return spec
		}),
	)
}

func toPrt[T any](v T) *T {
	return &v
}
