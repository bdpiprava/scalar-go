package scalargo

// ServerOverride represents server override configuration for UI display.
// This is a simplified version of model.Server, used to quickly override
// server URLs without needing the full OpenAPI server specification.
type ServerOverride struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// Server is alias exists for backwards compatibility but will be removed in a future version.
// The name was changed to avoid collision with model.Server.
//
// Deprecated: use ServerOverride instead.
type Server = ServerOverride

// WithServers servers to override the openapi spec servers
func WithServers(servers ...ServerOverride) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyServers] = servers
	}
}
