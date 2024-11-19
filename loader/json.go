package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// readJSONFile reads a JSON file and unmarshalls it into the provided data structure.
func readJSONFile[T any](path string) (T, error) {
	var data T
	if !isJSONFile(path) {
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

// isJSONFile checks if the file is a JSON file.
func isJSONFile(path string) bool {
	return strings.HasSuffix(path, ".json")
}
