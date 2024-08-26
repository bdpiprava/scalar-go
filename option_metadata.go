package scalargo

// MetaData metadata information for Scalar UI
type MetaData map[string]any

type MetaOption func(MetaData)

// WithMetaDataOpts add metadata
func WithMetaDataOpts(metadataOpts ...MetaOption) func(*Options) {
	metadata := MetaData{}

	for _, opt := range metadataOpts {
		opt(metadata)
	}

	return func(o *Options) {
		o.MetaData = metadata
	}
}

// WithTitle add title with value in metadata
func WithTitle(title string) func(MetaData) {
	return WithKeyValue("title", title)
}

// WithKeyValue add metadata with key and value
func WithKeyValue(key string, value any) func(MetaData) {
	return func(md MetaData) {
		md[key] = value
	}
}
