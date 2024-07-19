package model

// Components holds a set of reusable objects for different aspects of the OAS.
type Components struct {
	Schemas    GenericObject `yaml:"schemas" json:"schemas"`
	Parameters GenericObject `yaml:"parameters" json:"parameters"`
}
