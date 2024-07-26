package model

// Components holds a set of reusable objects for different aspects of the OAS.
type Components struct {
	Schemas         GenericObject `yaml:"schemas" json:"schemas"`
	Parameters      GenericObject `yaml:"parameters" json:"parameters"`
	Responses       GenericObject `yaml:"responses" json:"responses"`
	Examples        GenericObject `yaml:"examples" json:"examples"`
	RequestBodies   GenericObject `yaml:"requestBodies" json:"requestBodies"`
	Headers         GenericObject `yaml:"headers" json:"headers"`
	SecuritySchemes GenericObject `yaml:"securitySchemes" json:"securitySchemes"`
	Link            GenericObject `yaml:"links" json:"links"`
	Callbacks       GenericObject `yaml:"callbacks" json:"callbacks"`
	PathItems       GenericObject `yaml:"pathItems" json:"pathItems"`
}
