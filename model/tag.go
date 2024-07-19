package model

// Tag structure is generated from "#/$defs/tag".
type Tag struct {
	Name         string                 `json:"name" yaml:"name"`
	Description  string                 `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// ExternalDocumentation structure is generated from "#/$defs/external-documentation".
type ExternalDocumentation struct {
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string  `json:"url" yaml:"url"`
}
