package loader

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/bdpiprava/scalar-go/model"
)

// readYamlFile reads a YAML file and unmarshalls it into the provided data structure.
func readYamlFile[T any](path string) (T, error) {
	var data T
	if !isYamlFile(path) {
		return data, fmt.Errorf("file '%s' is not a YAML file, supported extensions are [yml|yaml]", path)
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(contentBytes, &data)
	if err != nil {
		return data, err
	}

	return data, err
}

// isYamlFile checks if the file is a YAML file.
func isYamlFile(path string) bool {
	return strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")
}

// readDirRecursively reads a directory recursively and returns as model.GenericObject
// rootDir is the original root directory to validate all paths against
// dir is the current directory being read (must be validated to be within rootDir)
func readDirRecursively(rootDir string, dir string, key string) (*model.GenericObject, error) {
	data := model.GenericObject{}
	if !exists(dir) {
		return &data, nil
	}

	// Convert rootDir to absolute path for consistent validation
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve root directory: %w", err)
	}

	// Convert dir to absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve directory: %w", err)
	}

	// Validate that dir is within rootDir
	if !strings.HasPrefix(absDir, absRoot+string(filepath.Separator)) && absDir != absRoot {
		return nil, fmt.Errorf("path traversal detected: '%s' escapes root directory '%s'", dir, rootDir)
	}

	files, err := os.ReadDir(absDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()

		// Construct the target path
		targetPath := filepath.Join(absDir, fileName)

		// Validate that target is within root (defense in depth)
		if !strings.HasPrefix(targetPath, absRoot+string(filepath.Separator)) && targetPath != absRoot {
			return nil, fmt.Errorf("path traversal detected: '%s' escapes root directory", fileName)
		}

		if file.IsDir() {
			subData, err := readDirRecursively(absRoot, targetPath, key)
			if err != nil {
				return &data, err
			}

			maps.Copy(data, *subData)
			continue
		}

		ext := filepath.Ext(fileName)
		fileContent, err := readYamlFile[model.GenericObject](targetPath)
		if err != nil {
			return nil, err
		}

		if len(fileContent) == 0 {
			continue
		}

		if content, ok := fileContent[key]; ok {
			if objContent, ok := content.(model.GenericObject); ok {
				maps.Copy(data, objContent)
			} else {
				return nil, fmt.Errorf("expected object for key '%s' in file '%s', got %T", key, fileName, content)
			}
		} else {
			data[strings.TrimSuffix(fileName, ext)] = fileContent
		}
	}

	return &data, nil
}
