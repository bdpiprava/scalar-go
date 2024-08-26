package scalargo

// Server represtnts server override configuration
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// WithServers servers to override the openapi spec servers
func WithServers(servers ...Server) func(*Options) {
	return func(o *Options) {
		o.Servers = servers
	}
}
