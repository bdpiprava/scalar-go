package loader

import (
	"maps"
	"os"
	"path/filepath"

	"github.com/bdpiprava/scalar-go/model"
	"github.com/bdpiprava/scalar-go/sanitizer"
)

// LoadWithName reads the API specification from the provided root directory
func LoadWithName(rootDir string, apiFileName string) (*model.Spec, error) {
	content, err := readYamlFile[model.Spec](filepath.Join(rootDir, apiFileName))
	if err != nil {
		return nil, err
	}

	specContent := &content
	specContent.Paths = initializeIfNil(specContent.Paths)
	specContent.Components.Schemas = initializeIfNil(specContent.Components.Schemas)
	specContent.Components.Parameters = initializeIfNil(specContent.Components.Parameters)

	paths, err := readDirRecursively(filepath.Join(rootDir, "paths"), "paths")
	if err != nil {
		return nil, err
	}
	maps.Copy(specContent.Paths, *paths)

	schemas, err := readDirRecursively(filepath.Join(rootDir, "schemas"), "schemas")
	if err != nil {
		return nil, err
	}
	maps.Copy(specContent.Components.Schemas, *schemas)

	return sanitizer.Sanitize(specContent), nil
}

// Load reads the API specification from the provided root directory
func Load(rootDir string) (*model.Spec, error) {
	return LoadWithName(rootDir, "api.yaml")
}

func initializeIfNil(obj model.GenericObject) model.GenericObject {
	if obj != nil {
		return obj
	}
	return model.GenericObject{}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
