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
func readDirRecursively(dir string, key string) (*model.GenericObject, error) {
	data := model.GenericObject{}
	if !exists(dir) {
		return &data, nil
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fileName := file.Name()
		if file.IsDir() {
			subData, err := readDirRecursively(filepath.Join(dir, fileName), key)
			if err != nil {
				return &data, err
			}

			maps.Copy(data, *subData)
			continue
		}

		ext := filepath.Ext(fileName)
		fileContent, err := readYamlFile[model.GenericObject](filepath.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		if len(fileContent) == 0 {
			continue
		}

		if content, ok := fileContent[key]; ok {
			maps.Copy(data, content.(model.GenericObject))
		} else {
			data[strings.TrimSuffix(fileName, ext)] = fileContent
		}
	}

	return &data, nil
}
