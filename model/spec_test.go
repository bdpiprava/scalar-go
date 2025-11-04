package model_test

import (
	"sort"
	"testing"

	"github.com/bdpiprava/scalar-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentedPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec model.Spec
		want []model.DocumentedPath
	}{
		{
			name: "should return documented paths with methods",
			spec: model.Spec{
				Paths: model.GenericObject{
					"/users": model.GenericObject{
						"get":  model.GenericObject{"summary": "Get users"},
						"post": model.GenericObject{"summary": "Create user"},
					},
					"/pets": model.GenericObject{
						"get": model.GenericObject{"summary": "Get pets"},
					},
				},
			},
			want: []model.DocumentedPath{
				{Path: "/pets", Method: "get"},
				{Path: "/users", Method: "get"},
				{Path: "/users", Method: "post"},
			},
		},
		{
			name: "should return empty slice for empty paths",
			spec: model.Spec{
				Paths: model.GenericObject{},
			},
			want: []model.DocumentedPath{},
		},
		{
			name: "should handle single path with multiple methods",
			spec: model.Spec{
				Paths: model.GenericObject{
					"/api": model.GenericObject{
						"get":    model.GenericObject{},
						"post":   model.GenericObject{},
						"put":    model.GenericObject{},
						"delete": model.GenericObject{},
						"patch":  model.GenericObject{},
					},
				},
			},
			want: []model.DocumentedPath{
				{Path: "/api", Method: "delete"},
				{Path: "/api", Method: "get"},
				{Path: "/api", Method: "patch"},
				{Path: "/api", Method: "post"},
				{Path: "/api", Method: "put"},
			},
		},
		{
			name: "should handle paths with special characters",
			spec: model.Spec{
				Paths: model.GenericObject{
					"/users/{userId}": model.GenericObject{
						"get": model.GenericObject{},
					},
					"/items/{item_id}/details": model.GenericObject{
						"get": model.GenericObject{},
					},
				},
			},
			want: []model.DocumentedPath{
				{Path: "/items/{item_id}/details", Method: "get"},
				{Path: "/users/{userId}", Method: "get"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.spec.DocumentedPaths()

			// Sort for consistent comparison
			sort.Slice(got, func(i, j int) bool {
				if got[i].Path == got[j].Path {
					return got[i].Method < got[j].Method
				}
				return got[i].Path < got[j].Path
			})

			sort.Slice(tc.want, func(i, j int) bool {
				if tc.want[i].Path == tc.want[j].Path {
					return tc.want[i].Method < tc.want[j].Method
				}
				return tc.want[i].Path < tc.want[j].Path
			})

			require.Len(t, got, len(tc.want))
			for i := range tc.want {
				assert.Equal(t, tc.want[i], got[i])
			}
		})
	}
}

func TestDocumentedPath_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path model.DocumentedPath
		want string
	}{
		{
			name: "should format path with get method",
			path: model.DocumentedPath{Path: "/users", Method: "get"},
			want: "get_/users",
		},
		{
			name: "should format path with post method",
			path: model.DocumentedPath{Path: "/pets", Method: "post"},
			want: "post_/pets",
		},
		{
			name: "should format path with parameters",
			path: model.DocumentedPath{Path: "/users/{id}", Method: "get"},
			want: "get_/users/{id}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.path.String()
			assert.Equal(t, tc.want, got)
		})
	}
}
