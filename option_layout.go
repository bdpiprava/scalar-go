package scalargo

// Layout represents different layout options
type Layout string

const (
	LayoutModern  Layout = "modern"
	LayoutClassic Layout = "classic"
)

// WithLayout sets the layout for the Scalar UI
func WithLayout(layout Layout) func(*Options) {
	return func(o *Options) {
		o.Layout = layout
	}
}
