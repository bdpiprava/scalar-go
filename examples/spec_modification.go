package examples

import (
	"fmt"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/data"
	"github.com/bdpiprava/scalar-go/model"
)

// ExampleBasicModification demonstrates basic spec information changes
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
func ExampleDynamicInfo() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecDir(specDir),
		scalargo.WithBaseFileName(specFileName),
		scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
			originalTitle := spec.Info.Title
			spec.Info.Title = fmt.Sprintf("%s (Generated at Runtime)", originalTitle)

			spec.Info.Description = toPrt(fmt.Sprintf(
				"%s\n\n**Note:** This documentation was generated dynamically and includes runtime modifications.",
				*spec.Info.Description,
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

			var description string
			if spec.Info.Description != nil {
				description = *spec.Info.Description
			}
			spec.Info.Description = toPrt(description + pathInfo)

			return spec
		}),
	)
}

// ExampleForSpecBytes demonstrates loading a spec from embedded bytes
func ExampleForSpecBytes() (string, error) {
	return scalargo.NewV2(
		scalargo.WithSpecBytes(data.PetStoreSpec),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Pet Store API (Embedded)"),
			scalargo.WithKeyValue("Description", "Self-contained build with embedded spec"),
		),
	)
}

func toPrt[T any](v T) *T {
	return &v
}
