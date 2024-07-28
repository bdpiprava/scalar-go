package model

// TagGroup represent the element of the x-tagGroups
type TagGroup struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Tags        []string `yaml:"tags" json:"tags"`
}
