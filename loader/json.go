package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// readJsonFile reads a JSON file and unmarshalls it into the provided data structure.
func readJsonFile[T any](path string) (T, error) {
	var data T
	if !isJsonFile(path) {
		return data, fmt.Errorf("file '%s' is not a JSON file, supported extensions are [JSON]", path)
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(contentBytes, &data)
	if err != nil {
		return data, err
	}

	return data, err
}

// isJsonFile checks if the file is a JSON file.
func isJsonFile(path string) bool {
	return strings.HasSuffix(path, ".json")
}
