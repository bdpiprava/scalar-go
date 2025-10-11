package model

// Components holds a set of reusable objects for different aspects of the OAS.
type Components struct {
	Schemas         GenericObject `yaml:"schemas,omitempty" json:"schemas,omitempty"`
	Parameters      GenericObject `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Responses       GenericObject `yaml:"responses,omitempty" json:"responses,omitempty"`
	Examples        GenericObject `yaml:"examples,omitempty" json:"examples,omitempty"`
	RequestBodies   GenericObject `yaml:"requestBodies,omitempty" json:"requestBodies,omitempty"`
	Headers         GenericObject `yaml:"headers,omitempty" json:"headers,omitempty"`
	SecuritySchemes GenericObject `yaml:"securitySchemes,omitempty" json:"securitySchemes,omitempty"`
	Link            GenericObject `yaml:"links,omitempty" json:"links,omitempty"`
	Callbacks       GenericObject `yaml:"callbacks,omitempty" json:"callbacks,omitempty"`
	PathItems       GenericObject `yaml:"pathItems,omitempty" json:"pathItems,omitempty"`
}
