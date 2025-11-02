package loader

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"github.com/bdpiprava/scalar-go/model"
	"github.com/bdpiprava/scalar-go/sanitizer"
	"gopkg.in/yaml.v3"
)

// validatePath ensures the resolved path is within the allowed root directory
// This prevents path traversal attacks using "../" sequences
func validatePath(rootDir, targetPath string) (string, error) {
	// Clean and resolve the root directory to absolute path
	cleanRoot, err := filepath.Abs(filepath.Clean(rootDir))
	if err != nil {
		return "", fmt.Errorf("invalid root directory: %w", err)
	}

	// Clean the target path to remove ".." and "./" sequences
	cleanTarget := filepath.Clean(targetPath)

	// Join and resolve to absolute path
	fullPath, err := filepath.Abs(filepath.Join(cleanRoot, cleanTarget))
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	// Ensure fullPath starts with cleanRoot followed by separator
	// This prevents escaping the root directory via path traversal
	if !strings.HasPrefix(fullPath, cleanRoot+string(filepath.Separator)) && fullPath != cleanRoot {
		return "", fmt.Errorf("path traversal detected: '%s' escapes root directory '%s'", targetPath, rootDir)
	}

	return fullPath, nil
}

// LoadFromDir reads the API specification from the provided root directory
func LoadFromDir(rootDir string, apiFileName string) (*model.Spec, error) {
	// Validate the main API file path
	apiFilePath, err := validatePath(rootDir, apiFileName)
	if err != nil {
		return nil, err
	}

	content, err := readFile[model.Spec](apiFilePath)
	if err != nil {
		return nil, err
	}

	specContent := &content
	specContent.Paths = initializeIfNil(specContent.Paths)
	specContent.Components.Schemas = initializeIfNil(specContent.Components.Schemas)
	specContent.Components.Parameters = initializeIfNil(specContent.Components.Parameters)
	specContent.Components.Responses = initializeIfNil(specContent.Components.Responses)

	// Validate paths subdirectory
	pathsDir, err := validatePath(rootDir, "paths")
	if err != nil {
		return nil, err
	}
	paths, err := readDirRecursively(rootDir, pathsDir, "paths")
	if err != nil {
		return nil, err
	}
	maps.Copy(specContent.Paths, *paths)

	// Validate responses subdirectory
	responsesDir, err := validatePath(rootDir, "responses")
	if err != nil {
		return nil, err
	}
	responses, err := readDirRecursively(rootDir, responsesDir, "responses")
	if err != nil {
		return nil, err
	}
	maps.Copy(specContent.Components.Responses, *responses)

	// Validate schemas subdirectory
	schemasDir, err := validatePath(rootDir, "schemas")
	if err != nil {
		return nil, err
	}
	schemas, err := readDirRecursively(rootDir, schemasDir, "schemas")
	if err != nil {
		return nil, err
	}
	maps.Copy(specContent.Components.Schemas, *schemas)

	return sanitizer.Sanitize(specContent), nil
}

// Load reads the API specification from the provided root directory
func Load(rootDir string) (*model.Spec, error) {
	return LoadFromDirRoot(rootDir)
}

// LoadWithName reads the API specification from the provided root directory
func LoadWithName(rootDir, apiFileName string) (*model.Spec, error) {
	return LoadFromDir(rootDir, apiFileName)
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
