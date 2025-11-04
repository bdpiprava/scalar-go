package scalargo

// HTTPClientConfig configures the default HTTP client to use in the Scalar UI
type HTTPClientConfig struct {
	TargetKey string `json:"targetKey"` // Target language/platform (e.g., "node", "php", "python")
	ClientKey string `json:"clientKey"` // Specific client library (e.g., "undici", "guzzle", "requests")
}

// DocumentSource represents a single OpenAPI document source for multi-document configurations
type DocumentSource struct {
	Title   string `json:"title,omitempty"`   // Display title for this document
	Slug    string `json:"slug,omitempty"`    // URL slug for routing
	URL     string `json:"url,omitempty"`     // URL to fetch the OpenAPI document
	Content string `json:"content,omitempty"` // Inline OpenAPI document content (mutually exclusive with URL)
	Default bool   `json:"default,omitempty"` // Whether this is the default document to display
}

// SorterOption defines sorting behavior for tags and operations
type SorterOption string

const (
	// SorterAlpha sorts alphabetically
	SorterAlpha SorterOption = "alpha"
	// SorterMethod sorts by HTTP method (for operations only)
	SorterMethod SorterOption = "method"
)

// OperationTitleSource defines where to get operation titles from
type OperationTitleSource string

const (
	// OperationTitleSourceSummary uses the summary field
	OperationTitleSourceSummary OperationTitleSource = "summary"
	// OperationTitleSourcePath uses the path
	OperationTitleSourcePath OperationTitleSource = "path"
)

// SchemaPropertiesOrder defines how to order schema properties
type SchemaPropertiesOrder string

const (
	// SchemaPropertiesOrderAlpha sorts properties alphabetically
	SchemaPropertiesOrderAlpha SchemaPropertiesOrder = "alpha"
	// SchemaPropertiesOrderPreserve preserves the order from the spec
	SchemaPropertiesOrderPreserve SchemaPropertiesOrder = "preserve"
)

// ShowToolbarOption defines when to display the developer tools toolbar
type ShowToolbarOption string

const (
	// ShowToolbarAlways displays the toolbar in all environments
	ShowToolbarAlways ShowToolbarOption = "always"
	// ShowToolbarLocalhost displays toolbar only on localhost or similar hosts (Scalar default)
	ShowToolbarLocalhost ShowToolbarOption = "localhost"
	// ShowToolbarNever never displays the toolbar
	ShowToolbarNever ShowToolbarOption = "never"
)
