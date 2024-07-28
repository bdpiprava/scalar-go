package model

type GenericObject map[string]any

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

func (s Spec) DocumentedPaths() []DocumentedPath {
	paths := make([]DocumentedPath, 0)
	for path, methods := range s.Paths {
		for method, _ := range methods.(GenericObject) {
			paths = append(paths, DocumentedPath{Path: path, Method: method})
		}
	}
	return paths
}
