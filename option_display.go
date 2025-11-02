package scalargo

// WithHideSearch hides or shows the search functionality in the Scalar UI
func WithHideSearch(hide bool) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyHideSearch] = hide
	}
}

// WithShowOperationId shows or hides operation IDs in the Scalar UI
func WithShowOperationId(show bool) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyShowOperationId] = show
	}
}

// WithDefaultHttpClient sets the default HTTP client for code examples
func WithDefaultHttpClient(target, client string) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyDefaultHttpClient] = HttpClientConfig{
			TargetKey: target,
			ClientKey: client,
		}
	}
}

// WithTagsSorter sets how tags are sorted in the sidebar
func WithTagsSorter(sorter SorterOption) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyTagsSorter] = string(sorter)
	}
}

// WithOperationsSorter sets how operations are sorted within tags
func WithOperationsSorter(sorter SorterOption) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyOperationsSorter] = string(sorter)
	}
}

// WithOperationTitleSource sets where to get operation titles from (summary or path)
func WithOperationTitleSource(source OperationTitleSource) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyOperationTitleSource] = string(source)
	}
}

// WithOrderSchemaPropertiesBy sets how schema properties are ordered
func WithOrderSchemaPropertiesBy(order SchemaPropertiesOrder) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyOrderSchemaPropertiesBy] = string(order)
	}
}

// WithPersistAuth enables or disables persisting authentication credentials in localStorage
func WithPersistAuth(persist bool) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyPersistAuth] = persist
	}
}

// WithCustomCss sets custom CSS in the configuration object (different from WithOverrideCSS which injects CSS in <style> tag)
func WithCustomCss(css string) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyCustomCss] = css
	}
}

// WithMultipleSources configures multiple OpenAPI document sources for multi-version API documentation
func WithMultipleSources(sources ...DocumentSource) func(*Options) {
	return func(o *Options) {
		o.Configurations[keySources] = sources
	}
}
