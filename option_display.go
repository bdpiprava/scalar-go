package scalargo

// WithHideSearch hides or shows the search functionality in the Scalar UI
func WithHideSearch(hide bool) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyHideSearch] = hide
	}
}

// WithShowOperationID shows or hides operation IDs in the Scalar UI
func WithShowOperationID(show bool) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyShowOperationID] = show
	}
}

// WithDefaultHTTPClient sets the default HTTP client for code examples
func WithDefaultHTTPClient(target, client string) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyDefaultHTTPClient] = HTTPClientConfig{
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

// WithCustomCSS sets custom CSS in the configuration object (different from WithOverrideCSS which injects CSS in <style> tag)
func WithCustomCSS(css string) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyCustomCSS] = css
	}
}

// WithMultipleSources configures multiple OpenAPI document sources for multi-version API documentation
func WithMultipleSources(sources ...DocumentSource) func(*Options) {
	return func(o *Options) {
		o.Configurations[keySources] = sources
	}
}

// WithShowToolbar controls the visibility of the developer tools toolbar
// Accepts: ShowToolbarAlways, ShowToolbarLocalhost, or ShowToolbarNever
// Default in this library: ShowToolbarNever (differs from Scalar's default of localhost)
func WithShowToolbar(visibility ShowToolbarOption) func(*Options) {
	return func(o *Options) {
		o.Configurations[keyShowToolbar] = string(visibility)
	}
}
