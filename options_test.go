package scalargo_test

import (
	"testing"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/stretchr/testify/assert"
)

func TestWithHiddenClients(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		options        []scalargo.Option
		wantHiddenClients any
	}{
		{
			name: "should set hidden clients with single client",
			options: []scalargo.Option{
				scalargo.WithHiddenClients("curl"),
			},
			wantHiddenClients: []string{"curl"},
		},
		{
			name: "should set hidden clients with multiple clients",
			options: []scalargo.Option{
				scalargo.WithHiddenClients("curl", "wget", "httpie"),
			},
			wantHiddenClients: []string{"curl", "wget", "httpie"},
		},
		{
			name: "should set hidden clients to nil slice when no clients provided",
			options: []scalargo.Option{
				scalargo.WithHiddenClients(),
			},
			wantHiddenClients: []string(nil),
		},
		{
			name: "should overwrite previous hidden clients when called twice",
			options: []scalargo.Option{
				scalargo.WithHiddenClients("curl"),
				scalargo.WithHiddenClients("wget", "httpie"),
			},
			wantHiddenClients: []string{"wget", "httpie"},
		},
		{
			name: "should ignore WithHiddenClients when WithHideAllClients is called first",
			options: []scalargo.Option{
				scalargo.WithHideAllClients(),
				scalargo.WithHiddenClients("curl", "wget"),
			},
			wantHiddenClients: true,
		},
		{
			name: "should be overridden by WithHideAllClients when called after",
			options: []scalargo.Option{
				scalargo.WithHiddenClients("curl", "wget"),
				scalargo.WithHideAllClients(),
			},
			wantHiddenClients: true,
		},
		{
			name: "should handle WithHiddenClients called between WithHideAllClients calls",
			options: []scalargo.Option{
				scalargo.WithHideAllClients(),
				scalargo.WithHiddenClients("curl"),
				scalargo.WithHideAllClients(),
			},
			wantHiddenClients: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := &scalargo.Options{
				Configurations: make(map[string]any),
			}

			for _, opt := range tc.options {
				opt(opts)
			}

			got := opts.Configurations["hiddenClients"]
			assert.Equal(t, tc.wantHiddenClients, got)
		})
	}
}
