package model

type GenericObject map[string]any

type Spec struct {
	OpenAPI    string        `yaml:"openapi" json:"openapi"`
	Info       Info          `yaml:"info" json:"info"`
	Paths      GenericObject `yaml:"paths" json:"paths"`
	Servers    []Server      `yaml:"servers" json:"servers"`
	Tags       []Tag         `yaml:"tags" json:"tags"`
	Components Components    `yaml:"components" json:"components"`
}
