package loader

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/bdpiprava/scalar-go/model"
	"github.com/bdpiprava/scalar-go/sanitizer"
	"gopkg.in/yaml.v3"
)

// LoadFromDir reads the API specification from the provided root directory
func LoadFromDir(rootDir string, apiFileName string) (*model.Spec, error) {
	content, err := readFile[model.Spec](filepath.Join(rootDir, apiFileName))
	if err != nil {
		return nil, err
	}

	specContent := &content
	specContent.Paths = initializeIfNil(specContent.Paths)
	specContent.Components.Schemas = initializeIfNil(specContent.Components.Schemas)
	specContent.Components.Parameters = initializeIfNil(specContent.Components.Parameters)
	specContent.Components.Responses = initializeIfNil(specContent.Components.Responses)

	paths, err := readDirRecursively(filepath.Join(rootDir, "paths"), "paths")
	if err != nil {
		return nil, err
	}
	maps.Copy(specContent.Paths, *paths)

	responses, err := readDirRecursively(filepath.Join(rootDir, "responses"), "responses")
	if err != nil {
		return nil, err
	}
	maps.Copy(specContent.Components.Responses, *responses)

	schemas, err := readDirRecursively(filepath.Join(rootDir, "schemas"), "schemas")
	if err != nil {
		return nil, err
	}
	maps.Copy(specContent.Components.Schemas, *schemas)

	return sanitizer.Sanitize(specContent), nil
}

// LoadFromDirRoot reads the API specification from the provided root directory
func LoadFromDirRoot(rootDir string) (*model.Spec, error) {
	return LoadFromDir(rootDir, "api.yaml")
}

// LoadFromBytes reads the API specification from the provided bytes in either YAML or JSON format
func LoadFromBytes(bytes []byte) (*model.Spec, error) {
	specContent := &model.Spec{}

	// Try YAML first
	err := yaml.Unmarshal(bytes, specContent)
	if err == nil {
		return sanitizer.Sanitize(specContent), nil
	}

	// Try JSON if YAML fails
	err = json.Unmarshal(bytes, specContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse as YAML or JSON: %w", err)
	}

	return sanitizer.Sanitize(specContent), nil
}

func initializeIfNil(obj model.GenericObject) model.GenericObject {
	if obj != nil {
		return obj
	}
	return model.GenericObject{}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// readFile reads a file and unmarshalls it into the provided data structure.
func readFile[T any](path string) (data T, err error) {
	if data, err = readYamlFile[T](path); err == nil {
		return
	} else if data, err = readJSONFile[T](path); err == nil {
		return
	}
	return data, fmt.Errorf("file '%s' is not a YAML or JSON file, supported extensions are [yml|yaml|json]", path)
}
