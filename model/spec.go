package model

import "fmt"

// GenericObject represets the generic yaml or json object where key is always string and value can be anything
type GenericObject map[string]any

// DocumentedPath holds information about documented path with method
type DocumentedPath struct {
	Path   string
	Method string
}

// Spec represents the OpenAPI spec definition
type Spec struct {
	OpenAPI    string        `yaml:"openapi" json:"openapi"`
	Info       Info          `yaml:"info" json:"info"`
	Paths      GenericObject `yaml:"paths" json:"paths"`
	Servers    []Server      `yaml:"servers" json:"servers"`
	Tags       []Tag         `yaml:"tags" json:"tags"`
	TagsGroup  []TagGroup    `yaml:"x-tagGroups" json:"x-tagGroups"`
	Components Components    `yaml:"components" json:"components"`
}

// DocumentedPaths returns the list of path in the spec
func (s Spec) DocumentedPaths() []DocumentedPath {
	paths := make([]DocumentedPath, 0)
	for path, methods := range s.Paths {
		for method := range methods.(GenericObject) {
			paths = append(paths, DocumentedPath{Path: path, Method: method})
		}
	}
	return paths
}

// String returns method and path joined by '_'
func (d DocumentedPath) String() string {
	return fmt.Sprintf("%s_%s", d.Method, d.Path)
}
