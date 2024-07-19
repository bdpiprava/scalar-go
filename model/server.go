package model

// Server structure is generated from "#/$defs/server".
type Server struct {
	URL         string                    `json:"url" yaml:"url"`
	Description *string                   `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// ServerVariable structure is generated from "#/$defs/server-variable".
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"`
	Description *string  `json:"description,omitempty" yaml:"description,omitempty"`
}
